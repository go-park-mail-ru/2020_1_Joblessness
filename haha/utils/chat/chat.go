package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
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

func (r *RoomInstance) SendGeneratedMessage(message *Message) {
	if err := r.messenger.SaveMessage(message); err == nil {
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
	c.Socket.Close()
}

func (c *Chatter) Write() {
	for msg := range c.Send {
		golog.Errorf("Send by %d: %s", c.ID, msg)
		if err := c.Socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.Socket.Close()
}
