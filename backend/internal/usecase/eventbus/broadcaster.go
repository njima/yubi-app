package eventbus

import (
	"context"
	"log/slog"
	"sync"
)

// Broadcaster fans out the result of a single fetch to all subscribers of a
// key. When a Bus notification arrives for a key, the fetch function is called
// once and the result is sent to every subscriber channel.
type Broadcaster[T any] struct {
	bus    *Bus
	fetch  func(ctx context.Context, key string) (T, error)
	parent context.Context

	mu   sync.Mutex
	keys map[string]*broadcasterKey[T]
}

type broadcasterKey[T any] struct {
	subs   map[chan T]struct{}
	cancel context.CancelFunc
}

// NewBroadcaster creates a Broadcaster that listens on bus notifications and
// calls fetch once per notification, broadcasting the result to all
// subscribers of that key. The parent context is used to derive per-key
// goroutine contexts; cancelling it stops all broadcaster goroutines.
func NewBroadcaster[T any](parent context.Context, bus *Bus, fetch func(ctx context.Context, key string) (T, error)) *Broadcaster[T] {
	return &Broadcaster[T]{
		bus:    bus,
		fetch:  fetch,
		parent: parent,
		keys:   make(map[string]*broadcasterKey[T]),
	}
}

// Subscribe returns a channel that receives the fetched value each time the
// underlying bus fires for the given key. The broadcaster for this key is
// started lazily on the first subscriber and stopped when the last subscriber
// unsubscribes.
func (b *Broadcaster[T]) Subscribe(key string) chan T {
	ch := make(chan T, 1)

	b.mu.Lock()
	defer b.mu.Unlock()

	bk, ok := b.keys[key]
	if !ok {
		ctx, cancel := context.WithCancel(b.parent)
		bk = &broadcasterKey[T]{
			subs:   make(map[chan T]struct{}),
			cancel: cancel,
		}
		b.keys[key] = bk
		go b.loop(ctx, key)
	}
	bk.subs[ch] = struct{}{}
	return ch
}

// Unsubscribe removes ch from the subscriber set for the given key. If this
// was the last subscriber, the broadcaster goroutine for this key is stopped.
func (b *Broadcaster[T]) Unsubscribe(key string, ch chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	bk, ok := b.keys[key]
	if !ok {
		return
	}
	delete(bk.subs, ch)
	if len(bk.subs) == 0 {
		bk.cancel()
		delete(b.keys, key)
	}
}

// loop subscribes to the underlying bus for the given key and, on each
// notification, fetches data once and fans it out to all subscribers.
func (b *Broadcaster[T]) loop(ctx context.Context, key string) {
	busCh := b.bus.Subscribe(key)
	defer b.bus.Unsubscribe(key, busCh)

	for {
		select {
		case <-ctx.Done():
			return
		case <-busCh:
			val, err := b.fetch(ctx, key)
			if err != nil {
				slog.Error("broadcaster: fetch failed", "key", key, "error", err)
				continue
			}

			// Snapshot subscribers under lock, then fan out without holding it.
			b.mu.Lock()
			bk, ok := b.keys[key]
			var snapshot []chan T
			if ok {
				snapshot = make([]chan T, 0, len(bk.subs))
				for ch := range bk.subs {
					snapshot = append(snapshot, ch)
				}
			}
			b.mu.Unlock()

			for _, ch := range snapshot {
				select {
				case ch <- val:
				default:
					// Drain stale value and replace with fresh one.
					select {
					case <-ch:
					default:
					}
					select {
					case ch <- val:
					default:
					}
				}
			}
		}
	}
}
