package users

import (
	"context"
	"log/slog"

	"github.com/kviatkovsky/User-Management-gRPC/internal/domain/models/users"
	"github.com/kviatkovsky/User-Management-gRPC/internal/domain/models/users/repositories"
	"github.com/kviatkovsky/User-Management-gRPC/internal/grpc/user"
	"github.com/kviatkovsky/User-Management-gRPC/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Logger *slog.Logger
	Repo   repositories.Repository
	user.UnimplementedUserServiceServer
}

func (s *Server) GetUsers(ctx context.Context, in *user.GetUsersRequest) (*user.GetUsersResponse, error) {
	const op = "users/users.GetUsers"
	usersResp := user.GetUsersResponse{}

	log := s.Logger.With(
		slog.String("op", op),
	)

	log.Info("getting users")
	userCol, err := s.Repo.FindAll(ctx)
	if err != nil {
		log.Error("Error getting users", "error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, usr := range userCol {
		usersResp.Users = append(usersResp.Users, &user.UserResponse{
			Id:    usr.ID,
			Email: usr.Email,
		})
	}

	return &usersResp, nil
}

func (s *Server) CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	const op = "users/users.CreateUser"
	log := s.Logger.With(
		slog.String("op", op),
	)

	usersResp := user.CreateUserResponse{}
	var userToCreate users.User
	userToCreate.Email = in.GetEmail()
	userToCreate.PassHash = []byte(in.GetPassword())

	err := s.Repo.Create(ctx, &userToCreate)
	if err != nil {
		log.Error("Error creating user", "error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	usersResp.Id = userToCreate.ID

	return &usersResp, nil
}

func (s *Server) GetUserByID(ctx context.Context, in *user.GetUserByIDRequest) (*user.GetUserByIDResponse, error) {
	const op = "users/users.GetUserByID"
	log := s.Logger.With(
		slog.String("op", op),
	)

	log.Info("getting user by ID")
	usr, err := s.Repo.FindOne(ctx, in.GetId())
	if err != nil {
		log.Error("Error getting users", "error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.GetUserByIDResponse{
		User: &user.UserResponse{
			Id:    usr.ID,
			Email: usr.Email,
		},
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	const op = "users/users.UpdateUser"
	log := s.Logger.With(
		slog.String("op", op),
	)

	log.Info("Updating user by ID")

	usr, err := s.Repo.FindOne(ctx, in.GetId())

	if err != nil {
		log.Error("Error getting user", "error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	utils.PrepareUserToUpdate(&usr, in)

	err = s.Repo.Update(ctx, &usr)
	if err != nil {
		log.Error("Error updating users", "error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.UpdateUserResponse{
		User: &user.UserResponse{
			Id:    usr.ID,
			Email: usr.Email,
		},
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, in *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	const op = "users/users.DeleteUser"
	log := s.Logger.With(
		slog.String("op", op),
	)

	log.Info("deleting user by ID")
	err := s.Repo.Delete(ctx, in.GetId())
	if err != nil {
		log.Error("Error deleting user", "error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.DeleteUserResponse{
		Id: in.GetId(),
	}, nil
}

func (s *Server) GetUserByEmail(ctx context.Context, in *user.GetUserByEmailRequest) (*user.GetUserByEmailResponse, error) {
	const op = "users/users.GetUserByEmail"
	log := s.Logger.With(
		slog.String("op", op),
	)

	log.Info("getting user by ID")
	usr, err := s.Repo.FindByEmail(ctx, in.GetEmail())
	if err != nil {
		log.Error("Error getting users", "error", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &user.GetUserByEmailResponse{
		User: &user.UserResponse{
			Id:    usr.ID,
			Email: usr.Email,
		},
	}, nil
}
