package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"zene/core/types"
)

func createAlbumArtTable(ctx context.Context) {
	tableName := "album_art"
	schema := `CREATE TABLE IF NOT EXISTS album_art (
		musicbrainz_album_id TEXT PRIMARY KEY,
		date_modified TEXT NOT NULL
	);`
	createTable(ctx, tableName, schema)
}

func createArtistArtTable(ctx context.Context) {
	tableName := "artist_art"
	schema := `CREATE TABLE IF NOT EXISTS artist_art (
		musicbrainz_artist_id TEXT PRIMARY KEY,
		date_modified TEXT NOT NULL
	);`
	createTable(ctx, tableName, schema)
}

func SelectAlbumArtByMusicBrainzAlbumId(ctx context.Context, musicbrainzAlbumId string) (types.AlbumArtRow, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in SelectAlbumArtByMusicBrainzAlbumId: %v", err)
		return types.AlbumArtRow{}, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT musicbrainz_album_id, date_modified FROM album_art WHERE musicbrainz_album_id = $musicbrainz_album_id;`)
	defer stmt.Finalize()
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

func InsertAlbumArtRow(ctx context.Context, musicbrainzAlbumId string, dateModified string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in InsertAlbumArtRow: %v", err)
		return err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO album_art (musicbrainz_album_id, date_modified)
		VALUES ($musicbrainz_album_id, $date_modified)
		ON CONFLICT(musicbrainz_album_id) DO UPDATE SET date_modified=excluded.date_modified
	 	WHERE excluded.date_modified>album_art.date_modified;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_album_id", musicbrainzAlbumId)
	stmt.SetText("$date_modified", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert album art row: %v", err)
	}
	return nil
}

func SelectArtistSubDirectories(ctx context.Context, musicbrainzArtistId string) ([]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in SelectArtistSubDirectories: %v", err)
		return []string{}, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT DISTINCT f.dir_path FROM track_metadata m JOIN files f ON f.id = m.file_id WHERE m.musicbrainz_artist_id = $musicbrainz_artist_id;`)
	defer stmt.Finalize()
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

func SelectArtistArtByMusicBrainzArtistId(ctx context.Context, musicbrainzArtistId string) (types.ArtistArtRow, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in SelectArtistArtByMusicBrainzArtistId: %v", err)
		return types.ArtistArtRow{}, err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`SELECT musicbrainz_artist_id, date_modified FROM artist_art WHERE musicbrainz_artist_id = $musicbrainz_artist_id;`)
	defer stmt.Finalize()
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

func InsertArtistArtRow(ctx context.Context, musicbrainzArtistId string, dateModified string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := DbPool.Take(ctx)
	if err != nil {
		log.Printf("failed to take a db conn from the pool in InsertArtistArtRow: %v", err)
		return err
	}
	defer DbPool.Put(conn)

	stmt := conn.Prep(`INSERT INTO artist_art (musicbrainz_artist_id, date_modified)
	VALUES ($musicbrainz_artist_id, $date_modified)
	ON CONFLICT(musicbrainz_artist_id) DO UPDATE SET date_modified=excluded.date_modified
	 WHERE excluded.date_modified>artist_art.date_modified;`)
	defer stmt.Finalize()
	stmt.SetText("$musicbrainz_artist_id", musicbrainzArtistId)
	stmt.SetText("$date_modified", time.Now().Format(time.RFC3339Nano))

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert artist art row: %v", err)
	}
	return nil
}
