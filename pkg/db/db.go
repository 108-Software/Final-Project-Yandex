package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}


	if err := db.Ping(); err != nil {
		return err
	}

	
	if install {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}


func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
