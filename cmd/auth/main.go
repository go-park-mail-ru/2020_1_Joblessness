package main

import (
	"github.com/kataras/golog"
	"google.golang.org/grpc"
	"joblessness/authService/grpc"
	"joblessness/authService/server"
	"joblessness/haha/auth/repository/postgres"
	"joblessness/haha/utils/database"
	"net"
)

func main() {
	db, err := database.OpenDatabase()
	if err != nil {
		golog.Error(err.Error())
		return
	}
	defer db.Close()

	repository := authPostgres.NewAuthRepository(db)

	listen, err := net.Listen("tcp", "127.0.0.1:8004")
	if err != nil {
		golog.Error(err.Error())
		return
	}

	server := grpc.NewServer()
	authGrpc.RegisterAuthServer(server, authServer.NewAuthServer(repository))
	err = server.Serve(listen)
	if err != nil {
		golog.Error("Server auth failed")
	}
}
