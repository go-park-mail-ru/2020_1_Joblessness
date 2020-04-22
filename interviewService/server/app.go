package interviewServer

//go:generate cd ./interviewService/rpc && protoc -I=. interview.proto --go_out=plugins=grpc:.

import (
	"golang.org/x/net/context"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	grpcModels "joblessness/haha/models/grpc"
	interviewRpc "joblessness/interviewService/rpc"
)

type server struct {
	interviewRepo interviewInterfaces.InterviewRepository
}

func (s *server) IsOrganizationVacancy(cc context.Context, in *interviewRpc.IsParameters) (*interviewRpc.Status, error) {
	err := s.interviewRepo.IsOrganizationVacancy(in.VacancyID, in.UserID)
	if err != nil {
		return &interviewRpc.Status{Code: 500}, err
	}

	return &interviewRpc.Status{Code: 200}, nil
}

func (s *server) ResponseSummary(cc context.Context, in *interviewRpc.SendSummary) (*interviewRpc.Status, error) {
	sendSummary := grpcModels.TransformSendSummaryBase(in)

	err := s.interviewRepo.ResponseSummary(sendSummary)
	if err != nil {
		return &interviewRpc.Status{Code: 500}, err
	}

	return &interviewRpc.Status{Code: 200}, nil
}

func (s *server) SaveMessage(cc context.Context, in *interviewRpc.Message) (*interviewRpc.Status, error) {
	message := grpcModels.TransformMessageBase(in)

	err := s.interviewRepo.SaveMessage(message)
	if err != nil {
		return &interviewRpc.Status{Code: 500}, err
	}

	return &interviewRpc.Status{Code: 200}, nil
}

func (s *server) GetHistory(cc context.Context, in *interviewRpc.ChatParameters) (*interviewRpc.Messages, error) {
	chatParams := grpcModels.TransformChatParamsBase(in)

	res, err := s.interviewRepo.GetHistory(chatParams)
	if err != nil {
		return nil, err
	}

	return grpcModels.TransformMessagesRPC(&res), nil
}

func (s *server) GetResponseCredentials(cc context.Context, in *interviewRpc.CredentialsParams) (*interviewRpc.SummaryCredentials, error) {
	res, err := s.interviewRepo.GetResponseCredentials(in.SummaryID, in.VacancyID)
	if err != nil {
		return nil, err
	}

	return grpcModels.TransformSummaryCredentialsRPC(res), nil
}

func (s *server) GetConversations(cc context.Context, in *interviewRpc.CredentialsParams) (*interviewRpc.Conversations, error) {
	res, err := s.interviewRepo.GetConversations(in.UserID)
	if err != nil {
		return nil, err
	}

	return grpcModels.TransformConversationsRPC(res), nil
}

func NewInterviewServer(u interviewInterfaces.InterviewRepository) interviewRpc.InterviewServer {
	return &server{interviewRepo: u}
}