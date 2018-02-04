package pubsub

type IPublisher interface {
	Publish(msg IMessage) error
}
