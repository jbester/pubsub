package pubsub

//  ICallbackEndpoint interface is a connection to a given bus.
//  It allows you to:
//    - subscribe to messages with callbacks
//    - control dispatch of incoming messages
//    - free sources
type ICallbackEndpoint interface {
	IDispatcher
	ISubscriber
	IIdentifiableEndpoint
	IPublisher

	// Close this subscribe and free any resources associated
	Close()
}


type CallbackEndpoint struct {
	BusEndpoint
}

func NewBusEndpoint(bus *Bus, id ChannelId, maxQueueLength int) *CallbackEndpoint {
	endpoint := &CallbackEndpoint{
		BusEndpoint: BusEndpoint {
			bus:       bus,
			ChannelId: id,
			channel:   make(chan IMessage, maxQueueLength),
		},
	}

	endpoint.router = NewCallbackRouter(bus, endpoint)

	return endpoint
}


// Subscribe to a message
func (endpoint *CallbackEndpoint) SubscribeMessage(typeId MessageId, fn Function) {
	endpoint.ids = append(endpoint.ids, typeId)
	endpoint.router.(ICallbackMessageRouter).BindMessage(typeId, fn)
	endpoint.bus.registerMessageInterest(typeId, endpoint)
}

