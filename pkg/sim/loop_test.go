package sim

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/smart-core-os/sc-playground/pkg/sim/input"
	"github.com/smart-core-os/sc-playground/pkg/timeline"
	"github.com/smart-core-os/sc-playground/pkg/timeline/skiplisttl"
)

func TestLoop_Advance(t *testing.T) {
	scrubber := &TestScrubber{}
	capturer := &TestCapturer{}
	tl := skiplisttl.New()
	loop := &Loop{
		Model: scrubber,
		Input: capturer,
		TL:    tl,
	}

	ctx := context.Background()
	err := loop.Advance(ctx, time.Unix(1000, 0))
	if err != nil {
		t.Fatal(err)
	}

	scrubber.assertCalls(t, time.Unix(1000, 0))
	capturer.assertCalls(t, &CaptureCall{ctx, tl, &TestSession{Committed: true}})
}

type TestScrubber struct {
	calls []time.Time
}

func (ts *TestScrubber) Scrub(t time.Time) error {
	ts.calls = append(ts.calls, t)
	return nil
}

func (ts *TestScrubber) assertCalls(tt *testing.T, t ...time.Time) {
	if diff := cmp.Diff(t, ts.calls); diff != "" {
		tt.Fatalf("Scrubber (-want, +got)\n%v", diff)
	}
}

type TestCapturer struct {
	calls []*CaptureCall
}

func (t *TestCapturer) assertCalls(tt *testing.T, calls ...*CaptureCall) {
	if diff := cmp.Diff(calls, t.calls, cmpopts.IgnoreUnexported(CaptureCall{})); diff != "" {
		tt.Fatalf("Capturer (-want, +got)\n%v", diff)
	}
}

func (t *TestCapturer) Capture(ctx context.Context, tl timeline.AddTL) (input.Session, error) {
	session := &TestSession{}
	t.calls = append(t.calls, &CaptureCall{ctx, tl, session})
	return session, nil
}

type CaptureCall struct {
	ctx     context.Context
	tl      timeline.AddTL
	Session *TestSession
}

type TestSession struct {
	Committed bool
	Rejected  error
}

func (t *TestSession) Commit() {
	t.Committed = true
}

func (t *TestSession) Reject(err error) {
	t.Rejected = err
}
