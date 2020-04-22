package interviewGrpc

import (
	"context"
	"google.golang.org/grpc"
	baseModels "joblessness/haha/models/base"
	grpcModels "joblessness/haha/models/grpc"
	"joblessness/haha/utils/chat"
	interviewRpc "joblessness/interviewService/rpc"
)

type InterviewGrpcRepository struct {
	handler interviewRpc.InterviewClient
}

func NewInterviewGrpcRepository(conn *grpc.ClientConn) *InterviewGrpcRepository {
	return &InterviewGrpcRepository{
		handler: interviewRpc.NewInterviewClient(conn),
	}
}

func (r *InterviewGrpcRepository) IsOrganizationVacancy(vacancyID, userID uint64) (err error) {
	_, err = r.handler.IsOrganizationVacancy(context.Background(), &interviewRpc.IsParameters{UserID: userID, VacancyID: vacancyID})

	return err
}

func (r *InterviewGrpcRepository) ResponseSummary(sendSummary *baseModels.SendSummary) (err error) {
	_, err = r.handler.ResponseSummary(context.Background(), grpcModels.TransformSendSummaryRPC(sendSummary))

	return err
}

func (r *InterviewGrpcRepository) SaveMessage(message *chat.Message) (err error) {
	_, err = r.handler.SaveMessage(context.Background(), grpcModels.TransformMessageRPC(message))

	return err
}

func (r *InterviewGrpcRepository) GetHistory(parameters *baseModels.ChatParameters) (result baseModels.Messages, err error) {
	res, err := r.handler.GetHistory(context.Background(), grpcModels.TransformChatParamsRPC(parameters))

	return *grpcModels.TransformMessagesBase(res), err
}

func (r *InterviewGrpcRepository) GetResponseCredentials(summaryID, vacancyID uint64) (result *baseModels.SummaryCredentials, err error) {
	res, err := r.handler.GetResponseCredentials(context.Background(), &interviewRpc.CredentialsParams{
		SummaryID: summaryID,
		VacancyID: vacancyID,
	})

	return grpcModels.TransformSummaryCredentialsBase(res), err
}

func (r *InterviewGrpcRepository) GetConversations(userID uint64) (result baseModels.Conversations, err error) {
	res, err := r.handler.GetConversations(context.Background(), &interviewRpc.CredentialsParams{
		UserID: userID,
	})

	return grpcModels.TransformConversationsBase(res), err
}