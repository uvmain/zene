package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

const (
	dbFile = "files.db"
)

func main() {
	// Get the directory to scan from the environment variable
	dirToScan := os.Getenv("SCAN_DIR")
	if dirToScan == "" {
		// dirToScan = "E:\\music"
		dirToScan = "Y:\\Music"
	}

	// Open SQLite database in WAL mode
	conn, err := sqlite.OpenConn(dbFile, sqlite.OpenReadWrite|sqlite.OpenCreate)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer conn.Close()
	// Enable WAL mode
	err = sqlitex.ExecuteTransient(conn, "PRAGMA journal_mode=WAL;", nil)
	if err != nil {
		log.Fatalf("Failed to set WAL mode: %v", err)
	}

	// Create the table if it doesn't exist
	err = createTable(conn)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Scan the directory and populate the database
	err = scanDirectory(conn, dirToScan)
	if err != nil {
		log.Fatalf("Failed to scan directory: %v", err)
	}

	// Set up HTTP handlers
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllFiles(conn, w)
	})
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		handleGetFileByName(conn, w, r)
	})

	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable(conn *sqlite.Conn) error {
	stmt := `CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dir_path TEXT NOT NULL,
		filename TEXT NOT NULL,
		date_added TEXT NOT NULL,
		mdate TEXT NOT NULL
	);`
	return sqlitex.ExecuteTransient(conn, stmt, nil)
}

func scanDirectory(conn *sqlite.Conn, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		stmt, err := conn.Prepare(`INSERT INTO files (dir_path, filename, date_added, mdate)
		VALUES ($dir_path, $filename, $date_added, $mdate);`)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %v", err)
		}
		defer stmt.Finalize()

		stmt.SetText("$dir_path", filepath.Dir(path))
		stmt.SetText("$filename", info.Name())
		stmt.SetText("$date_added", time.Now().Format(time.RFC3339))
		stmt.SetText("$mdate", info.ModTime().Format(time.RFC3339))

		// Execute the statement
		_, err = stmt.Step()
		if err != nil {
			return fmt.Errorf("failed to insert file: %v", err)
		}
		return nil
	})
}

type Row = struct {
	Id        string
	DirPath   string
	Filename  string
	DateAdded string
	Mdate     string
}

func handleGetAllFiles(conn *sqlite.Conn, w http.ResponseWriter) {
	stmt, _ := conn.Prepare(`SELECT dir_path, filename, date_added, mdate FROM files;`)

	var rows []Row

	for {
		if hasRow, err := stmt.Step(); err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		} else if !hasRow {
			break
		} else {
			var row Row
			row.Id = stmt.GetText("id")
			row.DirPath = stmt.GetText("dir_path")
			row.Filename = stmt.GetText("filename")
			row.DateAdded = stmt.GetText("date_added")
			row.Mdate = stmt.GetText("mdate")
			rows = append(rows, row)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rows); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func handleGetFileByName(conn *sqlite.Conn, w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing filename parameter", http.StatusBadRequest)
		return
	}

	stmt, _ := conn.Prepare(`SELECT id, dir_path, filename, date_added, mdate FROM files WHERE filename = $filename;`)
	stmt.SetText("$filename", filename)

	if hasRow, err := stmt.Step(); err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	} else if !hasRow {
		http.Error(w, "File not found", http.StatusNotFound)
	} else {
		var row Row
		row.Id = stmt.GetText("id")
		row.DirPath = stmt.GetText("dir_path")
		row.Filename = stmt.GetText("filename")
		row.DateAdded = stmt.GetText("date_added")
		row.Mdate = stmt.GetText("mdate")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(row); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
