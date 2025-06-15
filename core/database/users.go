package database

import (
	"context"
	"fmt"
	"zene/core/logic"
	"zene/core/types"
)

func CreateUsersTable(ctx context.Context) {
	tableName := "users"
	schema := `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE
	);`
	createTable(ctx, tableName, schema)
}

func GetUserByUsername(ctx context.Context, username string) (types.User, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return types.User{}, fmt.Errorf("failed to take a db conn from the pool in GetUserByUsername: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, password_hash, created_at, is_admin FROM users WHERE username = $username`)
	defer stmt.Finalize()
	stmt.SetText("$username", username)

	if hasRow, err := stmt.Step(); err != nil {
		return types.User{}, fmt.Errorf("failed to select user from users: %v", err)
	} else if !hasRow {
		return types.User{}, fmt.Errorf("User not found")
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
		return types.User{}, fmt.Errorf("failed to take a db conn from the pool in GetUserByUsername: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, password_hash, created_at, is_admin FROM users WHERE id = $id`)
	defer stmt.Finalize()
	stmt.SetInt64("$id", id)

	if hasRow, err := stmt.Step(); err != nil {
		return types.User{}, fmt.Errorf("failed to select user from users: %v", err)
	} else if !hasRow {
		return types.User{}, fmt.Errorf("User not found")
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
		return []types.User{}, fmt.Errorf("failed to take a db conn from the pool in GetUserByUsername: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT id, username, created_at, is_admin, password_hash FROM users ORDER BY id ASC`)
	defer stmt.Finalize()

	var users []types.User
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []types.User{}, fmt.Errorf("failed to step through users: %w", err)
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

func UpsertUser(ctx context.Context, user types.User) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("failed to take a db conn from the pool in UpsertUser: %v", err)
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
		return fmt.Errorf("failed to upsert user: %v", err)
	}
	return nil
}

func DeleteUser(ctx context.Context, username string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	conn, err := DbPool.Take(ctx)
	if err != nil {
		return fmt.Errorf("failed to take a db conn from the pool in DeleteUser: %v", err)
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`DELETE FROM users WHERE username = $username`)
	defer stmt.Finalize()

	stmt.SetText("$username", username)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
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
