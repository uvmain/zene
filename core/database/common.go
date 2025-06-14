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
		return false, fmt.Errorf("error preparing stmt for doesTableExist: %v", err)
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

func createTable(ctx context.Context, tableName string, createSql string) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("Database: failed to take a db conn from the pool in createTable: %v", err)
		return
	}
	defer DbPool.Put(conn)

	tableExists, err := doesTableExist(tableName, conn)

	if err != nil {
		log.Printf("Database: Error checking if %s table exists: %s", tableName, err)
	}

	if !tableExists {
		stmt := createSql
		err := sqlitex.ExecuteTransient(conn, stmt, nil)
		if err != nil {
			log.Fatalf("Database: Failed to create %s table: %v", tableName, err)
		} else {
			log.Printf("Database: %s table created", tableName)
		}
	} else {
		log.Printf("Database: %s table already exists", tableName)
	}
}

func createTrigger(ctx context.Context, triggerName string, triggerSQL string) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in createTrigger: %v", err)
		return
	}
	defer DbPool.Put(conn)

	stmt, err := conn.Prepare("SELECT name FROM sqlite_master WHERE type='trigger' AND name=$triggername")
	if err != nil {
		log.Fatalf("Failed to prepare stmt for createTrigger: %v", err)
	}
	defer stmt.Finalize()
	stmt.SetText("$triggername", triggerName)

	hasRow, err := stmt.Step()
	if hasRow {
		log.Printf("%s trigger already exists", triggerName)
	} else if err != nil {
		log.Printf("Error checking for %s trigger: %s", triggerName, err)
	} else {
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

func createIndex(ctx context.Context, indexName, indexTable, indexColumn string, indexUnique bool) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in createIndex: %v", err)
		return
	}
	defer DbPool.Put(conn)

	stmt, err := conn.Prepare("SELECT name FROM sqlite_master WHERE type='index' AND name=?")
	if err != nil {
		log.Fatalf("Failed to prepare stmt for createIndex: %v", err)
	}
	defer stmt.Finalize()
	stmt.BindText(1, indexName)

	hasRow, err := stmt.Step()
	if err != nil {
		log.Printf("Database: Error checking for %s index: %s", indexName, err)
		return
	}
	if hasRow {
		log.Printf("Database: %s index already exists", indexName)
		return
	}

	var sql string
	if indexUnique {
		sql = fmt.Sprintf("CREATE UNIQUE INDEX %q ON %q (%q);", indexName, indexTable, indexColumn)
	} else {
		sql = fmt.Sprintf("CREATE INDEX %q ON %q (%q);", indexName, indexTable, indexColumn)
	}

	stmt2, err := conn.Prepare(sql)
	if err != nil {
		log.Printf("Database: Failed to prepare CREATE INDEX for %s: %v", indexName, err)
		return
	}
	defer stmt2.Finalize()

	_, err = stmt2.Step()
	if err != nil {
		log.Printf("Database: Error creating %s index: %s", indexName, err)
		return
	}

	log.Printf("Database: %s index created", indexName)
}
