package pubsub

import (
	"reflect"
)

type Function interface{}

type messageCallbackMap map[MessageId][]reflect.Value

type CallbackRouter struct {
	messageCallbacks messageCallbackMap
	bus              *Bus
	endpoint         *CallbackEndpoint
	commands         chan func()
}

type IMessageRouter interface {
	Close()
	Dispatch(msg IMessage)
}

type ICallbackMessageRouter interface {
	IMessageRouter
	BindMessage(typeId MessageId, fn Function)
}

func NewCallbackRouter(bus *Bus, endpoint *CallbackEndpoint) *CallbackRouter {
	dispatcher := &CallbackRouter{
		messageCallbacks: make(messageCallbackMap),
		bus: bus,
		endpoint: endpoint,
		commands: make(chan func()),
	}

	// execute the commands in a single context
	go func() {
		defer close(dispatcher.commands)
		for {
			fn := <- dispatcher.commands
			if fn == nil {
				break
			}
			fn()
		}
	}()

	return dispatcher
}

func (router *CallbackRouter) Close() {
	router.commands <- nil
}

func (router *CallbackRouter) BindMessage(typeId MessageId, fn Function) {
	router.commands <- func() {
		callbacks := router.messageCallbacks[typeId]
		callbackValue := reflect.ValueOf(fn)
		for _, candidate := range callbacks {
			if candidate == callbackValue {
				return
			}
		}
		callbacks = append(callbacks, callbackValue)
		router.messageCallbacks[typeId] = callbacks
	}
}


func (router *CallbackRouter) Dispatch(msg IMessage) {
	router.commands <- func() {
		if callbacks, ok := router.messageCallbacks[msg.Id()]; ok {
			args := []reflect.Value{reflect.ValueOf(msg)}

			for _, callback := range callbacks {
				callback.Call(args)
			}
		}
	}
}


