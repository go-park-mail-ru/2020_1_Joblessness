package main

import (
	"github.com/kataras/golog"
	"google.golang.org/grpc"
	searchPostgres "joblessness/haha/search/repository/postgres"
	"joblessness/haha/utils/database"
	searchServer "joblessness/searchService/server"
	"net"
)

func main() {
	db, err := database.OpenDatabase()
	if err != nil {
		golog.Error(err.Error())
		return
	}
	defer db.Close()

	repo := searchPostgres.NewSearchRepository(db)
	list, err := net.Listen("tcp", ":8002")
	if err != nil {
		golog.Error(err.Error())
		return
	}

	server := grpc.NewServer()
	searchServer.NewSearchServer(repo)
	server.Serve(list)
}