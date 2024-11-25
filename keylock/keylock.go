//go:build !solution

package keylock

import "sync"

type KeyLock struct {
	mu         sync.Mutex
	dangerZone map[string]chan struct{}
}

func New() *KeyLock {
	return &KeyLock{dangerZone: make(map[string]chan struct{})}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	var blockedCh []chan struct{}
	l.mu.Lock()
	for _, key := range keys {
		if ch, exists := l.dangerZone[key]; exists {
			blockedCh = append(blockedCh, ch)
		}
	}
	if len(blockedCh) > 0 {
		l.mu.Unlock()
		for _, ch := range blockedCh {
			select {
			case <-ch:
			case <-cancel:
				return true, nil
			}
		}
		return l.LockKeys(keys, cancel)
	}
	for _, key := range keys {
		l.dangerZone[key] = make(chan struct{})
	}
	l.mu.Unlock()
	return false, func() {
		l.mu.Lock()
		defer l.mu.Unlock()
		for _, key := range keys {
			close(l.dangerZone[key])
			delete(l.dangerZone, key)
		}
	}
}
