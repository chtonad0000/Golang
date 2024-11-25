//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	maxCount      int
	interval      time.Duration
	ticker        *time.Ticker
	stopSignal    chan bool
	acquire       chan bool
	acquireLocker chan bool
	stopped       bool
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	var acquire chan bool
	if interval == 0 {
		acquire = make(chan bool, 1)
	} else {
		acquire = make(chan bool, maxCount)
	}

	limiter := &Limiter{
		maxCount:      maxCount,
		interval:      interval,
		stopSignal:    make(chan bool, 1),
		acquire:       acquire,
		acquireLocker: make(chan bool, 1),
	}
	if interval > 0 {
		limiter.ticker = time.NewTicker(interval)
		go func() {
			for {
				select {
				case <-limiter.ticker.C:
					length := len(limiter.acquire)
					for length > 0 {
						select {
						case <-limiter.acquire:
							length--
						default:
							length--
							continue
						}
					}
				case <-limiter.stopSignal:
					limiter.ticker.Stop()
					return
				}
			}
		}()
	}
	return limiter
}

func (l *Limiter) Acquire(ctx context.Context) error {
	for {
		select {
		case <-l.stopSignal:

			return ErrStopped
		default:
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				select {
				case l.acquire <- true:
					if l.interval == 0 {
						<-l.acquire
					}
					return nil
				default:
				}
			}
		}
	}
}

func (l *Limiter) Stop() {
	l.stopSignal <- true
}
