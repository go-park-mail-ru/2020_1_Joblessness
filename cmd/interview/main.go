package main

import (
	"github.com/kataras/golog"
	"google.golang.org/grpc"
	interviewPostgres "joblessness/haha/interview/repository/postgres"
	"joblessness/haha/utils/database"
	interviewRpc "joblessness/interviewService/rpc"
	interviewServer "joblessness/interviewService/server"
	"net"
)

func main() {
	db, err := database.OpenDatabase()
	if err != nil {
		golog.Error(err.Error())
		return
	}
	defer db.Close()

	repo := interviewPostgres.NewInterviewRepository(db)
	list, err := net.Listen("tcp", "127.0.0.1:8003")
	if err != nil {
		golog.Error(err.Error())
		return
	}

	server := grpc.NewServer()
	interviewRpc.RegisterInterviewServer(server, interviewServer.NewInterviewServer(repo))
	err = server.Serve(list)
	if err != nil {
		golog.Error("Server interview failed: ", err)
	}
}
