package pubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//  Test timed dispatch works
func TestTimeDispatch(t *testing.T) {
	bus := NewBus()
	sub := bus.NewCallbackEndpoint(10)
	result := make(chan bool)

	sub.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	err := bus.Publish(MessageA{})
	assert.NoError(t, err)
	assert.True(t, sub.TimedDispatch(time.Millisecond*50), "Message not received at endpoint level")

	select {
	case <-result:
	case <-time.After(time.Millisecond * 50):
		assert.Fail(t, "Callback not executed")
	}
}

//  Test timed dispatch times out
func TestTimeDispatchFailure(t *testing.T) {
	bus := NewBus()
	sub := bus.NewCallbackEndpoint(10)
	result := make(chan bool)

	sub.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	assert.False(t, sub.TimedDispatch(time.Millisecond*50), "expected timeout")
}

//  Test try dispatch works
func TestTryDispatchSuccess(t *testing.T) {
	bus := NewBus()
	sub := bus.NewCallbackEndpoint(10)
	result := make(chan bool)

	sub.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	err := bus.Publish(MessageA{})

	assert.NoError(t, err, "error publishing")

	<-time.After(time.Millisecond * 50)

	assert.True(t, sub.TryDispatch(), "message not received at endpoint level")

	select {
	case <-result:
	case <-time.After(time.Millisecond * 50):
		assert.Fail(t, "callback not exected")
	}
}

//  Test try dispatch times out
func TestTryDispatchFailure(t *testing.T) {
	bus := NewBus()
	sub := bus.NewCallbackEndpoint(10)
	result := make(chan bool)

	sub.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	assert.False(t, sub.TryDispatch(), "unexpected message received at endpoint level")
}

//  Test dispatch works
func TestDispatch(t *testing.T) {
	bus := NewBus()
	sub := bus.NewCallbackEndpoint(10)
	result := make(chan bool)

	sub.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	err := sub.Publish(MessageA{})
	assert.NoError(t, err, "error publishing")

	sub.Dispatch()

	select {
	case <-result:
	case <-time.After(time.Millisecond * 50):
		assert.Fail(t, "Callback not executed")
	}
}
