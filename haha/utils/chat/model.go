package chat

import (
	"time"
)

type Message struct {
	Message   string    `json:"message"`
	UserOneID uint64    `json:"userOneId,omitempty"`
	UserOne   string    `json:"userOne,omitempty"`
	UserTwoID uint64    `json:"userTwoId,omitempty"`
	UserTwo   string    `json:"userTwo,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	VacancyID uint64    `json:"vacancyId,omitempty"`
}
