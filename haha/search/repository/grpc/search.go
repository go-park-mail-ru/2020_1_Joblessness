package searchGrpc

import (
	"context"
	"google.golang.org/grpc"
	"io"
	baseModels "joblessness/haha/models/base"
	grpcModels "joblessness/haha/models/grpc"
	"joblessness/searchService/rpc"
)

type SearchGrpcRepository struct {
	handler searchRpc.SearchClient
}

func NewSearchGrpcRepository(conn *grpc.ClientConn) *SearchGrpcRepository {
	return &SearchGrpcRepository{
		handler: searchRpc.NewSearchClient(conn),
	}
}

func (r *SearchGrpcRepository) SearchPersons(params *baseModels.SearchParams) (result []*baseModels.Person, err error) {
	res, err := r.handler.SearchPersons(context.Background(), grpcModels.TransformParamsRPC(params))
	if err != nil {
		return result, err
	}

	result = make([]*baseModels.Person, 0)
	for {
		var person *searchRpc.Person
		person, err = res.Recv()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
		result = append(result, grpcModels.TransformPersonBase(person))
	}

	return result, err
}

func (r *SearchGrpcRepository) SearchOrganizations(params *baseModels.SearchParams) (result []*baseModels.Organization, err error) {
	res, err := r.handler.SearchOrganizations(context.Background(), grpcModels.TransformParamsRPC(params))
	if err != nil {
		return result, err
	}

	result = make([]*baseModels.Organization, 0)
	for {
		var org *searchRpc.Organization
		org, err = res.Recv()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
		result = append(result, grpcModels.TransformOrganizationBase(org))
	}

	return result, err
}

func (r *SearchGrpcRepository) SearchVacancies(params *baseModels.SearchParams) (result []*baseModels.Vacancy, err error) {
	res, err := r.handler.SearchVacancies(context.Background(), grpcModels.TransformParamsRPC(params))
	if err != nil {
		return result, err
	}

	result = make([]*baseModels.Vacancy, 0)
	for {
		var vac *searchRpc.Vacancy
		vac, err = res.Recv()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
		result = append(result, grpcModels.TransformVacancyBase(vac))
	}

	return result, err
}
