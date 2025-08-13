package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
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

func doesViewExist(ctx context.Context, viewName string) (bool, error) {
	query := `SELECT name FROM sqlite_master WHERE type = 'view' AND name = ?`
	var name string
	err := DB.QueryRowContext(ctx, query, viewName).Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("checking if view exists: %v", err)
	}
	return true, nil
}

func getTableNameFromSchema(schema string) (string, error) {
	re := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)`)
	matches := re.FindStringSubmatch(schema)
	if len(matches) < 2 {
		logger.Printf("table name not found in schema: %v", schema)
		return "", fmt.Errorf("table name not found in schema")
	}
	return matches[1], nil
}

func getVirtualTableNameFromSchema(schema string) (string, error) {
	re := regexp.MustCompile(`(?i)CREATE\s+VIRTUAL\s+TABLE\s+(\w+)`)
	matches := re.FindStringSubmatch(schema)
	if len(matches) < 2 {
		logger.Printf("virtual table name not found in schema: %v", schema)
		return "", fmt.Errorf("virtual table name not found in schema")
	}
	return matches[1], nil
}

func getViewNameFromSchema(schema string) (string, error) {
	re := regexp.MustCompile(`(?i)CREATE\s+VIEW\s+(\w+)`)
	matches := re.FindStringSubmatch(schema)
	if len(matches) < 2 {
		logger.Println("view name not found in schema")
		return "", fmt.Errorf("view name not found in schema")
	}
	return matches[1], nil
}

func getTriggerNameFromSchema(schema string) (string, error) {
	re := regexp.MustCompile(`(?i)CREATE\s+TRIGGER\s+(\w+)`)
	matches := re.FindStringSubmatch(schema)
	if len(matches) < 2 {
		logger.Println("trigger name not found in schema")
		return "", fmt.Errorf("trigger name not found in schema")
	}
	return matches[1], nil
}

func createTable(ctx context.Context, schema string) {
	tableName, err := getTableNameFromSchema(schema)
	if err != nil {
		log.Fatalf("Error extracting table name from schema: %v", err)
	}
	tableExists, err := doesTableExist(ctx, tableName)
	if err != nil {
		log.Fatalf("Error checking if table %s exists: %v", tableName, err)
	}

	if !tableExists {
		_, err := DB.ExecContext(ctx, schema)
		if err != nil {
			log.Fatalf("Database: error creating %s table: %v", tableName, err)
		}
		logger.Printf("Database: %s table created", tableName)
	} else {
		logger.Printf("Database: %s table already exists", tableName)
	}
}

func createVirtualTable(ctx context.Context, schema string) {
	tableName, err := getVirtualTableNameFromSchema(schema)
	if err != nil {
		log.Fatalf("Error extracting table name from schema: %v", err)
	}
	tableExists, err := doesTableExist(ctx, tableName)
	if err != nil {
		log.Fatalf("Error checking if virtual table %s exists: %v", tableName, err)
	}

	if !tableExists {
		_, err := DB.ExecContext(ctx, schema)
		if err != nil {
			log.Fatalf("Database: error creating %s virtual table: %v", tableName, err)
		}
		logger.Printf("Database: %s virtual table created", tableName)
	} else {
		logger.Printf("Database: %s virtual table already exists", tableName)
	}
}

func createView(ctx context.Context, schema string) {
	viewName, err := getViewNameFromSchema(schema)
	if err != nil {
		log.Fatalf("Error extracting view name from schema: %v", err)
	}
	viewExists, err := doesViewExist(ctx, viewName)
	if err != nil {
		log.Fatalf("Error checking if view %v exists: %v", viewExists, err)
	}

	if !viewExists {
		_, err := DB.ExecContext(ctx, schema)
		if err != nil {
			log.Fatalf("Database: error creating %s view: %v", viewName, err)
		}
		logger.Printf("Database: %s view created", viewName)
	} else {
		logger.Printf("Database: %s view already exists", viewName)
	}
}

func createTrigger(ctx context.Context, schema string) {
	triggerName, err := getTriggerNameFromSchema(schema)
	if err != nil {
		log.Fatalf("Error extracting trigger name from schema: %v", err)
	}

	query := "SELECT name FROM sqlite_master WHERE type='trigger' AND name=?"
	var name string
	err = DB.QueryRowContext(ctx, query, triggerName).Scan(&name)

	if err == sql.ErrNoRows {
		_, err := DB.ExecContext(ctx, schema)
		if err != nil {
			log.Fatalf("Database: error creating %s trigger: %v", triggerName, err)
			return
		}
		logger.Printf("Database: %s trigger created", triggerName)
	} else if err != nil {
		log.Fatalf("Database: error checking for %s trigger: %v", triggerName, err)
	} else {
		logger.Printf("Database: %s trigger already exists", triggerName)
	}
}

func createIndex(ctx context.Context, indexName, indexTable string, indexColumns []string, indexUnique bool) {
	query := "SELECT name FROM sqlite_master WHERE type='index' AND name=?"
	var name string
	err := DB.QueryRowContext(ctx, query, indexName).Scan(&name)

	if err == sql.ErrNoRows {
		var sql string
		if indexUnique {
			sql = fmt.Sprintf("CREATE UNIQUE INDEX %q ON %q (%s);", indexName, indexTable, strings.Join(indexColumns, ", "))
		} else {
			sql = fmt.Sprintf("CREATE INDEX %q ON %q (%s);", indexName, indexTable, strings.Join(indexColumns, ", "))
		}

		_, err := DB.ExecContext(ctx, sql)
		if err != nil {
			log.Fatalf("Database: error creating %s index: %v", indexName, err)
			return
		}
		logger.Printf("Database: %s index created", indexName)
	} else if err != nil {
		log.Fatalf("Database: error checking for %s index: %v", indexName, err)
	} else {
		logger.Printf("Database: %s index already exists", indexName)
	}
}
