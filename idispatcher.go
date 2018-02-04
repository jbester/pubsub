package pubsub

import "time"

//  CallbackRouter provides the mechanism to control flow of messages
//  to the callbacks.
type IDispatcher interface {
	// Wait for a message and dispatch it to any registered callbacks.
	Dispatch()

	// Wait up to timeout for a message and dispatch it to any registered
	// callbacks.   Returns true if a message is dispatched
	TimedDispatch(timeout time.Duration) (success bool)

	// Check to see if a message is available and dispatch it immediately.
	// Returns true if a message is dispatched
	TryDispatch() (success bool)
}

