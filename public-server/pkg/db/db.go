package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB(dataSourceName string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatal("sqlx.Connect: ", err)
	}
	return db
}
