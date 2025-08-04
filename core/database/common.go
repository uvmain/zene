package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/logger"
)

func doesTableExist(ctx context.Context, tableName string) (bool, error) {
	query := `SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`
	var name string
	err := DB.QueryRowContext(ctx, query, tableName).Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("checking if table exists: %v", err)
	}
	return true, nil
}

func createTable(ctx context.Context, tableName string, createSql string) error {
	tableExists, err := doesTableExist(ctx, tableName)
	if err != nil {
		return fmt.Errorf("checking if %s table exists: %v", tableName, err)
	}

	if !tableExists {
		_, err := DB.ExecContext(ctx, createSql)
		if err != nil {
			return fmt.Errorf("create %s table: %v", tableName, err)
		}
		logger.Printf("Database: %s table created", tableName)
	} else {
		logger.Printf("Database: %s table already exists", tableName)
	}
	return nil
}

func createTrigger(ctx context.Context, triggerName string, triggerSQL string) {
	query := "SELECT name FROM sqlite_master WHERE type='trigger' AND name=?"
	var name string
	err := DB.QueryRowContext(ctx, query, triggerName).Scan(&name)

	if err == sql.ErrNoRows {
		_, err := DB.ExecContext(ctx, triggerSQL)
		if err != nil {
			logger.Printf("Database: error creating %s trigger: %v", triggerName, err)
			return
		}
		logger.Printf("Database: %s trigger created", triggerName)
	} else if err != nil {
		logger.Printf("Database: error checking for %s trigger: %v", triggerName, err)
	} else {
		logger.Printf("Database: %s trigger already exists", triggerName)
	}
}

func createIndex(ctx context.Context, indexName, indexTable, indexColumn string, indexUnique bool) {
	query := "SELECT name FROM sqlite_master WHERE type='index' AND name=?"
	var name string
	err := DB.QueryRowContext(ctx, query, indexName).Scan(&name)

	if err == sql.ErrNoRows {
		var sql string
		if indexUnique {
			sql = fmt.Sprintf("CREATE UNIQUE INDEX %q ON %q (%s);", indexName, indexTable, indexColumn)
		} else {
			sql = fmt.Sprintf("CREATE INDEX %q ON %q (%s);", indexName, indexTable, indexColumn)
		}

		_, err := DB.ExecContext(ctx, sql)
		if err != nil {
			logger.Printf("Database: error creating %s index: %v", indexName, err)
			return
		}
		logger.Printf("Database: %s index created", indexName)
	} else if err != nil {
		logger.Printf("Database: error checking for %s index: %v", indexName, err)
	} else {
		logger.Printf("Database: %s index already exists", indexName)
	}
}
