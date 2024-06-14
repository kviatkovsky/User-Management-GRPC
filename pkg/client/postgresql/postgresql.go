package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kviatkovsky/User-Management-gRPC/internal/config"
	repeatable "github.com/kviatkovsky/User-Management-gRPC/pkg/utils"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, cfg config.StorageConfig) (conn *pgx.Conn, err error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)

		defer cancel()

		conn, err = pgx.Connect(ctx, dsn)

		if err != nil {
			return err
		}

		return nil
	}, cfg.Attempts, 1*time.Second)

	if err != nil {
		log.Fatal("error connecting with tries to the database: ", err)
	}

	return conn, nil
}
