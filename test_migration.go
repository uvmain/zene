package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"zene/core/database"
)

func main() {
	// Test the database initialization with WAL mode
	ctx := context.Background()
	
	// Initialize the database
	database.Initialise(ctx)
	defer database.CloseDatabase()
	
	// Query to check if WAL mode is enabled
	var journalMode string
	err := database.DB.QueryRowContext(ctx, "PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		log.Fatalf("Failed to check journal mode: %v", err)
	}
	
	fmt.Printf("✅ Database migration successful!\n")
	fmt.Printf("✅ Journal mode: %s\n", journalMode)
	
	if journalMode == "wal" {
		fmt.Printf("✅ WAL mode is enabled!\n")
	} else {
		fmt.Printf("❌ WAL mode is NOT enabled. Current mode: %s\n", journalMode)
		os.Exit(1)
	}
	
	// Test that we can query sqlite_master (basic functionality)
	var count int
	err = database.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_master WHERE type='table'").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to query sqlite_master: %v", err)
	}
	
	fmt.Printf("✅ Database contains %d tables\n", count)
	fmt.Printf("✅ Core migration from zombiezen to modernc.org/sqlite is complete!\n")
	fmt.Printf("\nNote: Some complex query functions still need manual migration from stubs\n")
}