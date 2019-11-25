package onetimesecret

import (
	"database/sql"
	"log"

	// SQLlite
	_ "github.com/mattn/go-sqlite3"
)

// Datastore : Repository interface
type Datastore interface {
	GetSecretByToken(token string) (*Secret, error)
	GetSecretByTokenAndPassword(token string, password string) (*Secret, error)
}

// DB : Holds database connection
type DB struct {
	*sql.DB
}

// NewDB : Creates new database connection
func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	create := `
		CREATE TABLE IF NOT EXISTS Secrets (
		token 		VARCHAR(32) PRIMARY KEY,
		secret 		TEXT,
		password 	VARCHAR(255),
		expire 		TEXT DEFAULT (datetime('now', '10 minutes')),
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
	return &DB{db}, nil
}
