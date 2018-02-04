package pubsub

import "time"

type BusEndpoint struct {
	bus *Bus
	ChannelId
	channel chan IMessage
	router  IMessageRouter
	ids     []MessageId
}

//  Get the channel id
func (endpoint *BusEndpoint) Id() ChannelId {
	return endpoint.ChannelId
}

//  Un-subscribe and shutdown all the endpoint
func (endpoint *BusEndpoint) Close() {
	for _, id := range endpoint.ids {
		if err := endpoint.bus.unregisterMessageInterest(id, endpoint); err != nil {
			panic(err)
		}
	}
	endpoint.Shutdown()
}

//  Shutdown the router and the channel
func (endpoint *BusEndpoint) Shutdown() {
	endpoint.router.Close()
	close(endpoint.channel)
}

//  Enqueue the message for later dispatch
func (endpoint *BusEndpoint) Enqueue(msg IMessage) {
	endpoint.channel <- msg
}

func (endpoint *BusEndpoint) Dispatch() {
	msg := <-endpoint.channel
	endpoint.router.Dispatch(msg)
}

func (endpoint *BusEndpoint) TimedDispatch(timeout time.Duration) (success bool) {
	success = false
	select {
	case msg, ok := <-endpoint.channel:
		if ok {
			endpoint.router.Dispatch(msg)
			success = true
		}

	case <-time.After(timeout):
	}
	return
}

func (endpoint *BusEndpoint) TryDispatch() (success bool) {
	success = false
	select {
	case msg, ok := <-endpoint.channel:
		if ok {
			endpoint.router.Dispatch(msg)
			success = true
		}
	default:
	}
	return
}

func (endpoint *BusEndpoint) Publish(msg IMessage) error {
	return endpoint.bus.Publish(msg)
}
