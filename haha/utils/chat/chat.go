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
			golog.Infof("new chatter in room: %d", chatter.ID)
			r.Chatters[chatter.ID] = chatter
		case chatter := <-r.leaveChan:
			golog.Infof("chatter leaving room: %d", chatter.ID)
			delete(r.Chatters, chatter.ID)
			close(chatter.Send)
		case rawMessage := <-r.forwardChan:
			r.HandleMessage(rawMessage)
		}
	}
}

func (r *RoomInstance) SendGeneratedMessage(message *Message) error {
	err := r.messenger.SaveMessage(message)
	if err == nil {
		receiver, existReceiver := r.Chatters[message.UserTwoId]
		if existReceiver {
			rawMessage, err := json.Marshal(message)
			if err == nil {
				select {
				case receiver.Send <- rawMessage:
				default:
					delete(r.Chatters, receiver.ID)
					close(receiver.Send)
				}
			} else {
				golog.Errorf("Broken message: %+v", message)
				return err
			}
		}
	}
	return err
}

func (r *RoomInstance) HandleMessage(rawMessage []byte) {
	var message *Message
	err := json.Unmarshal(rawMessage, &message)
	if err != nil {
		golog.Infof("broken message received: %v", err)
	}
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
				close(author.Send)
				delete(r.Chatters, author.ID)
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
			golog.Errorf("Received by %d: %s", c.ID, msg)
			if len(msg) != 0 {
				c.Room.Forward(msg)
			} else {
				golog.Errorf("Received empty array by %d", c.ID)
			}
		} else {
			break
		}
	}

	err := c.Socket.Close()
	if err != nil {
		golog.Error("Socket closed with error: ", err)
	}
}

func (c *Chatter) Write() {
	for msg := range c.Send {
		golog.Errorf("Send by %d: %s", c.ID, msg)
		if err := c.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}

	err := c.Socket.Close()
	if err != nil {
		golog.Error("Socket closed with error: ", err)
	}
}
