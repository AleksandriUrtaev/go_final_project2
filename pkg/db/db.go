package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const defaultDBFile = "scheduler.db"

func Init() (*sql.DB, error) {
	// путь из переменной окр.
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	// проверяем файл
	_, err := os.Stat(dbFile)

	var nCreate bool

	if err == nil {
		nCreate = false
	} else if os.IsNotExist(err) {
		nCreate = true
	} else {
		return nil, fmt.Errorf("stat db file: %w", err)
	}

	// откр. БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	//
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	// создаём таб. при первом запуске
	if nCreate {
		if err := createTable(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	schema := `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128),
		comment TEXT,
		repeat VARCHAR(128)
	);

	CREATE INDEX idx_date ON scheduler(date);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	return nil
}
