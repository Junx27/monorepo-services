package database

import (
	"context"
	"fmt"
	"net/url"
	"transaction/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func New(ctx context.Context, cfg config.Config) error {
	pass := url.QueryEscape(cfg.DatabasePassword)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s",
		cfg.DatabaseUsername,
		pass,
		fmt.Sprintf("%s:%s", cfg.DatabaseHost, cfg.DatabasePort),
		cfg.DatabaseName,
	)

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}

	DB = conn

	return nil
}
