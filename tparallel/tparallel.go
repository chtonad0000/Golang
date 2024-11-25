//go:build !solution

package tparallel

import (
	"sync"
)

type T struct {
	block      chan bool
	wg         *sync.WaitGroup
	mu         *sync.Mutex
	parentDone chan bool
	signal     chan bool
}

func (t *T) Parallel() {
	t.signal <- true
	<-t.parentDone
}

func (t *T) Run(subtest func(t *T)) {
	subTest := &T{
		parentDone: t.block,
		wg:         &sync.WaitGroup{},
		block:      make(chan bool),
		mu:         &sync.Mutex{},
		signal:     make(chan bool),
	}
	t.mu.Lock()
	t.wg.Add(1)
	go func() {
		subtest(subTest)
		close(subTest.block)
		t.wg.Done()
	}()
	t.mu.Unlock()
	select {
	case <-subTest.block:
		subTest.wg.Wait()

	case <-subTest.signal:
	}
}

func Run(topTests []func(t *T)) {
	globalCh := make(chan bool)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for _, topTest := range topTests {
		t := &T{
			block:      make(chan bool),
			parentDone: globalCh,
			mu:         &sync.Mutex{},
			signal:     make(chan bool),
			wg:         &sync.WaitGroup{},
		}
		mu.Lock()
		wg.Add(1)
		go func() {
			topTest(t)
			close(t.block)
			wg.Done()
		}()
		mu.Unlock()
		select {
		case <-t.signal:
		case <-t.block:
			t.wg.Wait()
		}
	}
	close(globalCh)
	wg.Wait()
}
