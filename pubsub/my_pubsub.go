//go:build !solution

package pubsub

import (
	"context"
	"errors"
	"sync"
)

var _ Subscription = (*MySubscription)(nil)

var (
	closeErr = errors.New("pubSub is close")
	doneCtx  = errors.New("context is close")
)

type topic struct {
	msgs    []interface{}
	chs     []chan struct{}
	chClose chan struct{}
	muCh    sync.Mutex
	muMsg   sync.Mutex
	wg      *sync.WaitGroup
}

type MySubscription struct {
	ctx              *context.Context
	firstUnprocessed int
	cancel           func()
}

func (t *topic) newSubsribe(ind int, cb MsgHandler) MySubscription {
	ctx, cancel := context.WithCancel(context.Background())
	s := MySubscription{ctx: &ctx, firstUnprocessed: ind, cancel: cancel}

	sendMassage := func() {
		t.muMsg.Lock()
		msg := t.msgs[s.firstUnprocessed]
		t.muMsg.Unlock()
		cb(msg)
		s.firstUnprocessed++
	}
	t.wg.Add(1)
	go func() {
		for {
			t.muCh.Lock()
			ch := t.chs[s.firstUnprocessed]
			t.muCh.Unlock()
			select {
			case <-ctx.Done():
				t.wg.Done()
				return
			case <-ch:
				select {
				case <-ctx.Done():
					return
				default:
					sendMassage()
				}
			case <-t.chClose:
				select {
				case <-ch:
					sendMassage()
				default:
					cancel()
				}
			}
		}
	}()
	return s
}

func (s *MySubscription) Unsubscribe() {
	s.cancel()
}

var _ PubSub = (*MyPubSub)(nil)

type MyPubSub struct {
	topics  map[string]*topic
	mu      sync.Mutex
	chClose chan struct{}
	wg      *sync.WaitGroup
	isDone  chan struct{}
}

func NewPubSub() PubSub {
	p := &MyPubSub{mu: sync.Mutex{},
		topics:  make(map[string]*topic),
		chClose: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},
		isDone:  make(chan struct{}, 1),
	}
	return p
}

func (p *MyPubSub) Subscribe(subj string, cb MsgHandler) (Subscription, error) {
	select {
	case <-p.chClose:
		return nil, closeErr
	default:
		p.mu.Lock()
		if _, ok := p.topics[subj]; !ok {
			p.topics[subj] = &topic{
				chs:     []chan struct{}{make(chan struct{}, 1)},
				chClose: p.chClose,
				wg:      p.wg,
			}
		}
		topic := p.topics[subj]
		topic.muCh.Lock()
		ind := len(topic.chs) - 1
		topic.muCh.Unlock()
		s := p.topics[subj].newSubsribe(ind, cb)
		p.mu.Unlock()
		return &s, nil
	}
}

func (p *MyPubSub) Publish(subj string, msg interface{}) error {
	select {
	case <-p.chClose:
		return closeErr
	default:
		p.mu.Lock()
		topic := p.topics[subj]

		topic.muMsg.Lock()
		topic.msgs = append(topic.msgs, msg)
		topic.muMsg.Unlock()

		topic.muCh.Lock()
		topic.chs = append(topic.chs, make(chan struct{}, 1))
		close(topic.chs[len(topic.chs)-2])
		topic.muCh.Unlock()
		p.mu.Unlock()
		return nil
	}
}

func (p *MyPubSub) Close(ctx context.Context) error {
	select {
	case <-p.chClose:
		return closeErr
	default:
		close(p.chClose)
		go p.wait()
	}
	select {
	case <-ctx.Done():
		return doneCtx
	case <-p.isDone:
		return nil
	}
}

func (p *MyPubSub) wait() {
	p.wg.Wait()
	close(p.isDone)
}
