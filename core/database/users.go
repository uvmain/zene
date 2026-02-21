package database

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"zene/core/config"
	"zene/core/encryption"
	"zene/core/logger"
	"zene/core/logic"
	"zene/core/types"
)

func migrateUsers(ctx context.Context) {
	schema := `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT NOT NULL,
		scrobbling_enabled BOOLEAN NOT NULL DEFAULT 1,
    ldap_authenticated BOOLEAN NOT NULL DEFAULT 0,
    admin_role BOOLEAN NOT NULL DEFAULT 0,
    settings_role BOOLEAN NOT NULL DEFAULT 1,
    stream_role BOOLEAN NOT NULL DEFAULT 1,
    jukebox_role BOOLEAN NOT NULL DEFAULT 0,
    download_role BOOLEAN NOT NULL DEFAULT 0,
    upload_role BOOLEAN NOT NULL DEFAULT 0,
    playlist_role BOOLEAN NOT NULL DEFAULT 0,
    cover_art_role BOOLEAN NOT NULL DEFAULT 0,
    comment_role BOOLEAN NOT NULL DEFAULT 0,
    podcast_role BOOLEAN NOT NULL DEFAULT 0,
    share_role BOOLEAN NOT NULL DEFAULT 0,
    video_conversion_role BOOLEAN NOT NULL DEFAULT 0,
		max_bit_rate INTEGER NOT NULL DEFAULT 0
	);`
	createTable(ctx, schema)

	schema = `CREATE TABLE user_music_folders (
    user_id INTEGER NOT NULL,
    folder_id TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		FOREIGN KEY (folder_id) REFERENCES music_folders(id) ON DELETE CASCADE
		UNIQUE(user_id,folder_id)
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_user_music_folders_user", "user_music_folders", []string{"user_id"}, false)

	schema = `CREATE VIEW users_with_folders AS
		SELECT u.id AS user_id, u.username, u.email, u.password, u.scrobbling_enabled, u.ldap_authenticated, u.admin_role, u.settings_role,
			u.stream_role, u.jukebox_role, u.download_role, u.upload_role, u.playlist_role, u.cover_art_role,
			u.comment_role, u.podcast_role, u.share_role, u.video_conversion_role, u.max_bit_rate, GROUP_CONCAT(f.folder_id) AS music_folder_ids
		FROM users u
		LEFT JOIN user_music_folders f
		ON u.id = f.user_id
		GROUP BY u.id;`
	createView(ctx, schema)
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
	adminEmail := cmp.Or(config.AdminEmail, "admin@localhost")

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

	musicDirs, err := GetMusicFolders(ctx)
	if err != nil {
		logger.Printf("Error getting music folders: %v", err)
		return fmt.Errorf("getting music folders: %v", err)
	}

	var folderIds []int
	for _, folder := range musicDirs {
		folderIds = append(folderIds, folder.Id)
	}

	user := types.User{
		Username:            adminUsername,
		Email:               adminEmail,
		Password:            encryptedPassword,
		AdminRole:           true,
		ScrobblingEnabled:   true,
		SettingsRole:        true,
		StreamRole:          true,
		JukeboxRole:         true,
		DownloadRole:        true,
		UploadRole:          true,
		PlaylistRole:        true,
		CoverArtRole:        true,
		CommentRole:         true,
		PodcastRole:         true,
		ShareRole:           true,
		VideoConversionRole: true,
		MaxBitRate:          0,
		Folders:             folderIds,
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
	val := ctx.Value(types.ContextKey("userId"))
	userId, ok := val.(int)
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
	query := `SELECT * FROM users_with_folders where username = ?;`
	var row types.User
	var foldersString string

	err := DB.QueryRowContext(ctx, query, username).Scan(&row.Id, &row.Username, &row.Email, &row.Password, &row.ScrobblingEnabled, &row.LdapAuthenticated,
		&row.AdminRole, &row.SettingsRole, &row.StreamRole, &row.JukeboxRole, &row.DownloadRole, &row.UploadRole, &row.PlaylistRole,
		&row.CoverArtRole, &row.CommentRole, &row.PodcastRole, &row.ShareRole, &row.VideoConversionRole, &row.MaxBitRate, &foldersString)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	}
	if foldersString == "" {
		row.Folders = []int{}
	} else {
		row.Folders = logic.StringToIntSlice(foldersString)
	}
	return row, nil
}

func GetUserWithoutFoldersByUsername(ctx context.Context, username string) (types.User, error) {
	query := `SELECT * FROM users where username = ?;`
	var row types.User

	err := DB.QueryRowContext(ctx, query, username).Scan(&row.Id, &row.Username, &row.Email, &row.Password, &row.ScrobblingEnabled, &row.LdapAuthenticated,
		&row.AdminRole, &row.SettingsRole, &row.StreamRole, &row.JukeboxRole, &row.DownloadRole, &row.UploadRole, &row.PlaylistRole,
		&row.CoverArtRole, &row.CommentRole, &row.PodcastRole, &row.ShareRole, &row.VideoConversionRole, &row.MaxBitRate)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	}
	return row, nil
}

func GetUserById(ctx context.Context, id int) (types.User, error) {
	query := `SELECT * FROM users_with_folders WHERE user_id = ?;`
	var row types.User
	var foldersString string

	err := DB.QueryRowContext(ctx, query, id).Scan(&row.Id, &row.Username, &row.Email, &row.Password, &row.ScrobblingEnabled, &row.LdapAuthenticated,
		&row.AdminRole, &row.SettingsRole, &row.StreamRole, &row.JukeboxRole, &row.DownloadRole, &row.UploadRole, &row.PlaylistRole,
		&row.CoverArtRole, &row.CommentRole, &row.PodcastRole, &row.ShareRole, &row.VideoConversionRole, &row.MaxBitRate, &foldersString)
	if err == sql.ErrNoRows {
		return types.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return types.User{}, fmt.Errorf("selecting user from users: %v", err)
	}
	row.Folders = logic.StringToIntSlice(foldersString)
	return row, nil
}

func GetAllUsers(ctx context.Context) ([]types.User, error) {
	query := `SELECT * FROM users_with_folders order by user_id asc;`
	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return []types.User{}, fmt.Errorf("querying all users: %v", err)
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var row types.User
		var foldersString string
		err := rows.Scan(&row.Id, &row.Username, &row.Email, &row.Password, &row.ScrobblingEnabled, &row.LdapAuthenticated,
			&row.AdminRole, &row.SettingsRole, &row.StreamRole, &row.JukeboxRole, &row.DownloadRole, &row.UploadRole, &row.PlaylistRole,
			&row.CoverArtRole, &row.CommentRole, &row.PodcastRole, &row.ShareRole, &row.VideoConversionRole, &row.MaxBitRate, &foldersString)
		if err != nil {
			return []types.User{}, fmt.Errorf("scanning user row: %v", err)
		}
		row.Folders = logic.StringToIntSlice(foldersString)
		users = append(users, row)
	}

	if err := rows.Err(); err != nil {
		return []types.User{}, fmt.Errorf("rows error: %v", err)
	}

	return users, nil
}

func UpsertUser(ctx context.Context, user types.User) (int, error) {
	query := `INSERT INTO users (username, password, email, scrobbling_enabled, ldap_authenticated, admin_role, settings_role,
    	stream_role, jukebox_role, download_role, upload_role, playlist_role, cover_art_role, comment_role,
    	podcast_role, share_role, video_conversion_role, max_bit_rate)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(username) DO UPDATE SET
			password = excluded.password,
			email = excluded.email,
			ldap_authenticated = excluded.ldap_authenticated,
			admin_role = excluded.admin_role,
			scrobbling_enabled = excluded.scrobbling_enabled,
			settings_role = excluded.settings_role,
			stream_role = excluded.stream_role,
			jukebox_role = excluded.jukebox_role,
			download_role = excluded.download_role,
			upload_role = excluded.upload_role,
			playlist_role = excluded.playlist_role,
			cover_art_role = excluded.cover_art_role,
			comment_role = excluded.comment_role,
			podcast_role = excluded.podcast_role,
			share_role = excluded.share_role,
			video_conversion_role = excluded.video_conversion_role,
			max_bit_rate = excluded.max_bit_rate`

	_, err := DB.ExecContext(ctx, query, user.Username, user.Password, user.Email, user.ScrobblingEnabled,
		user.LdapAuthenticated, user.AdminRole, user.SettingsRole, user.StreamRole, user.JukeboxRole,
		user.DownloadRole, user.UploadRole, user.PlaylistRole, user.CoverArtRole,
		user.CommentRole, user.PodcastRole, user.ShareRole, user.VideoConversionRole, user.MaxBitRate)
	if err != nil {
		return 0, fmt.Errorf("upserting user in users table: %v", err)
	}

	upsertedUser, err := GetUserWithoutFoldersByUsername(ctx, user.Username)
	if err != nil {
		return 0, fmt.Errorf("getting upsertedUser: %v", err)
	}

	query = `DELETE FROM user_music_folders WHERE user_id = ?`
	_, err = DB.ExecContext(ctx, query, upsertedUser.Id)
	if err != nil {
		logger.Printf("Error deleting user music folders for user %d: %v", upsertedUser.Id, err)
		return 0, fmt.Errorf("deleting user music folders: %v", err)
	}

	for _, folderId := range user.Folders {
		query = `INSERT INTO user_music_folders (user_id, folder_id) VALUES (?, ?) ON CONFLICT(user_id, folder_id) DO NOTHING`
		_, err = DB.ExecContext(ctx, query, upsertedUser.Id, folderId)
		if err != nil {
			logger.Printf("Error inserting user music folder %d for user %d: %v", folderId, upsertedUser.Id, err)
			return 0, fmt.Errorf("inserting user music folder %d: %v", folderId, err)
		}
	}

	adminUser, err := GetUserByContext(ctx)
	if err != nil {
		logger.Printf("user %s upserted, failed to get name of the creating admin user: %v", user.Username, err)
	} else {
		logger.Printf("user %s upserted by admin user %s", user.Username, adminUser.Username)
	}

	return upsertedUser.Id, nil
}

func DeleteUserById(ctx context.Context, id int) error {
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

func anyAdminUsersExist(ctx context.Context) (bool, error) {
	query := `SELECT 1 FROM users where admin_role = true LIMIT 1`
	var exists int
	err := DB.QueryRowContext(ctx, query).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func GetEncryptedPasswordFromDB(ctx context.Context, username string) (string, int, error) {
	query := `SELECT password, id FROM users WHERE username = ?`
	var encryptedPassword string
	var userId int

	err := DB.QueryRowContext(ctx, query, username).Scan(&encryptedPassword, &userId)
	if err == sql.ErrNoRows {
		return "", 0, fmt.Errorf("user not found")
	} else if err != nil {
		return "", 0, fmt.Errorf("selecting encrypted password for user %s: %v", username, err)
	}

	return encryptedPassword, userId, nil
}
