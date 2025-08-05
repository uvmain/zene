package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/config"
	"zene/core/encryption"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func createUsersTable(ctx context.Context) error {
	tableName := "users"
	schema := `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    encrypted_password TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE
	);`
	err := createTable(ctx, tableName, schema)
	return err
}

func CreateAdminUserIfRequired(ctx context.Context) error {
	adminUserExists, err := anyAdminUsersExist(ctx)
	if err != nil {
		return fmt.Errorf("checking if any admin users exist: %v", err)
	}

	if adminUserExists {
		logger.Printf("Admin user already exists, skipping creation")
		return nil
	}

	adminUsername := config.AdminUsername
	adminPassword := config.AdminPassword

	if adminUsername == "" {
		adminUsername = "admin"
		logger.Println("Admin username not set in configuration, using default 'admin'")
	}

	if adminPassword == "" {
		logger.Logger.Println("admin password not set in configuration, generating a random one")
		adminPassword, err = logic.GenerateRandomPassword(12)
		if err != nil {
			logger.Printf("Error generating random password for admin user: %v", err)
			return fmt.Errorf("generating random password for admin user: %v", err)
		} else {
			logger.Printf("** Generated random password for admin user: %s", adminPassword)
		}
	}

	encryptedPassword, err := encryption.EncryptAES(adminPassword)
	if err != nil {
		logger.Printf("Error encrypting admin password: %v", err)
		return fmt.Errorf("encrypting admin password: %v", err)
	}

	user := types.User{
		Username:          adminUsername,
		EncryptedPassword: encryptedPassword,
		IsAdmin:           true,
	}

	_, err = UpsertUser(ctx, user)
	if err != nil {
		logger.Printf("Error upserting admin user: %v", err)
		return fmt.Errorf("upserting admin user: %v", err)
	}

	logger.Printf("Admin user %s created successfully", adminUsername)
	return nil
}

func GetUserByContext(ctx context.Context) (types.User, error) {
	val := ctx.Value("userId")
	userId, ok := val.(int64)
	if !ok {
		return types.User{}, fmt.Errorf("userId missing or invalid in context")
	}

	user, err := GetUserById(ctx, userId)
	if err != nil {
		return types.User{}, fmt.Errorf("failed to get user from database: %v", err)
	}

	return user, nil
}

func GetUserByUsername(ctx context.Context, username string) (types.User, error) {
	query := `SELECT id, username, encrypted_password, created_at, is_admin FROM users WHERE username = ?`
	var row types.User

	err := DB.QueryRowContext(ctx, query, username).Scan(&row.Id, &row.Username, &row.EncryptedPassword, &row.CreatedAt, &row.IsAdmin)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	}
	return row, nil
}

func GetUserById(ctx context.Context, id int64) (types.User, error) {
	query := `SELECT id, username, encrypted_password, created_at, is_admin FROM users WHERE id = ?`
	var row types.User

	err := DB.QueryRowContext(ctx, query, id).Scan(&row.Id, &row.Username, &row.EncryptedPassword, &row.CreatedAt, &row.IsAdmin)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	}
	return row, nil
}

func GetAllUsers(ctx context.Context) ([]types.User, error) {
	query := `SELECT id, username, created_at, is_admin, encrypted_password FROM users ORDER BY id ASC`
	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return []types.User{}, fmt.Errorf("querying all users: %v", err)
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var row types.User
		err := rows.Scan(&row.Id, &row.Username, &row.CreatedAt, &row.IsAdmin, &row.EncryptedPassword)
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
		INSERT INTO users (username, encrypted_password, created_at, is_admin)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(username) DO UPDATE SET
			encrypted_password = excluded.encrypted_password,
			is_admin = excluded.is_admin`

	result, err := DB.ExecContext(ctx, query, user.Username, user.EncryptedPassword, logic.GetCurrentTimeFormatted(), user.IsAdmin)
	if err != nil {
		return 0, fmt.Errorf("upserting user: %v", err)
	}

	rowID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert ID: %v", err)
	}

	adminUser, err := logic.GetUsernameFromContext(ctx)
	if err != nil {
		logger.Printf("user %s created, failed to get name of the creating admin user: %v", user.Username, err)
	} else {
		logger.Printf("user %s created by admin user %s", user.Username, adminUser)
	}

	return rowID, nil
}

func DeleteUserByUsername(ctx context.Context, username string) error {
	query := `DELETE FROM users WHERE username = ?`
	_, err := DB.ExecContext(ctx, query, username)
	if err != nil {
		return fmt.Errorf("deleting user: %v", err)
	}

	adminUser, err := logic.GetUsernameFromContext(ctx)
	if err != nil {
		logger.Printf("user %s deleted, failed to get admin user: %v", username, err)
	} else {
		logger.Printf("user %s deleted by admin user %s", username, adminUser)
	}

	return nil
}

func DeleteUserById(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting user: %v", err)
	}

	adminUser, err := logic.GetUsernameFromContext(ctx)
	if err != nil {
		logger.Printf("user id %d deleted, failed to get admin user: %v", id, err)
	} else {
		logger.Printf("user id %d deleted by admin user %s", id, adminUser)
	}

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

func anyAdminUsersExist(ctx context.Context) (bool, error) {
	query := `SELECT 1 FROM users where is_admin = true LIMIT 1`
	var exists int
	err := DB.QueryRowContext(ctx, query).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func GetEncryptedPasswordFromDB(ctx context.Context, username string) (string, int64, error) {
	query := `SELECT encrypted_password, id FROM users WHERE username = ?`
	var encryptedPassword string
	var userId int64

	err := DB.QueryRowContext(ctx, query, username).Scan(&encryptedPassword, &userId)
	if err == sql.ErrNoRows {
		return "", 0, fmt.Errorf("user not found")
	} else if err != nil {
		return "", 0, fmt.Errorf("selecting encrypted password for user %s: %v", username, err)
	}

	return encryptedPassword, userId, nil
}
