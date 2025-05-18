package database

import (
	"log"

	"zombiezen.com/go/sqlite/sqlitex"
)

func doesTableExist(tableName string) (bool, error) {
	stmt := stmtDoesTableExist
	stmt.Reset()
	stmt.ClearBindings()
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
		err := sqlitex.ExecuteTransient(DbRW, stmt, nil)
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
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtCreateTriggerIfNotExists
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$triggername", triggerName)

	hasRow, err := stmt.Step()
	if hasRow {
		log.Printf("%s trigger already exists", triggerName)
	} else if err != nil {
		log.Printf("Error checking for %s trigger: %s", triggerName, err)
	} else {
		log.Printf("Creating %s trigger", triggerName)
		stmt, err := DbRW.Prepare(triggerSQL)
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
