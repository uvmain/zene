package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/logger"
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
	query := `SELECT id, username, password_hash, created_at, is_admin FROM users WHERE username = ?`
	var row types.User

	err := DB.QueryRowContext(ctx, query, username).Scan(&row.Id, &row.Username, &row.PasswordHash, &row.CreatedAt, &row.IsAdmin)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	}
	return row, nil
}

func GetUserById(ctx context.Context, id int64) (types.User, error) {
	query := `SELECT id, username, password_hash, created_at, is_admin FROM users WHERE id = ?`
	var row types.User

	err := DB.QueryRowContext(ctx, query, id).Scan(&row.Id, &row.Username, &row.PasswordHash, &row.CreatedAt, &row.IsAdmin)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	}
	return row, nil
}

func GetAllUsers(ctx context.Context) ([]types.User, error) {
	query := `SELECT id, username, created_at, is_admin, password_hash FROM users ORDER BY id ASC`
	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return []types.User{}, fmt.Errorf("querying all users: %v", err)
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var row types.User
		err := rows.Scan(&row.Id, &row.Username, &row.CreatedAt, &row.IsAdmin, &row.PasswordHash)
		if err != nil {
			return []types.User{}, fmt.Errorf("scanning user row: %v", err)
		}
		users = append(users, row)
	}

	if err := rows.Err(); err != nil {
		return []types.User{}, fmt.Errorf("rows error: %v", err)
	}

	return users, nil
}

func UpsertUser(ctx context.Context, user types.User) (int64, error) {
	query := `
		INSERT INTO users (username, password_hash, is_admin)
		VALUES (?, ?, ?)
		ON CONFLICT(username) DO UPDATE SET
			password_hash = excluded.password_hash,
			is_admin = excluded.is_admin`

	result, err := DB.ExecContext(ctx, query, user.Username, user.PasswordHash, user.IsAdmin)
	if err != nil {
		return 0, fmt.Errorf("upserting user: %v", err)
	}

	rowID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert ID: %v", err)
	}

	return rowID, nil
}

func DeleteUserByUsername(ctx context.Context, username string) error {
	query := `DELETE FROM users WHERE username = ?`
	_, err := DB.ExecContext(ctx, query, username)
	if err != nil {
		return fmt.Errorf("deleting user: %v", err)
	}
	logger.Printf("Deleted user with username: %s", username)
	return nil
}

func DeleteUserById(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting user: %v", err)
	}
	logger.Printf("Deleted user with ID: %d", id)
	return nil
}

func AnyUsersExist(ctx context.Context) (bool, error) {
	query := `SELECT 1 FROM users LIMIT 1`
	var exists int
	err := DB.QueryRowContext(ctx, query).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
