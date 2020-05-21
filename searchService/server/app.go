package searchServer

//go:generate cd ./searchService/rpc && protoc -I=. search.proto --go_out=plugins=rpc:.

import (
	baseModels "joblessness/haha/models/base"
	grpcModels "joblessness/haha/models/grpc"
	searchInterfaces "joblessness/haha/search/interfaces"
	"joblessness/searchService/rpc"
)

type server struct {
	searchRepo searchInterfaces.SearchRepository
}

func (s *server) SearchPersons(in *searchRpc.SearchParams, stream searchRpc.Search_SearchPersonsServer) error {
	params := &baseModels.SearchParams{}
	if in != nil {
		params.Desc = in.Desc
		params.Request = in.Request
		params.Since = in.Since
	} else {
		return nil
	}

	res, err := s.searchRepo.SearchPersons(params)
	if err != nil {
		return err
	}

	for _, a := range res {
		ar := grpcModels.TransformPersonRPC(a)

		if err := stream.Send(ar); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) SearchOrganizations(in *searchRpc.SearchParams, stream searchRpc.Search_SearchOrganizationsServer) error {
	params := &baseModels.SearchParams{}
	if in != nil {
		params.Desc = in.Desc
		params.Request = in.Request
		params.Since = in.Since
	} else {
		return nil
	}

	res, err := s.searchRepo.SearchOrganizations(params)
	if err != nil {
		return err
	}

	for _, a := range res {
		ar := grpcModels.TransformOrganizationRPC(a)

		if err := stream.Send(ar); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) SearchVacancies(in *searchRpc.SearchParams, stream searchRpc.Search_SearchVacanciesServer) error {
	params := &baseModels.SearchParams{}
	if in != nil {
		params.Desc = in.Desc
		params.Request = in.Request
		params.Since = in.Since
	} else {
		return nil
	}

	res, err := s.searchRepo.SearchVacancies(params)
	if err != nil {
		return err
	}

	for _, a := range res {
		ar := grpcModels.TransformVacanciesRPC(a)

		if err := stream.Send(ar); err != nil {
			return err
		}
	}
	return nil
}

func NewSearchServer(u searchInterfaces.SearchRepository) searchRpc.SearchServer {
	return &server{searchRepo: u}
}
