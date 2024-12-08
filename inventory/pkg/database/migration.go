package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"inventory/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg config.Config) error {
	pass := url.QueryEscape(cfg.DatabasePassword)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DatabaseUsername,
		pass,
		fmt.Sprintf("%s:%s", cfg.DatabaseHost, cfg.DatabasePort),
		cfg.DatabaseName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationPath),
		"postgres", driver)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	return nil
}
