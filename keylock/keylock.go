//go:build !solution

package keylock

import (
	"sync"
)

type KeyLock struct {
	mu sync.Mutex
	mp sync.Map
}

func New() *KeyLock {
	return &KeyLock{}
}

func (l *KeyLock) getCh(key string) chan struct{} {
	ch, _ := l.mp.LoadOrStore(key, make(chan struct{}, 1))
	return ch.(chan struct{})
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	firstBan := 0
	for {
		select {
		case <-cancel:
			return true, nil
		case l.getCh(keys[firstBan]) <- struct{}{}:
			<-l.getCh(keys[firstBan])

			l.mu.Lock()
		loop:
			for i := 0; i <= len(keys); i++ {
				if i == len(keys) {
					l.mu.Unlock()
					return false, func() {
						for i := 0; i < len(keys); i++ {
							<-l.getCh(keys[i])
						}
					}
				}
				select {
				case l.getCh(keys[i]) <- struct{}{}:
					continue loop
				default:
					firstBan = i
					for j := 0; j < i; j++ {
						<-l.getCh(keys[j])
					}
					l.mu.Unlock()
					break loop
				}
			}
		}
	}
}
