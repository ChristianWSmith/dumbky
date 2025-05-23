package db

import (
	"database/sql"
	"dumbky/internal/log"
	"strconv"
)

func migrate(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT)`)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = db.Exec(`INSERT OR IGNORE INTO meta (key, value) VALUES ('schema_version', '0')`)
	if err != nil {
		log.Error(err)
		return err
	}

	var versionStr string
	err = db.QueryRow(`SELECT value FROM meta WHERE key = 'schema_version'`).Scan(&versionStr)
	if err != nil {
		log.Error(err)
		return err
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		log.Error(err)
		return err
	}

	switch version {
	case 0:
		if err := migrateToV1(db); err != nil {
			log.Error(err)
			return err
		}
		_, err = db.Exec(`UPDATE meta SET value = '1' WHERE key = 'schema_version'`)
		if err != nil {
			log.Error(err)
			return err
		}
		version = 1
		fallthrough
	default:
	}

	return nil
}

func migrateToV1(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE collections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(name)
	);
	CREATE TABLE requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		collection_name TEXT NOT NULL,
		name TEXT NOT NULL,
		payload TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (collection_name) REFERENCES collections(name),
		UNIQUE(collection_name, name)
	);
	CREATE TABLE environments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		payload TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(name)
	);
	INSERT INTO collections (name) VALUES ('');`) // TODO: collection name
	if err != nil {
		log.Error(err)
	}
	return err
}
