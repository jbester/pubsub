package pubsub

type ChannelId uint64

//  Interface allows channels to be uniquely identified
type IIdentifiableEndpoint interface {
	//  Get the unique channel id
	Id() ChannelId
}

