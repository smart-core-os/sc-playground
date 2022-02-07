package input

import (
	"context"
	"testing"
	"time"

	"github.com/smart-core-os/sc-playground/pkg/timeline"
	"github.com/smart-core-os/sc-playground/pkg/timeline/skiplisttl"
)

func TestQueue(t *testing.T) {
	q := NewQueue()
	t.Cleanup(func() {
		_ = q.Close()
	})
	ctx := context.Background()

	// q.Dispatch blocks until Session.Commit or Reject is called
	// Session.Commit and Reject block until Dispatch.done is called
	// q.Capture doesn't really block, but needs to wait for events to be dispatched.

	// dispatch routine
	var dispatchError error
	go func() {
		done, err := q.Dispatch(ctx, time.Unix(1000, 0), "one", "two")
		dispatchError = err
		defer func() {
			done()
		}()
	}()

	tl := skiplisttl.New()
	// Wait for the dispatch to be recorded
	q.Wait(ctx)

	// Begin the input session
	session, err := q.Capture(ctx, tl)
	if err != nil {
		t.Cleanup(func() { session.Reject(err) })
		t.Fatal(err)
	}

	// this blocks waiting for the done func returned by Dispatch to be called
	session.Commit()

	// check the timeline
	if dispatchError != nil {
		t.Fatal(dispatchError)
	}
	want := skiplisttl.New()
	want.Add(time.Unix(1000, 0), "one", "two")
	if equal, reason := timeline.EqualExplain(want, tl); !equal {
		t.Fatalf(reason)
	}
}
