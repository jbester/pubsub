package pubsub

//  Identifiable message queue interface
type IMessageQueue interface {
	IIdentifiableEndpoint
	Shutdown()
	Enqueue(msg IMessage)
}
