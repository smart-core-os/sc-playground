package simulated

import (
	"github.com/smart-core-os/sc-playground/internal/util/clock"
	"sync"
	"time"
)

// Clock is an implementation of clock.Clock that is manually stepped through time by calling Advance.
// The value returned by Now will only change when Advance is called.
// The Clock can be stopped (as defined in clock.Clock) by calling Stop. After it has been stopped, it is forbidden
// to call any methods on the clock. You must call Stop after you have finished using the Clock, or a goroutine
// will be leaked.
type Clock struct {
	// state accessible directly from the clock.Clock functions, to allow Now to work.
	timeM sync.RWMutex
	time  time.Time

	// communicate with worker
	advanceCh chan advanceRequest // close this channel to stop the clock
	pendingCh chan timer
}

var _ clock.Clock = &Clock{}

// NewClock constructs a new Clock, with start as the initial value of Now.
func NewClock(start time.Time) *Clock {
	c := &Clock{
		time:      start,
		advanceCh: make(chan advanceRequest),
		pendingCh: make(chan timer),
	}

	go c.run()

	return c
}

func (c *Clock) Now() time.Time {
	c.timeM.RLock()
	defer c.timeM.RUnlock()
	return c.time
}

func (c *Clock) At(t time.Time) <-chan time.Time {
	// channel with buffer lets the consumer leak the channel without us locking up
	ch := make(chan time.Time, 1)
	c.pendingCh <- timer{
		ch:   ch,
		next: t,
	}
	return ch
}

func (c *Clock) After(d time.Duration) <-chan time.Time {
	t := c.Now().Add(d)
	return c.At(t)
}

func (c *Clock) Every(d time.Duration) clock.Ticker {
	first := c.Now().Add(d)
	ch := make(chan time.Time)
	stop := make(chan struct{})
	c.pendingCh <- timer{
		ch:       ch,
		next:     first,
		stop:     stop,
		interval: d,
	}
	return ticker{
		ch:   ch,
		stop: stop,
	}
}

// Advance will step the simulation forward by d.
// If there are any pending waits (such as calls to After), ok == true and next will contain the simulation duration
// until the next pending wait is due to complete. Otherwise, ok == false and next will be an unspecified value.
func (c *Clock) Advance(d time.Duration) (next time.Duration, ok bool) {
	resCh := make(chan advanceResponse)
	c.advanceCh <- advanceRequest{
		duration: d,
		response: resCh,
	}
	res := <-resCh
	return res.untilNext, res.ok
}

// Stop will shut down the Clock. It will stop the worker goroutine used internally, and closes the channels for
// all pending waits. You must call Stop after you have finished using the Clock, or a goroutine will leak.
// After calling Stop, all future calls to any methods on the Clock are forbidden. They will not return meaningful
// results and may panic.
func (c *Clock) Stop() {
	close(c.advanceCh)
}

func (c *Clock) run() {
	var pending []timer
	defer func() {
		// cleanup pending
		for _, p := range pending {
			close(p.ch)
		}
	}()

	for {
		select {
		case add := <-c.pendingCh:
			// add new pending
			pending = append(pending, add)

		case req, ok := <-c.advanceCh:
			if !ok {
				// clock stopped
				return
			}

			pending = c.advance(pending, req)
		}
	}
}

func (c *Clock) advance(old []timer, req advanceRequest) (new []timer) {
	// advance current time
	now := c.Now().Add(req.duration)
	c.setTime(now)

	// process all pending timers
	for _, pending := range old {
		// if we haven't reached the right time yet, just carry it over
		if pending.next.After(now) {
			new = append(new, pending)
			continue
		}

		if pending.stop == nil {
			// non-repeating timer, send and discard
			pending.ch <- now
		} else {
			// repeating timer, add to new only if it's not stopped
			select {
			case <-pending.stop:
				// the timer was stopped, don't add to new
			case pending.ch <- now:
				// don't schedule the next one for the past, that could cause large backlogs to build up
				// instead, just skip over excess iterations
				for !pending.next.After(now) {
					pending.next = pending.next.Add(pending.interval)
				}
				new = append(new, pending)
			}
		}

	}

	// find the next timer that will fire
	var res advanceResponse
	for _, pending := range new {
		untilNext := pending.next.Sub(now)
		if !res.ok {
			res.ok = true
			res.untilNext = untilNext
		} else if untilNext < res.untilNext {
			res.untilNext = untilNext
		}
	}
	req.response <- res

	return
}

func (c *Clock) setTime(newT time.Time) {
	c.timeM.Lock()
	defer c.timeM.Unlock()
	c.time = newT
}

type timer struct {
	ch   chan<- time.Time
	next time.Time

	stop     <-chan struct{} //closed elsewhere to request a stop. nil if this timer is not repeating.
	interval time.Duration   //only used if stop != nil
}

type advanceRequest struct {
	duration time.Duration
	response chan<- advanceResponse
}

type advanceResponse struct {
	untilNext time.Duration
	ok        bool
}

type ticker struct {
	ch   <-chan time.Time
	stop chan<- struct{}
}

func (t ticker) C() <-chan time.Time {
	return t.ch
}

func (t ticker) Stop() {
	t.stop <- struct{}{}
	close(t.stop)
}

// SimulateFor will advance the simulated clock through the total duration specified.
// The clock will only ever be advanced to its next event. If the clock has no pending event, then it will be advanced
// by defaultStep.
// SimulateFor does not insert any real-time delays, so the simulation will complete as fast as possible.
// This function is useful in tests, or offline simulations.
func SimulateFor(clk *Clock, duration time.Duration, defaultStep time.Duration) {
	var next time.Duration
	for duration > 0 {
		duration -= next

		var ok bool
		next, ok = clk.Advance(next)
		if !ok {
			next = defaultStep
		}
	}
}

// SimulateUntilIdle will repeatedly advance the clock to the next pending event, until there are no more pending
// events.
func SimulateUntilIdle(clk *Clock) {
	var next time.Duration
	for {
		var ok bool
		next, ok = clk.Advance(next)

		if !ok {
			break
		}
	}
}
