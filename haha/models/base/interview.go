package baseModels

import (
	"joblessness/haha/utils/chat"
	"time"
)

//easyjson:json
type Messages struct {
	From []*chat.Message `json:"from"`
	To   []*chat.Message `json:"to"`
}

type ChatParameters struct {
	From uint64
	To   uint64
	Page uint64
}

type SummaryCredentials struct {
	UserID           uint64
	OrganizationID   uint64
	UserName         string
	OrganizationName string
}

//easyjson:json
type ConversationTitle struct {
	ChatterID     uint64    `json:"chatter_id"`
	Avatar        string    `json:"avatar"`
	ChatterName   string    `json:"chatter_name"`
	Tag           string    `json:"tag"`
	InterviewDate time.Time `json:"interview_date"`
}

//easyjson:json
type Conversations []*ConversationTitle
