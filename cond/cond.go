//go:build !solution

package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

type queueChan struct {
	stIn  []chan struct{}
	stOut []chan struct{}
}

func (q *queueChan) rebuild() {
	if len(q.stOut) == 0 {
		for len(q.stIn) > 0 {
			q.stOut = append(q.stOut, q.stIn[len(q.stIn)-1])
			q.stIn = q.stIn[:len(q.stIn)-1]
		}
	}
}

func (q *queueChan) pop() {
	q.rebuild()
	q.stOut = q.stOut[:len(q.stOut)-1]
}

func (q *queueChan) front() chan struct{} {
	q.rebuild()
	return q.stOut[len(q.stOut)-1]
}

func (q *queueChan) push(ch chan struct{}) {
	q.stIn = append(q.stIn, ch)
}

func (q *queueChan) len() int {
	return len(q.stIn) + len(q.stOut)
}

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *sync.Mutex or *sync.RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
type Cond struct {
	L    Locker
	q    queueChan
	lock chan struct{}
}

// New returns a new Cond with Locker l.
func New(l Locker) *Cond {
	return &Cond{L: l, lock: make(chan struct{}, 1)}
}

// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//	c.L.Lock()
//	for !condition() {
//	    c.Wait()
//	}
//	... make use of condition ...
//	c.L.Unlock()
func (c *Cond) Wait() {
	ch := make(chan struct{}, 1)

	c.lock <- struct{}{}
	c.q.push(ch)
	<-c.lock

	c.L.Unlock()
	<-ch
	c.L.Lock()
}

func (c *Cond) wake(n int) {
	c.lock <- struct{}{}
	if n == -1 {
		n = c.q.len()
	}
	queueLen := c.q.len()
	for i := 0; i < min(n, queueLen); i++ {
		ch := c.q.front()
		c.q.pop()
		ch <- struct{}{}
	}
	<-c.lock
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Signal() {
	c.wake(1)
}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast() {
	c.wake(-1)
}
