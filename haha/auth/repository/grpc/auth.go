package authGrpcRepository

import (
	"context"
	"google.golang.org/grpc"
	"joblessness/authService/grpc"
	"strconv"
)

type repository struct {
	authClient authGrpc.AuthClient
}

func NewRepository(conn *grpc.ClientConn) *repository {
	return &repository{authClient: authGrpc.NewAuthClient(conn)}
}

func (r *repository) RegisterPerson(login, password, name string) (err error) {
	_, err = r.authClient.RegisterPerson(context.Background(), &authGrpc.UserRegister{
		Login:    login,
		Password: password,
		Name:     name,
	})

	return err
}

func (r *repository) RegisterOrganization(login, password, name string) (err error) {
	_, err = r.authClient.RegisterOrganization(context.Background(), &authGrpc.UserRegister{
		Login:    login,
		Password: password,
		Name:     name,
	})

	return err
}

func (r *repository) Login(login, password, SID string) (userID uint64, err error) {
	protoUserID, err := r.authClient.Login(context.Background(), &authGrpc.UserLoginParams{
		Login:    login,
		Password: password,
		Sid:      SID,
	})
	if err != nil {
		return userID, err
	}

	return strconv.ParseUint(protoUserID.Id, 10, 64)
}

func (r *repository) Logout(sessionId string) (err error) {
	_, err = r.authClient.Logout(context.Background(), &authGrpc.SessionID{Id: sessionId})

	return err
}

func (r *repository) SessionExists(sessionId string) (userID uint64, err error) {
	protoUserID, err := r.authClient.SessionExists(context.Background(), &authGrpc.SessionID{Id: sessionId})
	if err != nil {
		return userID, err
	}

	return strconv.ParseUint(protoUserID.Id, 10, 64)
}

func (r *repository) DoesUserExists(login string) (err error) {
	_, err = r.authClient.DoesUserExists(context.Background(), &authGrpc.UserLogin{Login: login})

	return err
}

func (r *repository) GetRole(userID uint64) (role string, err error) {
	protoRole, err := r.authClient.GetRole(context.Background(), &authGrpc.UserID{Id: string(userID)})
	if err != nil {
		return role, err
	}

	return protoRole.Role, err
}
