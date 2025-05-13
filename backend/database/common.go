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

func createTriggerIfNotExists(triggerName string, triggerSQL string) {
	checkQuery := Db.Prep("SELECT name FROM sqlite_master WHERE type='trigger' AND name=$triggername")
	checkQuery.SetText("$triggername", triggerName)

	hasRow, err := checkQuery.Step()
	if hasRow {
		log.Printf("%s trigger already exists", triggerName)
	} else if err != nil {
		log.Printf("Error checking for %s trigger: %s", triggerName, err)
	} else {
		log.Printf("Creating %s trigger", triggerName)
		stmt, err := Db.Prepare(triggerSQL)
		if err != nil {
			log.Printf("Error preparing %s trigger: %s", triggerName, err)
			return
		}
		defer stmt.Finalize()

		_, err = stmt.Step()
		if err != nil {
			log.Printf("Error creating %s trigger: %s", triggerName, err)
			return
		}
		log.Printf("%s trigger created", triggerName)
	}
}
