package database

import (
	"log"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var dbFile = "sqlite.db"
var Db *sqlite.Conn

func Initialise() {
	openDatabase()
	createFilesTable()
}

func openDatabase() {
	conn, err := sqlite.OpenConn(dbFile, sqlite.OpenReadWrite|sqlite.OpenCreate)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	Db = conn

	err = sqlitex.ExecuteTransient(conn, "PRAGMA journal_mode=WAL;", nil)
	if err != nil {
		log.Fatalf("Failed to set WAL mode: %v", err)
	}
}

func CloseDatabase() {
	if Db != nil {
		Db.Close()
	}
}

func createFilesTable() {
	stmt := `CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dir_path TEXT NOT NULL,
			filename TEXT NOT NULL,
			date_added TEXT NOT NULL,
			mdate TEXT NOT NULL
		);`
	err := sqlitex.ExecuteTransient(Db, stmt, nil)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
