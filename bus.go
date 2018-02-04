package pubsub

import (
	"sync"

	"github.com/pkg/errors"
)

type MessageId uint64
type subscriberMap map[MessageId][]IMessageQueue

type Bus struct {
	subscribers       subscriberMap
	lastUsedChannelId ChannelId
	commands          chan func()
}

var ErrPreviouslySubscribed = errors.New("previously subscribed")
var ErrNotRegistered = errors.New("not registered")
var ErrInvalidArgument = errors.New("invalid argument")
var ErrNoSubscribers = errors.New("no subscribers")

func NewBus() *Bus {
	bus := &Bus{
		subscribers:       make(subscriberMap, 0),
		lastUsedChannelId: 0,
		commands:          make(chan func()),
	}
	// execute the commands in a single context
	go func() {
		defer close(bus.commands)
		for {
			fn := <-bus.commands
			if fn == nil {
				break
			}
			fn()
		}
	}()
	return bus
}

func (bus *Bus) Close() {
	subs := make(map[IMessageQueue]bool)
	// get a canonical list of all subscribers
	for _, set := range bus.subscribers {
		for _, sub := range set {
			subs[sub] = true
		}
	}
	// close them all
	for sub, _ := range subs {
		sub.Shutdown()
	}
	bus.commands <- nil
}

func (bus *Bus) registerMessageInterest(typeId MessageId, subscriber *CallbackEndpoint) error {
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)

	bus.commands <- func() {
		subscribers := bus.subscribers[typeId]
		for _, candidate := range subscribers {
			if candidate.Id() == subscriber.Id() {
				err = ErrPreviouslySubscribed
				wg.Done()
				return
			}
		}
		subscribers = append(subscribers, subscriber)
		bus.subscribers[typeId] = subscribers
		err = nil
		wg.Done()
	}

	wg.Wait()
	return err
}

func (bus *Bus) unregisterMessageInterest(typeId MessageId, subscriber *BusEndpoint) error {
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)

	bus.commands <- func() {
		subscribers := bus.subscribers[typeId]
		for i, candidate := range subscribers {
			if candidate.Id() == subscriber.Id() {
				bus.subscribers[typeId] = append(subscribers[:i], subscribers[i+1:]...)

				wg.Done()
				err = nil
				return
			}
		}
		err = ErrNotRegistered
		wg.Done()
	}

	wg.Wait()
	return err
}

func (bus *Bus) NewCallbackEndpoint(maxQueueLength int) ICallbackEndpoint {
	wg := sync.WaitGroup{}
	wg.Add(1)

	id := ChannelId(0)
	bus.commands <- func() {
		bus.lastUsedChannelId++
		id = bus.lastUsedChannelId
		wg.Done()
	}
	wg.Wait()

	sub := NewBusEndpoint(bus, id, maxQueueLength)
	return sub
}

func (bus *Bus) Publish(msg IMessage) error {
	if msg == nil {
		return errors.Wrap(ErrInvalidArgument, "msg not be nil")
	}

	var err error
	var wg = sync.WaitGroup{}
	wg.Add(1)

	bus.commands <- func() {
		subscribers := bus.subscribers[msg.Id()]
		if len(subscribers) != 0 {
			for _, candidate := range subscribers {
				candidate.Enqueue(msg)
			}

			err = nil
		} else {
			err = ErrNoSubscribers
		}
		wg.Done()
	}

	wg.Wait()
	return err
}
