//go:build !solution

package batcher

import (
	"gitlab.com/slon/shad-go/batcher/slow"
	"sync"
)

type Batcher struct {
	mu          sync.Mutex
	value       *slow.Value
	readChannel chan struct{}
	running     chan struct{}
	v           interface{}
}

func NewBatcher(v *slow.Value) *Batcher {
	return &Batcher{
		value:       v,
		readChannel: make(chan struct{}, 1),
		running:     make(chan struct{}, 1),
	}
}

func (b *Batcher) Load() interface{} {
	select {
	case b.readChannel <- struct{}{}:
		b.mu.Lock()
		defer func() {
			b.mu.Unlock()
			b.running <- struct{}{}
		}()
		b.v = b.value.Load()
		<-b.readChannel

		return b.v
	case <-b.running:
		b.mu.Lock()
		defer b.mu.Unlock()
		b.running <- struct{}{}
		return b.v
	}
}
