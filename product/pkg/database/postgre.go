package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"product/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func New(ctx context.Context, cfg config.Config) error {
	pass := url.QueryEscape(cfg.DatabasePassword)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DatabaseUsername,
		pass,
		fmt.Sprintf("%s:%s", cfg.DatabaseHost, cfg.DatabasePort),
		cfg.DatabaseName,
	)

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatalln(err.Error())
	}

	DB = conn

	return nil
}

func Transaction(ctx context.Context, tf func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = tf(ctx, tx)
	if err != nil {
		log.Println(err.Error())
		tx.Rollback(ctx)
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println(err.Error())
		tx.Rollback(ctx)
	}

	return nil
}
