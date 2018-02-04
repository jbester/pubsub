package pubsub


//  ISuscriber interface is a broker connected to a given bus.
//  It allows you to subscribe to messages.
type ISubscriber interface {
	// Subscribe the specified callback to a specified message id
	SubscribeMessage(typeId MessageId, fn Function)
}


