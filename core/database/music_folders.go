package database

import (
	"context"
	"database/sql"
	"fmt"
	"zene/core/config"
	"zene/core/logger"
	"zene/core/types"
)

func migrateMusicFolders(ctx context.Context) {
	schema := `CREATE TABLE music_folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
	);`
	createTable(ctx, schema)

	for _, dir := range config.MusicDirs {
		if dir != "" {
			if len(dir) > 255 {
				dir = dir[:255] // limit directory name to 255 characters
			}
			err := InsertMusicFolder(ctx, dir)
			if err != nil {
				fmt.Printf("Error inserting music folder %s: %v\n", dir, err)
			}
		}
	}
}

func GetMusicFolders(ctx context.Context) ([]types.MusicFolder, error) {
	query := `SELECT id, name FROM music_folders`
	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying music folders: %v", err)
	}
	defer rows.Close()

	var folders []types.MusicFolder
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("scanning music folder name: %v", err)
		}
		folders = append(folders, types.MusicFolder{
			Id:   id,
			Name: name,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over music folders: %v", err)
	}

	return folders, nil
}

func GetMusicFolderById(ctx context.Context, id int) (types.MusicFolder, error) {
	query := `SELECT id, name FROM music_folders where id = ?`
	var row types.MusicFolder
	err := DB.QueryRowContext(ctx, query, id).Scan(&row.Id, &row.Name)
	if err == sql.ErrNoRows {
		return types.MusicFolder{}, fmt.Errorf("music folder with id %d not found", id)
	} else if err != nil {
		return types.MusicFolder{}, fmt.Errorf("querying music folder: %v", err)
	}
	return row, nil
}

func InsertMusicFolder(ctx context.Context, name string) error {
	query := `SELECT COUNT(*) FROM music_folders WHERE name = ?`
	var count int
	err := DB.QueryRowContext(ctx, query, name).Scan(&count)
	if err != nil {
		return fmt.Errorf("checking if music folder exists: %v", err)
	}
	if count > 0 {
		logger.Printf("Music folder '%s' already exists, skipping insert.", name)
		return nil
	}

	query = `INSERT INTO music_folders (name) VALUES (?)`
	_, err = DB.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("inserting music folder: %v", err)
	}
	logger.Printf("Inserted music folder: %s", name)
	return nil
}
