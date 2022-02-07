// Package input provides tools for collecting and processing user input, typically API calls.
// The package defines the Dispatcher and Capturer types for collecting and processing these events.
//
// A source of events, say an API handler, would dispatch an event to a Dispatcher when an API request comes in.
// At some later point, an event processor captures all events via a Capturer.
// The event processor can notify the event source when it has completed processing the event, and the event source
// can notify the processor when it has finished reading any changes applied by this event.
//
// Most applications of this package will use a Queue as the underlying mediator, passing it as a Dispatcher to event
// sources, and as a Capturer to event processors.
package input
