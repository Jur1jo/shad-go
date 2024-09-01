//go:build !solution

package dupcall

import (
	"context"
	"sync"
	"sync/atomic"
)

type Call struct {
	hasRunFunc    chan struct{}
	cntRunFunc    atomic.Int32
	isInit        atomic.Bool
	mu            sync.Mutex
	res           interface{}
	err           error
	cancelCtxFunc func()
}

func (o *Call) init() {
	o.isInit.Store(true)
	o.hasRunFunc = make(chan struct{}, 1)
	o.hasRunFunc <- struct{}{}
}

func (o *Call) doDone() {
	o.mu.Lock()
	o.cntRunFunc.Add(-1)
	if o.cntRunFunc.Load() == 0 {
		o.cancelCtxFunc()
	}
	o.mu.Unlock()
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.mu.Lock()
	if !o.isInit.Load() {
		o.init()
	}
	o.cntRunFunc.Add(1)
	o.mu.Unlock()
	select {
	case _, ok := <-o.hasRunFunc:
		if ok {
			ctx, cancel := context.WithCancel(context.Background())
			o.cancelCtxFunc = cancel
			go func() {
				o.res, o.err = cb(ctx)
				close(o.hasRunFunc)
			}()
		}
		select {
		case <-ctx.Done():
			o.doDone()
			return nil, ctx.Err()
		case <-o.hasRunFunc:
			o.mu.Lock()
			o.cntRunFunc.Add(-1)
			if o.cntRunFunc.Load() == 0 {
				o.hasRunFunc = make(chan struct{}, 1)
				o.hasRunFunc <- struct{}{}
			}
			o.mu.Unlock()
			return o.res, o.err
		}
	case <-ctx.Done():
		o.doDone()
		return nil, ctx.Err()
	}
}
