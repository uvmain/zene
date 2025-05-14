package database

import "zene/types"

func createAlbumArtTable() {
	tableName := "album_art"
	schema := `CREATE TABLE IF NOT EXISTS album_art (
		musicbrainz_album_id TEXT PRIMARY KEY,
		date_added TEXT NOT NULL,
		date_modified TEXT NOT NULL
	);`
	createTable(tableName, schema)
}

func SelectAlbumArtByMusicBrainzAlbumId(musicbrainz_album_id string) (types.AlbumArtRow, error) {
	stmt, err := Db.Prepare(`SELECT musicbrainz_album_id, date_added, date_added FROM album_art WHERE musicbrainz_album_id = $musicbrainz_album_id;`)
	if err != nil {
		return types.AlbumArtRow{}, err
	}
	defer stmt.Finalize()

	stmt.SetText("$musicbrainz_album_id", musicbrainz_album_id)

	if hasRow, err := stmt.Step(); err != nil {
		return types.AlbumArtRow{}, err
	} else if !hasRow {
		return types.AlbumArtRow{}, nil
	} else {
		var row types.AlbumArtRow
		row.MusicbrainzAlbumId = stmt.GetText("musicbrainz_album_id")
		row.DateAdded = stmt.GetText("date_added")
		row.DateModified = stmt.GetText("date_modified")
		return row, nil
	}
}
