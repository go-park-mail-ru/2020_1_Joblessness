package grpcModels

import (
	"github.com/golang/protobuf/ptypes"
	baseModels "joblessness/haha/models/base"
	"joblessness/haha/utils/chat"
	interviewRpc "joblessness/interviewService/rpc"
)

func TransformMessageRPC(m *chat.Message) *interviewRpc.Message {
	if m == nil {
		return nil
	}

	created, _ := ptypes.TimestampProto(m.Created)

	res := &interviewRpc.Message{
		Message:   m.Message,
		UserOneId: m.UserOneId,
		UserOne:   m.UserOne,
		UserTwoId: m.UserTwoId,
		UserTwo:   m.UserTwo,
		Created:   created,
	}
	return res
}

func TransformMessageBase(m *interviewRpc.Message) *chat.Message {
	if m == nil {
		return nil
	}

	created, _ := ptypes.Timestamp(m.Created)

	res := &chat.Message{
		Message:   m.Message,
		UserOneId: m.UserOneId,
		UserOne:   m.UserOne,
		UserTwoId: m.UserTwoId,
		UserTwo:   m.UserTwo,
		Created:   created,
	}
	return res
}

func TransformChatParamsRPC(c *baseModels.ChatParameters) *interviewRpc.ChatParameters {
	if c == nil {
		return nil
	}

	res := &interviewRpc.ChatParameters{
		From: c.From,
		To:   c.To,
		Page: c.Page,
	}
	return res
}

func TransformChatParamsBase(c *interviewRpc.ChatParameters) *baseModels.ChatParameters {
	if c == nil {
		return nil
	}

	res := &baseModels.ChatParameters{
		From: c.From,
		To:   c.To,
		Page: c.Page,
	}
	return res
}

func TransformMessagesRPC(m *baseModels.Messages) *interviewRpc.Messages {
	if m == nil {
		return nil
	}

	from := make([]*interviewRpc.Message, len(m.From))
	for i, mes := range m.From {
		created, _ := ptypes.TimestampProto(mes.Created)

		from[i] = &interviewRpc.Message{
			Message:   mes.Message,
			UserOneId: mes.UserOneId,
			UserOne:   mes.UserOne,
			UserTwoId: mes.UserTwoId,
			UserTwo:   mes.UserTwo,
			Created:   created,
		}
	}

	to := make([]*interviewRpc.Message, len(m.To))
	for i, mes := range m.To {
		created, _ := ptypes.TimestampProto(mes.Created)

		to[i] = &interviewRpc.Message{
			Message:   mes.Message,
			UserOneId: mes.UserOneId,
			UserOne:   mes.UserOne,
			UserTwoId: mes.UserTwoId,
			UserTwo:   mes.UserTwo,
			Created:   created,
		}
	}

	return &interviewRpc.Messages{
		From: from,
		To:   to,
	}
}

func TransformMessagesBase(m *interviewRpc.Messages) *baseModels.Messages {
	if m == nil {
		return nil
	}

	from := make([]*chat.Message, len(m.From))
	for i, mes := range m.From {
		created, _ := ptypes.Timestamp(mes.Created)

		from[i] = &chat.Message{
			Message:   mes.Message,
			UserOneId: mes.UserOneId,
			UserOne:   mes.UserOne,
			UserTwoId: mes.UserTwoId,
			UserTwo:   mes.UserTwo,
			Created:   created,
		}
	}

	to := make([]*chat.Message, len(m.To))
	for i, mes := range m.To {
		created, _ := ptypes.Timestamp(mes.Created)

		to[i] = &chat.Message{
			Message:   mes.Message,
			UserOneId: mes.UserOneId,
			UserOne:   mes.UserOne,
			UserTwoId: mes.UserTwoId,
			UserTwo:   mes.UserTwo,
			Created:   created,
		}
	}

	return &baseModels.Messages{
		From: from,
		To:   to,
	}
}

func TransformSummaryCredentialsRPC(s *baseModels.SummaryCredentials) *interviewRpc.SummaryCredentials {
	if s == nil {
		return nil
	}

	res := &interviewRpc.SummaryCredentials{
		UserID:           s.UserID,
		OrganizationID:   s.OrganizationID,
		UserName:         s.UserName,
		OrganizationName: s.OrganizationName,
	}
	return res
}

func TransformSummaryCredentialsBase(s *interviewRpc.SummaryCredentials) *baseModels.SummaryCredentials {
	if s == nil {
		return nil
	}

	res := &baseModels.SummaryCredentials{
		UserID:           s.UserID,
		OrganizationID:   s.OrganizationID,
		UserName:         s.UserName,
		OrganizationName: s.OrganizationName,
	}
	return res
}

func TransformConversationsRPC(c baseModels.Conversations) *interviewRpc.Conversations {
	if c == nil {
		return nil
	}

	res := make([]*interviewRpc.ConversationTitle, len(c))
	for i, title := range c {
		interviewDate, _ := ptypes.TimestampProto(title.InterviewDate)

		res[i] = &interviewRpc.ConversationTitle{
			ChatterId:     title.ChatterID,
			Avatar: title.Avatar,
			Tag:   title.Tag,
			ChatterName:   title.ChatterName,
			InterviewDate: interviewDate,
		}
	}

	return &interviewRpc.Conversations{Title: res}
}

func TransformConversationsBase(c *interviewRpc.Conversations) baseModels.Conversations {
	if c == nil {
		return nil
	}

	res := make(baseModels.Conversations, len(c.Title))
	for i, title := range c.Title {
		interviewDate, _ := ptypes.Timestamp(title.InterviewDate)

		res[i] = &baseModels.ConversationTitle{
			ChatterID:     title.ChatterId,
			Avatar: title.Avatar,
			Tag:           title.Tag,
			ChatterName:   title.ChatterName,
			InterviewDate: interviewDate,
		}
	}

	return res
}
