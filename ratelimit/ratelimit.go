//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	erase chan time.Time
	stop  chan struct{}
	ch    chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	limiter := Limiter{erase: make(chan time.Time, 2), stop: make(chan struct{}), ch: make(chan struct{}, maxCount)}

	go func() {
		for {
			select {
			case <-limiter.stop:
				return
			case t := <-limiter.erase:
				timer := time.NewTimer(t.Add(interval).Sub(time.Now()))
				select {
				case <-timer.C:
					<-limiter.ch
				case <-limiter.stop:
					timer.Stop()
					return
				}
				timer.Stop()
			}
		}
	}()
	return &limiter
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-l.stop:
		return ErrStopped
	default:
		select {
		case l.ch <- struct{}{}:
			l.erase <- time.Now()
			return nil
		case <-l.stop:
			return ErrStopped
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (l *Limiter) Stop() {
	fmt.Println("MEOW")
	close(l.stop)
}
