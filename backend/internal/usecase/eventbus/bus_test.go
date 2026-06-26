package eventbus

import (
	"sync"
	"testing"
	"time"
)

func TestSubscribeAndNotify(t *testing.T) {
	bus := NewBus()
	ch := bus.Subscribe("ep-1")

	bus.Notify("ep-1")

	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected notification on subscribed channel")
	}
}

func TestNotifyDoesNotBlockWhenBufferFull(t *testing.T) {
	bus := NewBus()
	ch := bus.Subscribe("ep-1")

	// Fill the buffer.
	bus.Notify("ep-1")
	// Second notify should not block even though the first hasn't been consumed.
	done := make(chan struct{})
	go func() {
		bus.Notify("ep-1")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Notify blocked when channel buffer was full")
	}

	// Drain the single buffered signal.
	<-ch
}

func TestNotifyRoutesToCorrectEpisode(t *testing.T) {
	bus := NewBus()
	ch1 := bus.Subscribe("ep-1")
	ch2 := bus.Subscribe("ep-2")

	bus.Notify("ep-1")

	select {
	case <-ch1:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected notification on ep-1 channel")
	}

	select {
	case <-ch2:
		t.Fatal("unexpected notification on ep-2 channel")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := NewBus()
	ch := bus.Subscribe("ep-1")
	bus.Unsubscribe("ep-1", ch)

	bus.Notify("ep-1")

	select {
	case <-ch:
		t.Fatal("should not receive notification after unsubscribe")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestUnsubscribeCleansEmptyEntry(t *testing.T) {
	bus := NewBus()
	ch := bus.Subscribe("ep-1")
	bus.Unsubscribe("ep-1", ch)

	bus.mu.RLock()
	defer bus.mu.RUnlock()
	if _, ok := bus.subs["ep-1"]; ok {
		t.Fatal("expected episode entry to be cleaned up after last unsubscribe")
	}
}

func TestConcurrentSubscribeNotify(t *testing.T) {
	bus := NewBus()
	const goroutines = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			ch := bus.Subscribe("ep-1")
			bus.Notify("ep-1")
			<-ch
			bus.Unsubscribe("ep-1", ch)
		}()
	}

	wg.Wait()
}

func TestMultipleSubscribers(t *testing.T) {
	bus := NewBus()
	ch1 := bus.Subscribe("ep-1")
	ch2 := bus.Subscribe("ep-1")

	bus.Notify("ep-1")

	for i, ch := range []chan struct{}{ch1, ch2} {
		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("subscriber %d did not receive notification", i)
		}
	}
}
