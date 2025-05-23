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

func createArtistArtTable() {
	tableName := "artist_art"
	schema := `CREATE TABLE IF NOT EXISTS artist_art (
		musicbrainz_artist_id TEXT PRIMARY KEY,
		date_modified TEXT NOT NULL
	);`
	createTable(tableName, schema)
}

func SelectAlbumArtByMusicBrainzAlbumId(musicbrainzAlbumId string) (types.AlbumArtRow, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

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
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtInsertAlbumArtRow
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$musicbrainz_album_id", musicbrainzAlbumId)
	stmt.SetText("$date_modified", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert album art row: %v", err)
	}
	return nil
}

func SelectArtistSubDirectories(musicbrainzArtistId string) ([]string, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectArtistSubDirectories
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$musicbrainz_artist_id", musicbrainzArtistId)

	var rows []string

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return []string{}, err
		} else if !hasRow {
			break
		}
		rows = append(rows, stmt.GetText("dir_path"))
	}

	return rows, nil
}

func SelectArtistArtByMusicBrainzArtistId(musicbrainzArtistId string) (types.ArtistArtRow, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtSelectArtistArtByMusicBrainzArtistId
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$musicbrainz_artist_id", musicbrainzArtistId)

	if hasRow, err := stmt.Step(); err != nil {
		return types.ArtistArtRow{}, err
	} else if !hasRow {
		return types.ArtistArtRow{}, nil
	} else {
		var row types.ArtistArtRow
		row.MusicbrainzArtistId = stmt.GetText("musicbrainz_artist_id")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}

func InsertArtistArtRow(musicbrainzArtistId string, dateModified string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	stmt := stmtInsertArtistArtRow
	stmt.Reset()
	stmt.ClearBindings()
	stmt.SetText("$musicbrainz_artist_id", musicbrainzArtistId)
	stmt.SetText("$date_modified", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert artist art row: %v", err)
	}
	return nil
}
