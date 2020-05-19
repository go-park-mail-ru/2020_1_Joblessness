package authGrpcRepository

import (
	"context"
	"google.golang.org/grpc"
	"joblessness/authService/grpc"
)

type AuthRepository struct {
	authClient authGrpc.AuthClient
}

func NewRepository(conn *grpc.ClientConn) *AuthRepository {
	return &AuthRepository{authClient: authGrpc.NewAuthClient(conn)}
}

func (r *AuthRepository) RegisterPerson(login, password, name string) (err error) {
	_, err = r.authClient.RegisterPerson(context.Background(), &authGrpc.UserRegister{
		Login:    login,
		Password: password,
		Name:     name,
	})

	return err
}

func (r *AuthRepository) RegisterOrganization(login, password, name string) (err error) {
	_, err = r.authClient.RegisterOrganization(context.Background(), &authGrpc.UserRegister{
		Login:    login,
		Password: password,
		Name:     name,
	})

	return err
}

func (r *AuthRepository) Login(login, password, SID string) (userID uint64, err error) {
	protoUserID, err := r.authClient.Login(context.Background(), &authGrpc.UserLoginParams{
		Login:    login,
		Password: password,
		Sid:      SID,
	})
	if err != nil {
		return userID, err
	}

	return protoUserID.Id, err
}

func (r *AuthRepository) Logout(sessionID string) (err error) {
	_, err = r.authClient.Logout(context.Background(), &authGrpc.SessionID{Id: sessionID})

	return err
}

func (r *AuthRepository) SessionExists(sessionID string) (userID uint64, err error) {
	protoUserID, err := r.authClient.SessionExists(context.Background(), &authGrpc.SessionID{Id: sessionID})
	if err != nil {
		return userID, err
	}

	return protoUserID.Id, err
}

func (r *AuthRepository) DoesUserExists(login string) (err error) {
	_, err = r.authClient.DoesUserExists(context.Background(), &authGrpc.UserLogin{Login: login})

	return err
}

func (r *AuthRepository) GetRole(userID uint64) (role string, err error) {
	protoRole, err := r.authClient.GetRole(context.Background(), &authGrpc.UserID{Id: userID})
	if err != nil {
		return role, err
	}

	return protoRole.Role, err
}
