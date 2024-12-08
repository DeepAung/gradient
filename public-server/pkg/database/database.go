package database

import (
	"log"
	"os"

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

func RunSQL(db *sqlx.DB, migrateSourceName string) {
	b, err := os.ReadFile(migrateSourceName)
	if err != nil {
		log.Fatal("MigrateDB: os.ReadFile: ", err)
	}
	db.MustExec(string(b))
}
