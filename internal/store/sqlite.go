package store

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	DB = db

	if err := CreateTable(); err != nil {
		return err
	}
	return nil
}

func CreateTable() error {
	_, err := DB.Exec(
		`CREATE TABLE IF NOT EXISTS expenses(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			amount INTEGER,
			category TEXT
		)`,
	)
	if err != nil {
		return err
	}
	return nil
}
