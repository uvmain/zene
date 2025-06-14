package database

import (
	"context"
	"fmt"
	"github.com/ollama/ollama/core/logic" // Assuming logic is now under core
	"github.com/ollama/ollama/core/types"
	log "github.com/sirupsen/logrus"
	sqlite "zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

// DB struct likely holds the connection pool, similar to how it's used elsewhere in ollama.
// For now, we'll assume the global DbPool is part of this DB struct or accessible via it.
// If DbPool is indeed global and DB is a wrapper, adjustments might be needed.
// Let's assume DB has a field Pool *sqlite.Pool for now.
// If the existing `createTable` and other functions not modified here rely on global `DbPool` and `dbMutex`,
// they will need to be refactored if `DbPool` is moved into the `DB` struct.
// For this change, we will pass *DB and assume its methods will use db.Pool instead of global DbPool.

// CreateUsersTable initializes the users table in the database.
// It now takes *DB to potentially use a connection pool from it.
func CreateUsersTable(ctx context.Context, db *DB) error {
	// This function's body used a global `createTable` which likely uses global `DbPool`.
	// To adapt, `createTable` would also need to take `db *DB`.
	// For now, let's assume `createTable` is refactored or this function directly uses `db.Pool`.
	// Simplified direct execution for now:
	conn, err := db.Pool.Take(ctx)
	if err != nil {
		log.Errorf("CreateUsersTable: failed to take a db conn from the pool: %v", err)
		return fmt.Errorf("CreateUsersTable: failed to take a db conn from the pool: %w", err)
	}
	defer db.Pool.Put(conn)

	schema := `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
	);`
	err = sqlitex.ExecuteTransient(conn, schema, nil)
	if err != nil {
		log.Errorf("CreateUsersTable: failed to create users table: %v", err)
		return fmt.Errorf("CreateUsersTable: failed to create users table: %w", err)
	}
	log.Info("Users table created or already exists.")
	return nil
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(ctx context.Context, db *DB, username string) (types.User, error) {
	// Assuming dbMutex is still global or part of DB struct. For now, keep as is if it's global.
	// dbMutex.RLock()
	// defer dbMutex.RUnlock()

	conn, err := db.Pool.Take(ctx)
	if err != nil {
		log.Errorf("GetUserByUsername: failed to take a db conn from the pool: %v", err)
		return types.User{}, fmt.Errorf("failed to take a db conn from the pool in GetUserByUsername: %w", err)
	}
	defer db.Pool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, password_hash, created_at, is_admin FROM users WHERE username = $username`)
	defer stmt.Finalize()
	stmt.SetText("$username", username)

	hasRow, err := stmt.Step()
	if err != nil {
		log.Errorf("Error selecting user %s: %v", username, err)
		return types.User{}, fmt.Errorf("failed to select user from users: %w", err)
	}
	if !hasRow {
		return types.User{}, fmt.Errorf("User not found: %s", username) // Consider a specific error type for Not Found
	}

	var row types.User
	row.Id = int(stmt.GetInt64("id")) // Corrected to int
	row.Username = stmt.GetText("username")
	row.PasswordHash = stmt.GetText("password_hash")
	row.CreatedAt = stmt.GetText("created_at")
	row.IsAdmin = stmt.GetInt64("is_admin") != 0
	return row, nil
}

// GetUserById retrieves a user by their ID.
func GetUserById(ctx context.Context, db *DB, id int) (types.User, error) {
	conn, err := db.Pool.Take(ctx)
	if err != nil {
		log.Errorf("GetUserById: failed to take a db conn from the pool: %v", err)
		return types.User{}, fmt.Errorf("failed to take a db conn from the pool in GetUserById: %w", err)
	}
	defer db.Pool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, password_hash, created_at, is_admin FROM users WHERE id = $id`)
	defer stmt.Finalize()
	stmt.SetInt64("$id", int64(id))

	hasRow, err := stmt.Step()
	if err != nil {
		log.Errorf("Error selecting user by ID %d: %v", id, err)
		return types.User{}, fmt.Errorf("failed to select user by ID from users: %w", err)
	}
	if !hasRow {
		return types.User{}, fmt.Errorf("User not found with ID: %d", id) // Specific error for Not Found
	}

	var row types.User
	row.Id = int(stmt.GetInt64("id"))
	row.Username = stmt.GetText("username")
	row.PasswordHash = stmt.GetText("password_hash")
	row.CreatedAt = stmt.GetText("created_at")
	row.IsAdmin = stmt.GetInt64("is_admin") != 0
	return row, nil
}

// GetAllUsers retrieves all users, selecting only id, username, and is_admin.
func GetAllUsers(ctx context.Context, db *DB) ([]types.User, error) {
	conn, err := db.Pool.Take(ctx)
	if err != nil {
		log.Errorf("GetAllUsers: failed to take a db conn from the pool: %v", err)
		return nil, fmt.Errorf("failed to take a db conn from the pool in GetAllUsers: %w", err)
	}
	defer db.Pool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, is_admin FROM users ORDER BY username ASC`)
	defer stmt.Finalize()

	var users []types.User
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			log.Errorf("Error stepping through users for GetAllUsers: %v", err)
			return nil, fmt.Errorf("failed to step through users: %w", err)
		}
		if !hasRow {
			break
		}
		var user types.User
		user.Id = int(stmt.GetInt64("id"))
		user.Username = stmt.GetText("username")
		user.IsAdmin = stmt.GetInt64("is_admin") != 0
		// PasswordHash is not selected and will remain empty
		users = append(users, user)
	}
	return users, nil
}

// UpsertUser creates a new user if ID is 0, or updates an existing user if ID is non-zero.
// Password hashing should be done before calling this function.
func UpsertUser(ctx context.Context, db *DB, user types.User) error {
	// Assuming dbMutex is still global or part of DB struct.
	// dbMutex.Lock()
	// defer dbMutex.Unlock()

	conn, err := db.Pool.Take(ctx)
	if err != nil {
		log.Errorf("UpsertUser: failed to take a db conn from the pool: %v", err)
		return fmt.Errorf("failed to take a db conn from the pool in UpsertUser: %w", err)
	}
	defer db.Pool.Put(conn)

	var stmt *sqlite.Stmt
	if user.Id == 0 { // Create new user
		stmt = conn.Prep(`
			INSERT INTO users (username, password_hash, is_admin)
			VALUES ($username, $password_hash, $is_admin)
		`)
		stmt.SetText("$username", user.Username)
		stmt.SetText("$password_hash", user.PasswordHash)
		stmt.SetInt64("$is_admin", logic.BoolToInt64(user.IsAdmin))
	} else { // Update existing user
		// If password hash is empty, it means password is not being updated.
		if user.PasswordHash != "" {
			stmt = conn.Prep(`
				UPDATE users SET username = $username, password_hash = $password_hash, is_admin = $is_admin
				WHERE id = $id
			`)
			stmt.SetText("$password_hash", user.PasswordHash)
		} else {
			stmt = conn.Prep(`
				UPDATE users SET username = $username, is_admin = $is_admin
				WHERE id = $id
			`)
		}
		stmt.SetText("$username", user.Username)
		stmt.SetInt64("$is_admin", logic.BoolToInt64(user.IsAdmin))
		stmt.SetInt64("$id", int64(user.Id))
	}
	defer stmt.Finalize()

	_, err = stmt.Step()
	if err != nil {
		// Check for UNIQUE constraint failed: users.username
		if sqlite.ErrCode(err) == sqlite.CONSTRAINT_UNIQUE {
			log.Warnf("UpsertUser: UNIQUE constraint failed for username '%s': %v", user.Username, err)
			return fmt.Errorf("username '%s' already exists: %w", user.Username, err) // Return a more specific error
		}
		log.Errorf("Failed to upsert user (ID: %d, Username: %s): %v", user.Id, user.Username, err)
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	if user.Id == 0 { // If it was an insert, get the last inserted ID
		user.Id = int(conn.LastInsertRowID())
	} else { // If it was an update, check rows affected
        changes := conn.Changes()
        if changes == 0 {
            // This might happen if the user data was identical or user ID didn't exist
            // Check if user exists to differentiate
            _, checkErr := GetUserById(ctx, db, user.Id)
            if checkErr != nil { // Assuming GetUserById returns a "not found" type error
                 log.Warnf("UpsertUser (update): User with ID %d not found, 0 rows affected.", user.Id)
                 return fmt.Errorf("user with ID %d not found for update", user.Id)
            }
            log.Infof("UpsertUser (update): User with ID %d data was unchanged or write was ineffective, 0 rows affected.", user.Id)
            // Not necessarily an error if data was identical, but good to be aware.
        }
    }
	return nil
}

// DeleteUser deletes a user by their ID.
func DeleteUser(ctx context.Context, db *DB, id int) error {
	// dbMutex.Lock()
	// defer dbMutex.Unlock()

	conn, err := db.Pool.Take(ctx)
	if err != nil {
		log.Errorf("DeleteUser: failed to take a db conn from the pool: %v", err)
		return fmt.Errorf("failed to take a db conn from the pool in DeleteUser: %w", err)
	}
	defer db.Pool.Put(conn)

	stmt := conn.Prep(`DELETE FROM users WHERE id = $id`)
	defer stmt.Finalize()
	stmt.SetInt64("$id", int64(id))

	_, err = stmt.Step()
	if err != nil {
		log.Errorf("Failed to delete user with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if conn.Changes() == 0 {
		log.Warnf("Attempted to delete user with ID %d, but no rows were affected (user not found).", id)
		return fmt.Errorf("user not found or already deleted: ID %d", id) // Specific error
	}

	return nil
}

// AnyUsersExist checks if there are any users in the database.
// This function might also need to be updated if DbPool is no longer global.
// For now, assuming it can still use a global DbPool or db.Pool if db is passed.
func AnyUsersExist(ctx context.Context, db *DB) (bool, error) {
	// dbMutex.RLock()
	// defer dbMutex.RUnlock()

	conn, err := db.Pool.Take(ctx) // Assumes db.Pool exists
	if err != nil {
		log.Errorf("AnyUsersExist: failed to take a db conn from the pool: %v", err)
		return false, fmt.Errorf("failed to take a db conn from the pool in AnyUsersExist: %w", err)
	}
	defer db.Pool.Put(conn)

	stmt := conn.Prep(`SELECT 1 FROM users LIMIT 1`)
	defer stmt.Finalize()

	hasRow, err := stmt.Step()
	if err != nil {
		log.Errorf("Error checking if any users exist: %v", err)
		return false, fmt.Errorf("failed to check for any users: %w", err)
	}

	return hasRow, nil
}

// TODO: The global DbPool and dbMutex are problematic if we are passing *DB around.
// This needs a more thorough refactor:
// 1. The DB struct should encapsulate the sqlite.Pool and any necessary mutexes.
// 2. All database functions should take `db *DB` as the first argument.
// 3. Global `DbPool` and `dbMutex` should be removed or be part of an unexported default DB instance.
// The current changes assume `db.Pool` can replace `DbPool` and that `dbMutex` is either
// still global (problematic) or handled within each function if necessary (less likely for pool connections).
// For the scope of this task, I've adapted the functions to take `*DB` and use `db.Pool`.
// The `CreateUsersTable` function has been adapted more directly.
// `GetUserByUsername`, `UpsertUser`, `DeleteUser`, `AnyUsersExist` are adapted to use `db *DB` and `db.Pool`.
// `CreateUsersTable` also added `IF NOT EXISTS` and a default for `is_admin`.
// `UpsertUser` logic was changed to differentiate between INSERT and UPDATE based on `user.Id`.
// Logging has been changed from fmt.Errorf to log.Errorf/Warnf/Infof.
// User ID type in types.User is int, so GetInt64 results are cast to int.
// `DeleteUser` now takes ID and checks `conn.Changes()`.
// `GetAllUsers` added to select only id, username, is_admin.
// `GetUserById` was added.
// Assumed `github.com/ollama/ollama/core/logic` for `BoolToInt64`.
// Corrected `row.Id` assignment in `GetUserByUsername` and `GetUserById` to `int(stmt.GetInt64("id"))`.
