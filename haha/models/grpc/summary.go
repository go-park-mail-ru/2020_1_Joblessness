package grpcModels

import (
	"github.com/golang/protobuf/ptypes"
	baseModels "joblessness/haha/models/base"
	interviewRpc "joblessness/interviewService/rpc"
)

func TransformSendSummaryRPC(s *baseModels.SendSummary) *interviewRpc.SendSummary {
	if s == nil {
		return nil
	}

	interviewDate, _ := ptypes.TimestampProto(s.InterviewDate)

	res := &interviewRpc.SendSummary{
		VacancyID:      s.VacancyID,
		SummaryID:      s.SummaryID,
		UserID:         s.UserID,
		OrganizationID: s.OrganizationID,
		InterviewDate:  interviewDate,
		Accepted:       s.Accepted,
		Denied:         s.Denied,
	}
	return res
}

func TransformSendSummaryBase(s *interviewRpc.SendSummary) *baseModels.SendSummary {
	if s == nil {
		return nil
	}

	interviewDate, _ := ptypes.Timestamp(s.InterviewDate)

	res := &baseModels.SendSummary{
		VacancyID:      s.VacancyID,
		SummaryID:      s.SummaryID,
		UserID:         s.UserID,
		OrganizationID: s.OrganizationID,
		InterviewDate:  interviewDate,
		Accepted:       s.Accepted,
		Denied:         s.Denied,
	}
	return res
}