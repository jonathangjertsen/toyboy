package util

import (
	"container/list"
	"context"
	"log/slog"
	"time"
)

type Broker[T any] struct {
	subscribers *list.List
	logger      *slog.Logger
	timeout     time.Duration
	In          chan T
	sub         chan chan *Subscriber[T]
	unsub       chan *Subscriber[T]
}

func NewBroker[T any](ctx context.Context, logger *slog.Logger, timeout time.Duration) *Broker[T] {
	b := &Broker[T]{
		subscribers: list.New(),
		logger:      logger,
		timeout:     timeout,
		In:          make(chan T, 1),
		sub:         make(chan chan *Subscriber[T]),
		unsub:       make(chan *Subscriber[T]),
	}
	go b.worker(ctx)
	return b
}

func (b *Broker[T]) worker(ctx context.Context) {
	for {
		select {
		case v := <-b.In:
			// handle incoming message
			for e := b.subscribers.Front(); e.Next() != nil; e = e.Next() {
				go func(sub *Subscriber[T], v T) {
					select {
					case sub.C <- v:
					case <-time.After(b.timeout):
						b.logger.ErrorContext(ctx, "Publish timed out", "timeout", b.timeout)
					}
				}(e.Value.(*Subscriber[T]), v)
			}
		case ch := <-b.sub:
			// subscribe
			sub := &Subscriber[T]{}
			sub.parent = b
			sub.C = make(chan T, 1)
			sub.elem = b.subscribers.PushBack(sub)
			ch <- sub
		case sub := <-b.unsub:
			// unsubscribe
			b.subscribers.Remove(sub.elem)
			close(sub.C)
		}
	}
}

func (b *Broker[T]) Subscribe() *Subscriber[T] {
	ch := make(chan *Subscriber[T], 1)
	b.sub <- ch
	return <-ch
}

type Subscriber[T any] struct {
	parent *Broker[T]
	elem   *list.Element
	C      chan T
}

func (s *Subscriber[T]) Unsubscribe() {
	s.parent.unsub <- s
}
