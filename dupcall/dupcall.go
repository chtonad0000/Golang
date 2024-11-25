package dupcall

import (
	"context"
	"fmt"
	"sync"
)

type Call struct {
	contexts map[string]*callState
	mu       sync.Mutex
}

type callState struct {
	done    chan struct{}
	waiters int
	value   interface{}
	err     error
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.mu.Lock()
	if o.contexts == nil {
		o.contexts = make(map[string]*callState)
	}
	o.mu.Unlock()
	key := fmt.Sprintf("%p", cb)
	o.mu.Lock()
	state, exists := o.contexts[key]
	if exists {
		state.waiters++
		o.mu.Unlock()
		select {
		case <-state.done:
			return state.value, state.err
		case <-ctx.Done():
			o.mu.Lock()
			state.waiters--
			if state.waiters == 0 {
				delete(o.contexts, key)
			}
			o.mu.Unlock()
			return nil, ctx.Err()
		}
	}
	state = &callState{
		done:    make(chan struct{}),
		waiters: 1,
	}
	o.contexts[key] = state
	o.mu.Unlock()
	go func() {
		defer close(state.done)
		state.value, state.err = cb(ctx)
		o.mu.Lock()
		if state.waiters == 1 {
			delete(o.contexts, key)
		} else {
			state.waiters--
		}
		o.mu.Unlock()
	}()
	select {
	case <-ctx.Done():
		o.mu.Lock()
		state.waiters--
		if state.waiters == 0 {
			delete(o.contexts, key)
		}
		o.mu.Unlock()
		return nil, ctx.Err()
	case <-state.done:
		return state.value, state.err
	}
}
