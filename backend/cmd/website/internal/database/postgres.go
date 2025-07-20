package database

import (
	"database/sql"

	"github.com/DeepAung/gradient/cmd/website/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func InitPostgresDB(cfg config.Database) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// TODO: how should i set these, and should i use connection pool?
	// db.SetConnMaxIdleTime(d time.Duration)
	// db.SetConnMaxLifetime(d time.Duration)
	// db.SetMaxIdleConns(n int)
	// db.SetMaxOpenConns(n int)

	return db
}
