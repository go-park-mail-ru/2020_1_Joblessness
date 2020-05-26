package interviewGrpc

import (
	"context"
	"google.golang.org/grpc"
	interviewInterfaces "joblessness/haha/interview/interfaces"
	baseModels "joblessness/haha/models/base"
	grpcModels "joblessness/haha/models/grpc"
	"joblessness/haha/utils/chat"
	interviewRpc "joblessness/interviewService/rpc"
)

type Repository struct {
	handler interviewRpc.InterviewClient
}

func NewInterviewGrpcRepository(conn *grpc.ClientConn) *Repository {
	return &Repository{
		handler: interviewRpc.NewInterviewClient(conn),
	}
}

func (r *Repository) IsOrganizationVacancy(vacancyID, userID uint64) (err error) {
	status, err := r.handler.IsOrganizationVacancy(context.Background(), &interviewRpc.IsParameters{UserID: userID, VacancyID: vacancyID})
	if err != nil {
		return err
	}

	switch status.Code {
	case 403:
		return interviewInterfaces.ErrOrganizationIsNotOwner
	case 200:
		return nil
	default:
		return err
	}
}

func (r *Repository) ResponseSummary(sendSummary *baseModels.SendSummary) (err error) {
	status, err := r.handler.ResponseSummary(context.Background(), grpcModels.TransformSendSummaryRPC(sendSummary))
	if err != nil {
		return err
	}

	switch status.Code {
	case 404:
		return interviewInterfaces.ErrNoSummaryToRefresh
	case 200:
		return nil
	default:
		return err
	}
}

func (r *Repository) SaveMessage(message *chat.Message) (err error) {
	_, err = r.handler.SaveMessage(context.Background(), grpcModels.TransformMessageRPC(message))

	return err
}

func (r *Repository) GetHistory(parameters *baseModels.ChatParameters) (result baseModels.Messages, err error) {
	res, err := r.handler.GetHistory(context.Background(), grpcModels.TransformChatParamsRPC(parameters))

	return *grpcModels.TransformMessagesBase(res), err
}

func (r *Repository) GetResponseCredentials(summaryID, vacancyID uint64) (result *baseModels.SummaryCredentials, err error) {
	res, err := r.handler.GetResponseCredentials(context.Background(), &interviewRpc.CredentialsParams{
		SummaryID: summaryID,
		VacancyID: vacancyID,
	})

	return grpcModels.TransformSummaryCredentialsBase(res), err
}

func (r *Repository) GetConversations(userID uint64) (result baseModels.Conversations, err error) {
	res, err := r.handler.GetConversations(context.Background(), &interviewRpc.CredentialsParams{
		UserID: userID,
	})

	return grpcModels.TransformConversationsBase(res), err
}
