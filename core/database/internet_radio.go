package database

import (
	"context"
	"fmt"
	"zene/core/types"
)

func migrateInternetRadio(ctx context.Context) {
	schema := `CREATE TABLE internet_radio (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		stream_url TEXT NOT NULL,
		homepage_url TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_internet_radio_user", "internet_radio", []string{"user_id"}, false)
}

func InsertInternetRadio(ctx context.Context, stationName string, streamUrl string, homepageUrl string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}

	query := `INSERT INTO internet_radio (user_id, name, stream_url, homepage_url)
		VALUES (?, ?, ?, ?)`
	_, err = DB.ExecContext(ctx, query, user.Id, stationName, streamUrl, homepageUrl)
	if err != nil {
		return fmt.Errorf("inserting internet radio: %v", err)
	}
	return nil
}

func GetInternetRadioStations(ctx context.Context) ([]types.InternetRadio, error) {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return []types.InternetRadio{}, fmt.Errorf("getting user from context: %v", err)
	}
	query := "SELECT ir.id, ir.name, ir.stream_url, ir.homepage_url FROM internet_radio ir join users u on ir.user_id = u.id WHERE u.id = ? order by ir.id asc;"
	rows, err := DB.QueryContext(ctx, query, user.Id)
	if err != nil {
		return []types.InternetRadio{}, fmt.Errorf("querying internet radio: %v", err)
	}
	defer rows.Close()

	var result []types.InternetRadio
	for rows.Next() {
		var row types.InternetRadio
		err := rows.Scan(&row.Id, &row.Name, &row.StreamUrl, &row.HomepageUrl)
		if err != nil {
			return []types.InternetRadio{}, fmt.Errorf("scanning internet radio row: %v", err)
		}
		result = append(result, row)
	}
	return result, nil
}

func UpdateInternetRadioStation(ctx context.Context, id string, name string, streamUrl string, homepageUrl string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}
	var args []interface{}
	query := `UPDATE internet_radio SET name = ?, stream_url = ?`
	args = append(args, name, streamUrl)

	if homepageUrl != "" {
		query += `, homepage_url = ?`
		args = append(args, homepageUrl)
	}

	query += ` WHERE id = ? AND user_id = ?`
	args = append(args, id, user.Id)

	_, err = DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("updating internet radio station: %v", err)
	}
	return nil
}

func DeleteInternetRadioStation(ctx context.Context, id string) error {
	user, err := GetUserByContext(ctx)
	if err != nil {
		return fmt.Errorf("getting user from context: %v", err)
	}
	query := `DELETE FROM internet_radio WHERE id = ? AND user_id = ?`
	_, err = DB.ExecContext(ctx, query, id, user.Id)
	if err != nil {
		return fmt.Errorf("deleting internet radio station: %v", err)
	}
	return nil
}
