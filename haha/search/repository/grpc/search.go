package searchGrpc

import (
	"google.golang.org/grpc"
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