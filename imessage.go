package pubsub

type IMessage interface {
	Id() MessageId
}

