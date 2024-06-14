package users

import (
	"context"

	"github.com/kviatkovsky/User-Management-gRPC/internal/grpc/user"
)

type User struct {
	ID       string
	Email    string
	PassHash []byte
}

type Service interface {
	GetUsers(ctx context.Context, in *user.GetUsersRequest) (*user.GetUsersResponse, error)
	CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.CreateUserResponse, error)
	GetUserByID(ctx context.Context, in *user.GetUserByIDRequest) (*user.GetUserByIDResponse, error)
	UpdateUser(ctx context.Context, in *user.UpdateUserRequest) (*user.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, in *user.DeleteUserRequest) (*user.DeleteUserResponse, error)
	GetUserByEmail(ctx context.Context, in *user.GetUserByEmailRequest) (*user.GetUserByEmailResponse, error)
}
