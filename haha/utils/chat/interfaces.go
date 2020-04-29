package chat

//go:generate mockgen -destination=./mock/room.go -package=mock joblessness/haha/utils/chat Room

type Messenger interface {
	SaveMessage(message *Message) (err error)
}

type Room interface {
	Run()
	SendGeneratedMessage(message *Message)
	HandleMessage(rawMessage []byte)
	Forward([]byte)
	Join(character *Chatter)
	Leave(character *Chatter)
}
