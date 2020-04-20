package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
)

type RoomInstance struct {
	forwardChan chan []byte
	joinChan    chan *Chatter
	leaveChan   chan *Chatter
	Chatters    map[uint64]*Chatter
	messenger   Messenger
}

func NewRoom(messenger Messenger) *RoomInstance {
	return &RoomInstance{
		forwardChan: make(chan []byte),
		joinChan:    make(chan *Chatter),
		leaveChan:   make(chan *Chatter),
		Chatters:    make(map[uint64]*Chatter),
		messenger:   messenger,
	}
}

func (r *RoomInstance) Forward(input []byte) {
	r.forwardChan <- input
}

func (r *RoomInstance) Join(character *Chatter) {
	r.joinChan <- character
}

func (r *RoomInstance) Leave(character *Chatter) {
	r.leaveChan <- character
}

func (r *RoomInstance) Run() {
	golog.Info("running chat room")
	for {
		select {
		case chatter := <-r.joinChan:
			golog.Infof("new chatter in room")
			r.Chatters[chatter.ID] = chatter
		case chatter := <-r.leaveChan:
			golog.Infof("chatter leaving room")
			delete(r.Chatters, chatter.ID)
			close(chatter.Send)
		case rawMessage := <-r.forwardChan:
			r.HandleMessage(rawMessage)
		}
	}
}

func (r *RoomInstance) SendGeneratedMessage(message *Message) {
	if err := r.messenger.SaveMessage(message); err == nil {
		receiver, existReceiver := r.Chatters[message.UserTwoId]
		if existReceiver {
			rawMessage, _ := json.Marshal(message)

			select {
			case receiver.Send <- rawMessage:
			default:
				delete(r.Chatters, receiver.ID)
				close(receiver.Send)
			}
		}
	}
}

func (r *RoomInstance) HandleMessage(rawMessage []byte) {
	var message *Message
	json.Unmarshal(rawMessage, &message)

	golog.Infof("chatter '%v' writing message to room, message: %v", message.UserOne, message.Message)

	if err := r.messenger.SaveMessage(message); err == nil {
		receiver, existReceiver := r.Chatters[message.UserTwoId]
		if existReceiver {
			select {
			case receiver.Send <- rawMessage:
			default:
				delete(r.Chatters, receiver.ID)
				close(receiver.Send)
			}
		}
		author, existAuthor := r.Chatters[message.UserOneId]
		if existAuthor {
			select {
			case author.Send <- rawMessage:
			default:
				delete(r.Chatters, author.ID)
				close(author.Send)
			}
		}
	}
}

type Chatter struct {
	ID     uint64
	Socket *websocket.Conn
	Send   chan []byte
	Room   Room
}

func (c *Chatter) Read() {
	for {
		if _, msg, err := c.Socket.ReadMessage(); err == nil {
			c.Room.Forward(msg)
		} else {
			break
		}
	}
	c.Socket.Close()
}

func (c *Chatter) Write() {
	for msg := range c.Send {
		if err := c.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.Socket.Close()
}