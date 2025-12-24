package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const (
	schema = "CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8) NOT NULL DEFAULT \"\", title VARCHAR(256) NOT NULL DEFAULT \"\", comment TEXT NOT NULL DEFAULT \"\", repeat VARCHAR(128) NOT NULL DEFAULT \"\"); CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);"
)

var db *sql.DB

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}

	if install {
		_, err = db.Exec(schema)
		return err
	}

	return nil
}
