package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDatabase : Initialize the database
func InitDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	create := `
		CREATE TABLE IF NOT EXISTS Secrets (
		token 		VARCHAR(32) PRIMARY KEY,
		secret 		TEXT,
		password 	VARCHAR(255),
		expire 		TEXT,
		maxviews 	INTEGER DEFAULT 1,
		views 		INTEGER
		);`

	stmt, err := db.Prepare(create)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	stmt.Exec()

	err = db.Ping()
	if err != nil {
		log.Fatal("Pinging database failed: ", err)
	}

	return db
}
