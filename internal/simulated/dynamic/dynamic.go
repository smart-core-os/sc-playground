package dynamic

// Completion provides a way to wait for an action to complete, by waiting on Done.
// Calls that start an action in the background can return an implementation of Completion
// to allow the caller to tell when the action is completed.
// Completion does not provide any way of conveying return or error values; use a different abstraction
// if these are required.
type Completion interface {
	// Done returns a channel which will be closed when the action completes.
	Done() <-chan struct{}
}
