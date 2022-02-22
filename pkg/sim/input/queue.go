package input

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/smart-core-os/sc-playground/pkg/timeline"
)

// ErrQueueClosed is returned by blocking calls when the Queue is closed.
var ErrQueueClosed = errors.New("input queue closed")

// Queue captures input events so they can be processed as a bulk at a later time.
// Queue is a Dispatcher, for event sources to contribute, and a Capturer for consumers to receive the events.
type Queue struct {
	buffer []event

	closed    chan struct{}
	closeOnce sync.Once

	comm   chan interface{}
	notify chan struct{}
}

func NewQueue() *Queue {
	q := &Queue{
		buffer: make([]event, 0, 1000),
		closed: make(chan struct{}),
		comm:   make(chan interface{}),
		notify: make(chan struct{}),
	}

	go q.run()

	return q
}

// Wait blocks until the queue has done something or ctx is done.
// Only one waiting caller will be woken.
// If q is closed, returns immediately.
func (q *Queue) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-q.notify:
	}
}

func (q *Queue) Dispatch(ctx context.Context, t time.Time, e ...interface{}) (func(), error) {
	msg := event{
		time:      t,
		value:     e,
		done:      make(chan struct{}),
		committed: make(chan struct{}),
		rejected:  make(chan error),
	}
	closeOnce := sync.Once{}
	done := func() {
		closeOnce.Do(func() {
			close(msg.done)
		})
	}
	select {
	case <-ctx.Done():
		done()
		return func() {}, ctx.Err()
	case <-q.closed:
		return func() {}, ErrQueueClosed
	case q.comm <- msg:
		select {
		case <-ctx.Done():
			done()
			return func() {}, ctx.Err()
		case <-q.closed:
			return func() {}, ErrQueueClosed
		case <-msg.committed:
			return done, nil
		case err := <-msg.rejected:
			done()
			return func() {}, err
		}
	}
}

func (q *Queue) Capture(ctx context.Context, tl timeline.AddTL) (Session, error) {
	s := &session{
		ctx:       ctx,
		tl:        tl,
		tlUpdated: make(chan struct{}),
		commit:    make(chan struct{}),
		reject:    make(chan error),
		done:      make(chan struct{}),
	}
	select {
	case <-q.closed:
		return nil, ErrQueueClosed
	case q.comm <- s:
		select {
		case <-q.closed:
			return nil, ErrQueueClosed
		case <-s.tlUpdated:
			return s, nil
		}
	}
}

func (q *Queue) Close() error {
	q.closeOnce.Do(func() {
		close(q.closed)
		close(q.notify)
	})
	return nil
}

func (q *Queue) run() {
	// Run is used to add thread safety. All interaction with q.buffer happen in this method, an alternative to locks.

	for {
		select {
		case <-q.closed:
			return
		case msg := <-q.comm:
			switch t := msg.(type) {
			case event:
				q.buffer = append(q.buffer, t)
			case *session:
				// process all the buffered events and signal waiting routines
				// 1. update the timeline
				for _, e := range q.buffer {
					for _, val := range e.value {
						t.tl.Add(e.time, val)
					}
				}

				// 2. notify that we updated the tl, causes q.Capture to return
				close(t.tlUpdated)

				// 3. wait for the session to be settled, or for it to be closed in another way
				//    This bit causes Dispatch calls to return
				select {
				case <-t.ctx.Done():
					q.rejectAll(t.ctx.Err())
				case <-q.closed:
					q.rejectAll(ErrQueueClosed)
				case err := <-t.reject:
					q.rejectAll(err)
				case <-t.commit:
					q.commitAll()
				}

				// 4. Wait for Done to be called on each Dispatch
				for _, e := range q.buffer {
					select {
					case <-q.closed:
					case <-e.done:
					}
				}

				// 5. empty the buffer
				q.buffer = q.buffer[:0]

				// 6. Close the session, causes Session.Commit or .Reject to return
				close(t.done)
			}
			// we just did something, notify a waiting party
			select {
			case <-q.closed: // might have got here because the queue was closed
			case q.notify <- struct{}{}:
			default:
			}
		}
	}
}

func (q *Queue) commitAll() {
	for _, e := range q.buffer {
		close(e.committed)
	}
}

func (q *Queue) rejectAll(err error) {
	for _, e := range q.buffer {
		select {
		case <-q.closed:
		case e.rejected <- err:
		}
	}
}

type event struct {
	time      time.Time
	value     []interface{}
	done      chan struct{} // closed by Dispatch response func
	committed chan struct{} // closed by Capture.Session.Commit
	rejected  chan error    // filled by Capture.Session.Reject
}

type session struct {
	ctx context.Context
	tl  timeline.AddTL

	tlUpdated  chan struct{}
	commit     chan struct{}
	reject     chan error
	settleOnce sync.Once // guards commit and reject
	done       chan struct{}
}

func (s *session) Commit() {
	s.settleOnce.Do(func() {
		close(s.commit)
	})
	<-s.done
}

func (s *session) Reject(err error) {
	s.settleOnce.Do(func() {
		s.reject <- err
	})
	<-s.done
}
