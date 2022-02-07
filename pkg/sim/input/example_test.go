package input

import (
	"context"
	"fmt"
	"time"

	"github.com/smart-core-os/sc-playground/pkg/timeline"
)

func Example_dispatch() {
	// Some API platform, we register a handler using their patterns.
	RegisterApiHandler(func(ctx context.Context, request *MyApiRequest) (*MyApiResponse, error) {
		// Do any validation of the request, check parameters, etc.
		// This step should be _pure_, and not access mutable state.
		if err := ValidateRequest(request); err != nil {
			return nil, err
		}

		// Convert the request into an event.
		// This step is optional if the request can be used directly as the event
		event := RequestToEvent(request)

		// Dispatch the event.
		// This call blocks until the event has been accepted by the event consumer,
		// after this call returns any shared state will reflect changes made by the event.
		done, err := dispatcher.Dispatch(ctx, time.Now(), event)
		if err != nil {
			// Errors can happen if the event pipeline is closed, or a context closes, or processing of events failed.
			return nil, err
		}

		// Make sure we notify the event consumer that we are done using the global state they set up using our event.
		defer done()

		// At this point, our event has been processed along with any other events that may have been dispatched.
		// We should inspect the shared state to see if our request was successful or not.
		return ComputeResponseMessage(request)
	})
}

func Example_capture() {
	// Typically a shared timeline, a collector of events and the input to our state calculations
	tl := ApplicationTimeline

	// Frame processes all inputs since the last frame
	frame := func(ctx context.Context, t time.Time) {
		// Capture returns when all unprocessed events have been added to the given timeline.
		// Capture returns a Session that allows us to signal when processing of the timeline has completed.
		session, err := capturer.Capture(ctx, tl)
		if err != nil {
			panic(fmt.Sprintf("Failed to capture input events: %v", err))
		}

		// Events have been captured into the timeline,
		// processes those events and update our state.
		err = ProcessTimeline(t, tl)

		// Now that processing has completed and any shared state updated,
		// let the event sources know that they can compute their responses.
		// These calls will block until all event sources that contributed to the
		// timeline have called done().
		if err != nil {
			session.Reject(err)
		} else {
			session.Commit()
		}

		// Processing is complete, we are free to wait for more events and begin another frame
	}

	// Run our frame in a loop (60 fps), processing events, calculating state, computing responses
	RunAt(time.Second/60, frame)
}

// helper vars and functions for the examples above

var dispatcher Dispatcher
var capturer Capturer
var ApplicationTimeline timeline.AddTL

type MyApiRequest struct{}
type MyApiResponse struct{}
type MyEvent struct{}

func ValidateRequest(_ *MyApiRequest) error {
	return nil
}
func RequestToEvent(_ *MyApiRequest) *MyEvent {
	return &MyEvent{}
}
func RegisterApiHandler(_ interface{}) {}
func ComputeResponseMessage(_ *MyApiRequest) (*MyApiResponse, error) {
	return nil, nil
}
func ProcessTimeline(_ time.Time, _ timeline.TL) error {
	return nil
}
func RunAt(_ time.Duration, _ func(ctx context.Context, t time.Time)) {}
