package db

import (
	"database/sql"
	"dumbky/internal/log"
	"dumbky/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
	dbPath, err := utils.GetDBFilePath()
	if err != nil {
		log.Error(err)
		return err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Error(err)
		return err
	}
	if err := db.Ping(); err != nil {
		log.Error(err)
		return err
	}
	DB = db

	if err := migrate(db); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
