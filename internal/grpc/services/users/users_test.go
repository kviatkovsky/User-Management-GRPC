package users

import (
	"context"
	"fmt"
	"testing"

	"github.com/kviatkovsky/User-Management-gRPC/internal/domain/models/users"
	"github.com/kviatkovsky/User-Management-gRPC/internal/grpc/user"
	"github.com/kviatkovsky/User-Management-gRPC/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	log := logger.SetupLogger("local")
	mockRepo := new(mockUserRepository)

	s := &Server{
		Logger: log,
		Repo:   mockRepo,
	}

	ctx := context.Background()
	req := &user.GetUsersRequest{}

	t.Run("Should return list of users", func(t *testing.T) {
		usrs := []*user.UserResponse{
			{Id: "success", Email: "user1@example.com"},
			{Id: "2", Email: "user2@example.com"},
		}

		res, _ := s.GetUsers(ctx, req)

		assert.Equal(t, res.Users, usrs)
	})
}

func TestGetUserById(t *testing.T) {
	log := logger.SetupLogger("local")
	mockRepo := new(mockUserRepository)

	s := &Server{
		Logger: log,
		Repo:   mockRepo,
	}

	ctx := context.Background()
	req := &user.GetUserByIDRequest{}

	t.Run("Should return user by ID", func(t *testing.T) {
		req.Id = "success"
		usr := &user.GetUserByIDResponse{
			User: &user.UserResponse{
				Id:    "success",
				Email: "user1@example.com",
			},
		}

		res, err := s.GetUserByID(ctx, req)

		assert.Equal(t, res, usr)
		assert.Equal(t, err, nil)
	})

	t.Run("Should throw an Error", func(t *testing.T) {
		req.Id = "2"
		expectErr := "rpc error: code = Internal desc = user not found"

		_, err := s.GetUserByID(ctx, req)

		assert.EqualError(t, err, expectErr)
	})
}

func TestGetUserByEmail(t *testing.T) {
	log := logger.SetupLogger("local")
	mockRepo := new(mockUserRepository)

	s := &Server{
		Logger: log,
		Repo:   mockRepo,
	}

	ctx := context.Background()
	req := &user.GetUserByEmailRequest{}

	t.Run("Should return user by Email", func(t *testing.T) {
		req.Email = "user1@example.com"
		usr := &user.GetUserByEmailResponse{
			User: &user.UserResponse{
				Id:    "success",
				Email: "user1@example.com",
			},
		}

		res, err := s.GetUserByEmail(ctx, req)

		assert.Equal(t, res, usr)
		assert.Equal(t, err, nil)

	})

	t.Run("Should throw an Error", func(t *testing.T) {
		req.Email = "error@example.com"
		expectErr := "rpc error: code = Internal desc = user not found"

		_, err := s.GetUserByEmail(ctx, req)

		assert.EqualError(t, err, expectErr)
	})
}

func TestCreateUser(t *testing.T) {
	log := logger.SetupLogger("local")
	mockRepo := new(mockUserRepository)

	s := &Server{
		Logger: log,
		Repo:   mockRepo,
	}

	ctx := context.Background()
	req := &user.CreateUserRequest{}

	t.Run("Create user success", func(t *testing.T) {
		usr := user.CreateUserResponse{Id: "user_id_1"}
		req.Email = "user1@success.com"

		res, err := s.CreateUser(ctx, req)

		assert.Equal(t, res.GetId(), usr.GetId())
		assert.Equal(t, err, nil)

	})

	t.Run("Create user failed", func(t *testing.T) {
		req.Email = "user1@failed.com"
		_, err := s.CreateUser(ctx, req)

		assert.EqualError(t, err, "rpc error: code = Internal desc = can't create user")
	})
}

func TestDeleteUser(t *testing.T) {
	log := logger.SetupLogger("local")
	mockRepo := new(mockUserRepository)

	s := &Server{
		Logger: log,
		Repo:   mockRepo,
	}

	ctx := context.Background()
	req := &user.DeleteUserRequest{}

	t.Run("Delete user success", func(t *testing.T) {
		usr := user.DeleteUserResponse{Id: "success"}
		req.Id = "success"

		res, err := s.DeleteUser(ctx, req)

		assert.Equal(t, res.GetId(), usr.GetId())
		assert.Equal(t, err, nil)
	})

	t.Run("Delete user Failed", func(t *testing.T) {
		req.Id = "error"

		_, err := s.DeleteUser(ctx, req)
		assert.EqualError(t, err, "rpc error: code = Internal desc = can't delete user")
	})
}

func TestUpdateUser(t *testing.T) {
	log := logger.SetupLogger("local")
	mockRepo := new(mockUserRepository)

	s := &Server{
		Logger: log,
		Repo:   mockRepo,
	}

	ctx := context.Background()
	req := &user.UpdateUserRequest{}

	t.Run("Success update user", func(t *testing.T) {
		req.Email = "user1@example.com"
		req.Id = "success"
		req.Password = "password"
		usr := user.UpdateUserResponse{
			User: &user.UserResponse{
				Email: "user1@success.com",
				Id:    "success",
			},
		}

		res, err := s.UpdateUser(ctx, req)

		assert.Equal(t, res.User, usr.User)
		assert.Equal(t, err, nil)
	})

	t.Run("Update user error getting user", func(t *testing.T) {
		req.Email = "user1@example.com"
		req.Id = "user_id"
		req.Password = "password"

		_, err := s.UpdateUser(ctx, req)

		expectErr := "rpc error: code = Internal desc = user not found"

		assert.EqualError(t, err, expectErr)
	})

	t.Run("Update user error getting user", func(t *testing.T) {
		req.Email = "user1@example.com"
		req.Id = "failed_update"
		req.Password = "password"

		_, err := s.UpdateUser(ctx, req)

		expectErr := "rpc error: code = Internal desc = error during update"

		assert.EqualError(t, err, expectErr)
	})

}

type mockUserRepository struct{}

func (m *mockUserRepository) Create(ctx context.Context, user *users.User) error {
	if user.Email == "user1@success.com" {
		user.ID = "user_id_1"
		return nil
	}

	return fmt.Errorf("can't create user")
}

func (m *mockUserRepository) FindAll(ctx context.Context) (u []users.User, err error) {
	usrs := []users.User{
		{ID: "success", Email: "user1@example.com"},
		{ID: "2", Email: "user2@example.com"},
	}

	return usrs, nil
}

func (m *mockUserRepository) FindOne(ctx context.Context, id string) (users.User, error) {
	if id == "success" {
		usr := users.User{
			ID:    "success",
			Email: "user1@example.com",
		}

		return usr, nil
	} else if id == "failed_update" {
		usr := users.User{
			ID:    "failed_update",
			Email: "user1@example.com",
		}

		return usr, nil
	}

	return users.User{}, fmt.Errorf("user not found")
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (users.User, error) {
	if email == "user1@example.com" {
		usr := users.User{
			ID:    "success",
			Email: "user1@example.com",
		}

		return usr, nil
	}

	return users.User{}, fmt.Errorf("user not found")
}

func (m *mockUserRepository) Update(ctx context.Context, user *users.User) error {
	if user.ID == "failed_update" {
		return fmt.Errorf("error during update")
	}

	user.ID = "success"
	user.Email = "user1@success.com"

	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	if id == "success" {
		return nil
	}

	return fmt.Errorf("can't delete user")
}
