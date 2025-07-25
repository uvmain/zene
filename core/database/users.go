package database

import (
	"context"
	"fmt"
	"zene/core/logic"
	"zene/core/types"
)

func createUsersTable(ctx context.Context) error {
	tableName := "users"
	schema := `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE
	);`
	err := createTable(ctx, tableName, schema)
	return err
}

func GetUserByUsername(ctx context.Context, username string) (types.User, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return types.User{}, fmt.Errorf("taking a db conn from the pool in GetUserByUsername: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, password_hash, created_at, is_admin FROM users WHERE username = $username`)
	defer stmt.Finalize()
	stmt.SetText("$username", username)

	if hasRow, err := stmt.Step(); err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	} else if !hasRow {
		return types.User{}, fmt.Errorf("user not found")
	} else {
		var row types.User
		row.Id = stmt.GetInt64("id")
		row.Username = stmt.GetText("username")
		row.CreatedAt = stmt.GetText("created_at")
		row.IsAdmin = stmt.GetInt64("is_admin") != 0
		row.PasswordHash = stmt.GetText("password_hash")
		return row, nil
	}
}

func GetUserById(ctx context.Context, id int64) (types.User, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return types.User{}, fmt.Errorf("taking a db conn from the pool in GetUserById: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, password_hash, created_at, is_admin FROM users WHERE id = $id`)
	defer stmt.Finalize()
	stmt.SetInt64("$id", id)

	if hasRow, err := stmt.Step(); err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	} else if !hasRow {
		return types.User{}, fmt.Errorf("user not found")
	} else {
		var row types.User
		row.Id = stmt.GetInt64("id")
		row.Username = stmt.GetText("username")
		row.CreatedAt = stmt.GetText("created_at")
		row.IsAdmin = stmt.GetInt64("is_admin") != 0
		row.PasswordHash = stmt.GetText("password_hash")
		return row, nil
	}
}

func GetAllUsers(ctx context.Context) ([]types.User, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return []types.User{}, fmt.Errorf("taking a db conn from the pool in GetAllUsers: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, created_at, is_admin, password_hash FROM users ORDER BY id ASC`)
	defer stmt.Finalize()

	var users []types.User
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []types.User{}, fmt.Errorf("stepping through users: %w", err)
		}
		if !hasRow {
			break
		}
		var row types.User
		row.Id = stmt.GetInt64("id")
		row.Username = stmt.GetText("username")
		row.CreatedAt = stmt.GetText("created_at")
		row.IsAdmin = stmt.GetInt64("is_admin") != 0
		row.PasswordHash = stmt.GetText("password_hash")
		users = append(users, row)
	}
	return users, nil
}

func UpsertUser(ctx context.Context, user types.User) (int64, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	rowId := int64(0)

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return rowId, fmt.Errorf("taking a db conn from the pool in UpsertUser: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`
		INSERT INTO users (username, password_hash, is_admin)
		VALUES ($username, $password_hash, $is_admin)
		ON CONFLICT(username) DO UPDATE SET
			password_hash = excluded.password_hash,
			is_admin = excluded.is_admin
	`)
	defer stmt.Finalize()

	stmt.SetText("$username", user.Username)
	stmt.SetText("$password_hash", user.PasswordHash)
	stmt.SetInt64("$is_admin", logic.BoolToInt64(user.IsAdmin))

	_, err = stmt.Step()
	if err != nil {
		return rowId, fmt.Errorf("upserting user: %v", err)
	}
	rowID := conn.LastInsertRowID()
	return rowID, nil
}

func DeleteUserByUsername(ctx context.Context, username string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in DeleteUserByUsername: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM users WHERE username = $username`)
	defer stmt.Finalize()

	stmt.SetText("$username", username)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("deleting user: %v", err)
	}

	return nil
}

func DeleteUserById(ctx context.Context, id int64) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("taking a db conn from the pool in DeleteUserById: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM users WHERE id = $id`)
	defer stmt.Finalize()

	stmt.SetInt64("$id", id)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("deleting user: %v", err)
	}

	return nil
}

func AnyUsersExist(ctx context.Context) (bool, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return false, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT 1 FROM users LIMIT 1`)
	defer stmt.Finalize()

	hasRow, err := stmt.Step()
	if err != nil {
		return false, err
	}

	return hasRow, nil
}
