package authServer

import (
	"context"
	"joblessness/authService/grpc"
	"joblessness/haha/auth/interfaces"
)

type server struct {
	authRepository authInterfaces.AuthRepository
}

func NewAuthServer(authRepository authInterfaces.AuthRepository) authGrpc.AuthServer {
	return &server{authRepository: authRepository}
}

func (s *server) RegisterPerson(ctx context.Context, in *authGrpc.UserRegister) (*authGrpc.Nothing, error) {
	err := s.authRepository.RegisterPerson(in.Login, in.Password, in.Name)

	return &authGrpc.Nothing{Dummy: true}, err
}

func (s *server) RegisterOrganization(ctx context.Context, in *authGrpc.UserRegister) (*authGrpc.Nothing, error) {
	err := s.authRepository.RegisterOrganization(in.Login, in.Password, in.Name)

	return &authGrpc.Nothing{Dummy: true}, err
}

func (s *server) Login(ctx context.Context, in *authGrpc.UserLoginParams) (*authGrpc.UserID, error) {
	userID, err := s.authRepository.Login(in.Login, in.Password, in.Sid)

	return &authGrpc.UserID{Id: userID}, err
}

func (s *server) Logout(ctx context.Context, in *authGrpc.SessionID) (*authGrpc.Nothing, error) {
	err := s.authRepository.Logout(in.Id)

	return &authGrpc.Nothing{Dummy: true}, err
}

func (s *server) SessionExists(ctx context.Context, in *authGrpc.SessionID) (*authGrpc.UserID, error) {
	userID, err := s.authRepository.SessionExists(in.Id)

	return &authGrpc.UserID{Id: userID}, err
}

func (s *server) DoesUserExists(ctx context.Context, in *authGrpc.UserLogin) (*authGrpc.Nothing, error) {
	err := s.authRepository.DoesUserExists(in.Login)

	return &authGrpc.Nothing{Dummy: true}, err
}

func (s *server) GetRole(ctx context.Context, in *authGrpc.UserID) (*authGrpc.Role, error) {
	role, err := s.authRepository.GetRole(in.Id)

	return &authGrpc.Role{Role: role}, err
}
