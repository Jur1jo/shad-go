//go:build !solution

package batcher

import (
	"fmt"
	"gitlab.com/slon/shad-go/batcher/slow"
	"sync"
	"sync/atomic"
)

type Batcher struct {
	v    *slow.Value
	mu   sync.Mutex
	muRw sync.RWMutex
	cur  any
	ch   atomic.Value
	lst  atomic.Value
}

func NewBatcher(v *slow.Value) *Batcher {
	b := Batcher{v: v}
	b.ch.Store(make(chan struct{}))
	return &b
}

func (b *Batcher) Load() any {
	ch := b.ch.Load().(chan struct{})
	if b.mu.TryLock() {
		b.ch.Store(make(chan struct{}))

		b.muRw.Lock()
		b.cur = b.v.Load()
		b.muRw.Unlock()

		close(ch)
		fmt.Println(b.ch.Load().(chan struct{}), ch)
		b.mu.Unlock()
	}
	fmt.Println("!", ch)
	<-ch
	b.muRw.RLock()
	defer b.muRw.RUnlock()
	return b.cur
}
