package pubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MessageA struct {
	value int
}

func (msg MessageA) Id() MessageId {
	return 1
}

type MessageB struct {
	value int
}

func (msg MessageB) Id() MessageId {
	return 2
}

//  verify single endpoint single message works
func TestSingleSubscriber(t *testing.T) {
	bus := NewBus()
	defer bus.Close()
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
		t.Fatalf("Callback not executed")
	}
}

//  verify multiple subscribers work
func TestMultipleSubscribers(t *testing.T) {
	bus := NewBus()
	defer bus.Close()
	sub1 := bus.NewCallbackEndpoint(10)
	sub2 := bus.NewCallbackEndpoint(1)

	received1 := make(chan bool)
	received2 := make(chan bool)

	sub1.SubscribeMessage(1, func(msg MessageA) {
		received1 <- true
	})

	sub2.SubscribeMessage(1, func(msg MessageA) {
		received2 <- true
	})

	err := bus.Publish(MessageA{})

	assert.NoError(t, err)
	assert.True(t, sub1.TimedDispatch(time.Millisecond*50), "Message not received at endpoint level")
	assert.True(t, sub2.TimedDispatch(time.Millisecond*50), "Message not received at endpoint level")

	var ok1, ok2 bool
	for !ok1 || !ok2 {
		select {
		case <-received1:
			ok1 = true
		case <-received2:
			ok2 = true
		case <-time.After(time.Millisecond * 50):
			assert.Fail(t, "Callback not executed subscriber1=", ok1, " subscriber2=", ok2)
		}
	}
}

//  verify multiple publishers work
func TestMultiplePublishers(t *testing.T) {
	bus := NewBus()
	defer bus.Close()
	sub := bus.NewCallbackEndpoint(10)
	result := make(chan int)

	sub.SubscribeMessage(1, func(msg MessageA) {
		result <- msg.value
	})

	go func() {
		err := bus.Publish(MessageA{1})
		assert.NoError(t, err)
	}()

	go func() {
		err := bus.Publish(MessageA{2})
		assert.NoError(t, err)
	}()

	total := 0
	for i := 0; i < 2; i++ {
		assert.True(t, sub.TimedDispatch(time.Millisecond*50), "message not received at endpoint level")
		select {
		case v := <-result:
			total += v
		case <-time.After(time.Millisecond * 50):
			assert.Fail(t, "callback not executed")
		}
	}

	assert.Equal(t, 3, total, "both messages not received properly")
}

//  verify close removes subscriptions
func TestSubClose(t *testing.T) {
	bus := NewBus()
	defer bus.Close()
	sub := bus.NewCallbackEndpoint(10)
	result := make(chan bool)

	sub.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	sub.Close()

	err := bus.Publish(MessageA{})

	if err == nil {
		t.Fatalf("expected no subscribers")
	}

}

//  verify close leaves other subscribers unaffected
func TestSubCloseMultiple(t *testing.T) {
	bus := NewBus()
	defer bus.Close()
	sub1 := bus.NewCallbackEndpoint(10)
	sub2 := bus.NewCallbackEndpoint(10)
	result := make(chan bool)

	sub1.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	sub2.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	sub1.Close()

	err := bus.Publish(MessageA{})

	// check that sub 2 still works

	assert.NoError(t, err, "error publishing")
	assert.True(t, sub2.TimedDispatch(time.Millisecond*50), "message not received at endpoint level")

	select {
	case <-result:
	case <-time.After(time.Millisecond * 50):
		assert.Fail(t, "callback not executed")
	}

}

//  verify nil does cause issues
func TestPublishNil(t *testing.T) {
	bus := NewBus()
	defer bus.Close()
	sub := bus.NewCallbackEndpoint(10)
	result := make(chan bool)

	sub.SubscribeMessage(1, func(msg MessageA) {
		result <- true
	})

	err := bus.Publish(nil)

	assert.Error(t, err)
}

//  verify multiple callbacks can be registered for a message
func TestMultipleSubscribe(t *testing.T) {
	bus := NewBus()
	defer bus.Close()
	sub1 := bus.NewCallbackEndpoint(10)
	count := 0

	sub1.SubscribeMessage(1, func(msg MessageA) {
		count++
	})
	sub1.SubscribeMessage(1, func(msg MessageA) {
		count += 10
	})

	err := bus.Publish(MessageA{})

	// check that sub 2 still works
	assert.NoError(t, err)

	for {
		if !sub1.TimedDispatch(time.Millisecond * 50) {
			break
		}
	}

	<-time.After(time.Millisecond * 50)
	assert.Equal(t, 11, count)
}
