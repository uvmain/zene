package database

import (
	"fmt"
	"time"
	"zene/types"
)

func createAlbumArtTable() {
	tableName := "album_art"
	schema := `CREATE TABLE IF NOT EXISTS album_art (
		musicbrainz_album_id TEXT PRIMARY KEY,
		date_modified TEXT NOT NULL
	);`
	createTable(tableName, schema)
}

func SelectAlbumArtByMusicBrainzAlbumId(musicbrainzAlbumId string) (types.AlbumArtRow, error) {
	stmt := stmtSelectAlbumArtByMusicBrainzAlbumId
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$musicbrainz_album_id", musicbrainzAlbumId)

	if hasRow, err := stmt.Step(); err != nil {
		return types.AlbumArtRow{}, err
	} else if !hasRow {
		return types.AlbumArtRow{}, nil
	} else {
		var row types.AlbumArtRow
		row.MusicbrainzAlbumId = stmt.GetText("musicbrainz_album_id")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}

func InsertAlbumArtRow(musicbrainzAlbumId string, dateModified string) error {
	stmt := stmtInsertAlbumArtRow
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$musicbrainz_album_id", musicbrainzAlbumId)
	stmt.SetText("$date_modified", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert Scans row: %v", err)
	}
	return nil
}
