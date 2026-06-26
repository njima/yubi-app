package eventbus

import "sync"

// Bus is a generic in-process pub/sub bus. Subscribers receive a signal on a
// buffered channel whenever Notify is called for the key they subscribed to.
type Bus struct {
	mu   sync.RWMutex
	subs map[string]map[chan struct{}]struct{}
}

func NewBus() *Bus {
	return &Bus{
		subs: make(map[string]map[chan struct{}]struct{}),
	}
}

// Subscribe returns a buffered (cap 1) channel that receives a signal
// each time Notify is called for the given key.
func (b *Bus) Subscribe(key string) chan struct{} {
	ch := make(chan struct{}, 1)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.subs[key] == nil {
		b.subs[key] = make(map[chan struct{}]struct{})
	}
	b.subs[key][ch] = struct{}{}
	return ch
}

// Unsubscribe removes the channel from the subscriber set and cleans up
// empty entries.
func (b *Bus) Unsubscribe(key string, ch chan struct{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if m, ok := b.subs[key]; ok {
		delete(m, ch)
		if len(m) == 0 {
			delete(b.subs, key)
		}
	}
}

// Notify sends a non-blocking signal to every subscriber of the given key.
func (b *Bus) Notify(key string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subs[key] {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
