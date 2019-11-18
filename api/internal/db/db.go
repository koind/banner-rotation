package db

import (
	"context"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/koind/banner-rotation/api/internal/config"
	"github.com/pkg/errors"
)

// Creates a connection pool and returns the connection itself
func IntPostgres(ctx context.Context, options config.Postgres) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", options.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "an error occurred while creating the connection pool")
	}

	db.SetMaxOpenConns(options.MaxOpenConns)
	db.SetMaxIdleConns(options.MaxIdleConns)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "an error occurred while connecting to the database")
	}

	return db, nil
}
