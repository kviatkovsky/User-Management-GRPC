package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kviatkovsky/User-Management-gRPC/internal/domain/models/users"
	"github.com/kviatkovsky/User-Management-gRPC/pkg/client/postgresql"
)

type PostgresqlUserRepository struct {
	client postgresql.Client
	logger *slog.Logger
}

func (r *PostgresqlUserRepository) Create(ctx context.Context, user *users.User) error {
	var id uuid.UUID
	q := `
		INSERT INTO users 
		    (email, password) 
		VALUES 
		    ($1, $2)
		RETURNING id
	`

	err := r.client.QueryRow(ctx, q, user.Email, user.PassHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Details: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			r.logger.Error("User creation failed", "error", newErr)

			return newErr
		}
		r.logger.Error("User creation failed with unknown error", "error", err.Error())

		return err
	}

	user.ID = id.String()
	r.logger.Info("User created successfully", "id", user.ID, "user", user)

	return nil
}

func (r *PostgresqlUserRepository) FindAll(ctx context.Context) (u []users.User, err error) {
	q := `
		Select id, email FROM users
	`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Details: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			r.logger.Error("Getting all users failed", "error", newErr)

			return nil, newErr
		}
		r.logger.Error("Getting all users failed", "error", err)
		return nil, err
	}

	usersCol := []users.User{}

	for rows.Next() {
		var usr users.User

		err = rows.Scan(&usr.ID, &usr.Email)
		if err != nil {
			r.logger.Error("Getting all users mapping failed", "error", err)
		}

		usersCol = append(usersCol, usr)
	}

	return usersCol, nil
}

func (r *PostgresqlUserRepository) FindOne(ctx context.Context, id string) (users.User, error) {
	var usr users.User
	q := `
		Select id, email FROM users WHERE id = $1
	`
	err := r.client.QueryRow(ctx, q, id).Scan(&usr.ID, &usr.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Details: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			r.logger.Error("Getting user by id failed", "error", newErr)

			return usr, newErr
		}
		r.logger.Error("Getting user by id failed", "error", err)
		return usr, err
	}

	return usr, nil
}

func (r *PostgresqlUserRepository) FindByEmail(ctx context.Context, email string) (users.User, error) {
	var usr users.User
	q := `
		Select id, email FROM users WHERE email = $1
	`
	err := r.client.QueryRow(ctx, q, email).Scan(&usr.ID, &usr.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Details: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			r.logger.Error("Getting user by Email failed", "error", newErr)

			return usr, newErr
		}
		r.logger.Error("Getting user by Email failed", "error", err)
		return usr, err
	}

	return usr, nil
}

func (r *PostgresqlUserRepository) Update(ctx context.Context, user *users.User) error {
	var (
		q    string
		args []interface{}
	)

	if len(user.PassHash) == 0 {
		q = `
			UPDATE users SET email = $1 WHERE id = $2
		`
		args = append(args, user.Email, user.ID)
	} else {
		q = `
			UPDATE users SET email = $1, password = $2 WHERE id = $3
		`
		args = append(args, user.Email, user.PassHash, user.ID)
	}
	tag, err := r.client.Exec(ctx, q, args...)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Details: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			r.logger.Error("Updating user failed", "error", newErr)

			return newErr
		}

		r.logger.Error("Updating user failed", "error", err)
	}

	r.logger.Info("Updated", "count", tag.RowsAffected())
	return nil
}

func (r *PostgresqlUserRepository) Delete(ctx context.Context, id string) error {
	q := `
		DELETE FROM users WHERE id = $1
	`

	tag, err := r.client.Exec(ctx, q, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			newErr := fmt.Errorf("SQL Error: %s, Details: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)
			r.logger.Error("Deleting user failed", "error", newErr)

			return newErr
		}

		r.logger.Error("Deleting user failed", "error", err)
	}

	r.logger.Info("removed", "count", tag.RowsAffected())
	return nil
}

func NewUserRepository(client postgresql.Client, logger *slog.Logger) Repository {
	return &PostgresqlUserRepository{
		client: client,
		logger: logger,
	}
}
