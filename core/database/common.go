package database

import (
	"context"
	"fmt"
	"log"
	"zene/core/logger"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func doesTableExist(tableName string, conn *sqlite.Conn) (bool, error) {
	stmt, err := conn.Prepare(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = $table_name;`)
	if err != nil {
		return false, fmt.Errorf("preparing stmt for doesTableExist: %v", err)
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

func createTable(ctx context.Context, tableName string, createSql string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in createTable: %v", err)
	}
	defer DbPool.Put(conn)

	tableExists, err := doesTableExist(tableName, conn)

	if err != nil {
		return fmt.Errorf("checking if %s table exists: %s", tableName, err)
	}

	if !tableExists {
		stmt := createSql
		err := sqlitex.ExecuteTransient(conn, stmt, nil)
		if err != nil {
			return fmt.Errorf("create %s table: %v", tableName, err)
		} else {
			logger.Printf("Database: %s table created", tableName)
			return nil
		}
	} else {
		logger.Printf("Database: %s table already exists", tableName)
		return nil
	}
}

func createTrigger(ctx context.Context, triggerName string, triggerSQL string) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		logger.Printf("taking a db conn from the pool in createTrigger: %v", err)
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
		logger.Printf("%s trigger already exists", triggerName)
	} else if err != nil {
		logger.Printf("Error checking for %s trigger: %s", triggerName, err)
	} else {
		stmt, err := conn.Prepare(triggerSQL)
		if err != nil {
			logger.Printf("Error preparing %s trigger: %s", triggerName, err)
			return
		}
		defer stmt.Finalize()

		_, err = stmt.Step()
		if err != nil {
			logger.Printf("Error creating %s trigger: %s", triggerName, err)
			return
		}
		logger.Printf("%s trigger created", triggerName)
	}
}

func createIndex(ctx context.Context, indexName, indexTable, indexColumn string, indexUnique bool) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		logger.Printf("taking a db conn from the pool in createIndex: %v", err)
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
		logger.Printf("Database: Error checking for %s index: %s", indexName, err)
		return
	}
	if hasRow {
		logger.Printf("Database: %s index already exists", indexName)
		return
	}

	var sql string
	if indexUnique {
		sql = fmt.Sprintf("CREATE UNIQUE INDEX %q ON %q (%s);", indexName, indexTable, indexColumn)
	} else {
		sql = fmt.Sprintf("CREATE INDEX %q ON %q (%s);", indexName, indexTable, indexColumn)
	}

	stmt2, err := conn.Prepare(sql)
	if err != nil {
		logger.Printf("Database: Failed to prepare CREATE INDEX for %s: %v", indexName, err)
		return
	}
	defer stmt2.Finalize()

	_, err = stmt2.Step()
	if err != nil {
		logger.Printf("Database: Error creating %s index: %s", indexName, err)
		return
	}

	logger.Printf("Database: %s index created", indexName)
}
