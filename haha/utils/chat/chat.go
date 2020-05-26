package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kataras/golog"
	"time"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type RoomInstance struct {
	forwardChan chan []byte
	joinChan    chan *Chatter
	leaveChan   chan *Chatter
	Chatters    map[uint64]map[*websocket.Conn]*Chatter
	messenger   Messenger
}

func NewRoom(messenger Messenger) *RoomInstance {
	return &RoomInstance{
		forwardChan: make(chan []byte),
		joinChan:    make(chan *Chatter),
		leaveChan:   make(chan *Chatter),
		Chatters:    make(map[uint64]map[*websocket.Conn]*Chatter),
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
			if r.Chatters[chatter.ID] == nil {
				r.Chatters[chatter.ID] = make(map[*websocket.Conn]*Chatter)
			}
			r.Chatters[chatter.ID][chatter.Socket] = chatter
		case chatter := <-r.leaveChan:
			golog.Infof("chatter leaving room: %d", chatter.ID)
			delete(r.Chatters[chatter.ID], chatter.Socket)
			if len(r.Chatters[chatter.ID]) == 0 {
				delete(r.Chatters, chatter.ID)
			}
			close(chatter.Send)
		case rawMessage := <-r.forwardChan:
			r.HandleMessage(rawMessage)
		}
	}
}

func (r *RoomInstance) SendGeneratedMessage(message *Message) error {
	err := r.messenger.SaveMessage(message)
	if err == nil {
		receivers, existReceivers := r.Chatters[message.UserTwoID]
		if existReceivers {
			rawMessage, err := json.Marshal(message)
			if err == nil {
				for _, receiver := range receivers{
					select {
					case receiver.Send <- rawMessage:
					default:
						delete(r.Chatters[receiver.ID], receiver.Socket)
						if len(r.Chatters[receiver.ID]) == 0 {
							delete(r.Chatters, receiver.ID)
						}
						close(receiver.Send)
					}
				}
			} else {
				golog.Errorf("Broken message: %+v", message)
				return err
			}
		}
	}
	golog.Errorf("Broken message generated: %+v", message)
	return err
}

func (r *RoomInstance) HandleMessage(rawMessage []byte) {
	var message *Message
	err := json.Unmarshal(rawMessage, &message)
	if err != nil {
		golog.Infof("Broken message received: %v", err)
	}
	golog.Infof("Chatter %d writing message to %d, message: %v", message.UserOneID, message.UserTwoID, message.Message)

	if err := r.messenger.SaveMessage(message); err == nil {
		receivers, existReceivers := r.Chatters[message.UserTwoID]
		if existReceivers {
			for _, receiver := range receivers{
				select {
				case receiver.Send <- rawMessage:
				default:
					delete(r.Chatters[receiver.ID], receiver.Socket)
					if len(r.Chatters[receiver.ID]) == 0 {
						delete(r.Chatters, receiver.ID)
					}
					close(receiver.Send)
				}
			}
		} else {
			golog.Infof("Receiver does not connected, message: %v", message.Message)
		}

		authors, existAuthors := r.Chatters[message.UserOneID]
		if existAuthors {
			for _, author := range authors {
				select {
				case author.Send <- rawMessage:
				default:
					delete(r.Chatters[author.ID], author.Socket)
					if len(r.Chatters[author.ID]) == 0 {
						delete(r.Chatters, author.ID)
					}
					close(author.Send)
				}
			}
		} else {
			golog.Infof("Author does not connected, message: %v", message.Message)
		}
	} else {
		golog.Errorf("Messages was not saved: %+v", message)
	}
}

type Chatter struct {
	ID     uint64
	Socket *websocket.Conn
	Send   chan []byte
	Room   Room
}

func (c *Chatter) Read() {
	var err error

	err = c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		golog.Error("Cannot set read deadline: ", err)
	}
	c.Socket.SetPongHandler(func(string) error {
		err = c.Socket.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			golog.Error("Cannot set read deadline: ", err)
		}
		return err
	})

	for {
		if _, msg, err := c.Socket.ReadMessage(); err == nil {
			golog.Infof("Read by %d: %s", c.ID, msg)
			if len(msg) != 0 {
				c.Room.Forward(msg)
			} else {
				golog.Infof("Read empty array by %d", c.ID)
			}
		} else {
			break
		}
	}
	golog.Error("Read from socket terminated: ", err)

	err = c.Socket.Close()
	if err != nil {
		golog.Error("Socket closed with error: ", err)
	}
}

func (c *Chatter) Write() {
	var err error
	ticker := time.NewTicker(pingPeriod)

	LOOP: for {
		select {
		case message, ok := <-c.Send:
			err = c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				golog.Error("Cannot set write deadline: ", err)
			}

			if !ok {
				err = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				break LOOP
			}

			golog.Errorf("Write by %d: %s", c.ID, message)
			err = c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				golog.Error("Cannot set write deadline: ", err)
			}

			if err = c.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
				break LOOP
			}
		case <-ticker.C:
			err = c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				golog.Error("Cannot set write deadline: ", err)
			}

			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				break LOOP
			}
		}
	}
	golog.Error("Write to socket terminated: ", err)

	ticker.Stop()
	err = c.Socket.Close()
	if err != nil {
		golog.Error("Socket closed with error: ", err)
	}
}