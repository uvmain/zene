package database

import (
	"log"

	"zombiezen.com/go/sqlite/sqlitex"
)

func doesTableExist(tableName string) (bool, error) {
	stmt, err := Db.Prepare(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = $table_name;`)

	if err != nil {
		return false, err
	}
	defer stmt.Finalize()

	stmt.SetText("$table_name", tableName)

	if hasRow, err := stmt.Step(); err != nil {
		return false, err
	} else if !hasRow {
		return false, nil
	} else {
		return true, nil
	}

}

func createTable(tableName string, createSql string) {
	table_exists, err := doesTableExist(tableName)

	if err != nil {
		log.Printf("Error checking if %s table exists: %s", tableName, err)
	}
	if !table_exists {
		stmt := createSql
		err := sqlitex.ExecuteTransient(Db, stmt, nil)
		if err != nil {
			log.Fatalf("Failed to create %s table: %v", tableName, err)
		} else {
			log.Printf("%s table created", tableName)
		}
	} else {
		log.Printf("%s table already exists", tableName)
	}
}
