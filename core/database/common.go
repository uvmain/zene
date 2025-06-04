package database

import (
	"context"
	"fmt"
	"log"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func doesTableExist(tableName string, conn *sqlite.Conn) (bool, error) {
	stmt, err := conn.Prepare(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = $table_name;`)
	if err != nil {
		return false, fmt.Errorf("error preparing doesTableExist stmt: %v", err)
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
	log.Printf("creating table %s", tableName)
	dbMutex.Lock()
	defer dbMutex.Unlock()
	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	tableExists, err := doesTableExist(tableName, conn)

	if err != nil {
		log.Printf("Error checking if %s table exists: %s", tableName, err)
	}
	if !tableExists {
		stmt := createSql
		err := sqlitex.ExecuteTransient(conn, stmt, nil)
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
	ctx := context.Background()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Println("failed to take a db conn from the pool")
	}
	defer DbPool.Put(conn)

	stmt, err := conn.Prepare("SELECT name FROM sqlite_master WHERE type='trigger' AND name=$triggername")
	if err != nil {
		log.Fatalf("Failed to prepare stmtCreateTriggerIfNotExists: %v", err)
	}
	defer stmt.Finalize()
	stmt.SetText("$triggername", triggerName)

	hasRow, err := stmt.Step()
	if hasRow {
		log.Printf("%s trigger already exists", triggerName)
	} else if err != nil {
		log.Printf("Error checking for %s trigger: %s", triggerName, err)
	} else {
		log.Printf("Creating %s trigger", triggerName)
		stmt, err := conn.Prepare(triggerSQL)
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
