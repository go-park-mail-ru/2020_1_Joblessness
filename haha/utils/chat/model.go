package chat

import (
	"time"
)

type Message struct {
	Message string `json:"message"`
	UserOneId  uint64 `json:"userOneId"`
	UserOne  string `json:"userOne"`
	UserTwoId  uint64 `json:"userTwoId"`
	UserTwo  string `json:"userTwo"`
	Created time.Time `json:"created"`
}