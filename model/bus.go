package model

import (
	"fmt"
	"sync"
)

type BusType interface {
	uint8 | uint16
}

type Bus[T BusType] struct {
	name       string
	writeValue T
	readValue  T
	readSync   *sync.Cond
	writeSync  *sync.Cond
	clock      *Clock
	rising     bool
}

func NewBus[T BusType](clock *Clock, name string) *Bus[T] {
	bus := &Bus[T]{
		clock:     clock,
		readSync:  sync.NewCond(&sync.Mutex{}),
		writeSync: sync.NewCond(&sync.Mutex{}),
		name:      name,
	}
	clock.AddRiseCallback(func(c Cycle) {
		bus.readSync.L.Lock()
		bus.writeSync.L.Lock()
		bus.rising = true
		bus.writeSync.L.Unlock()
		bus.readSync.L.Unlock()

		bus.writeSync.Broadcast()
	})
	clock.AddFallCallback(func(c Cycle) {
		bus.readSync.L.Lock()
		bus.writeSync.L.Lock()
		bus.readValue = bus.writeValue
		bus.rising = false
		bus.writeSync.L.Unlock()
		bus.readSync.L.Unlock()

		bus.readSync.Broadcast()
	})
	return bus
}

func (b *Bus[T]) Write(v T) {
	b.writeSync.L.Lock()
	for !b.rising {
		b.writeSync.Wait()
	}
	b.writeValue = v
	b.writeSync.L.Unlock()

	fmt.Printf("wrote 0x%x to %s\n", v, b.name)
}

func (b *Bus[T]) Read() T {
	b.readSync.L.Lock()
	for b.rising {
		b.readSync.Wait()
	}
	v := b.readValue
	b.readSync.L.Unlock()

	fmt.Printf("read 0x%x from %s\n", v, b.name)

	return v
}
