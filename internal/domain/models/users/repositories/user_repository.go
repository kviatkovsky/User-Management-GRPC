package repositories

import (
	"context"

	"github.com/kviatkovsky/User-Management-gRPC/internal/domain/models/users"
)

type Repository interface {
	Create(ctx context.Context, user *users.User) error
	FindAll(ctx context.Context) (u []users.User, err error)
	FindOne(ctx context.Context, id string) (users.User, error)
	FindByEmail(ctx context.Context, email string) (users.User, error)
	Update(ctx context.Context, user *users.User) error
	Delete(ctx context.Context, id string) error
}
