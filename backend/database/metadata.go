package database

import (
	"fmt"
	"log"
	"zene/types"
)

func createMetadataTable() {
	tableName := "metadata"
	schema := `CREATE TABLE IF NOT EXISTS track_metadata (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_id INTEGER,
		filename TEXT,
		format TEXT,
		duration TEXT,
		size TEXT,
		bitrate TEXT,
		title TEXT,
		artist TEXT,
		album TEXT,
		album_artist TEXT,
		genre TEXT,
		track_number TEXT,
		total_tracks TEXT,
		disc_number TEXT,
		total_discs TEXT,
		release_date TEXT,
		musicbrainz_artist_id TEXT,
		musicbrainz_album_id TEXT,
		musicbrainz_track_id TEXT,
		label TEXT
	);`
	createTable(tableName, schema)
}

func createMetadataTriggers() {
	createTriggerIfNotExists("track_metadata_after_delete_album_art", `CREATE TRIGGER track_metadata_after_delete_album_art AFTER DELETE ON track_metadata
	BEGIN
			DELETE FROM album_art WHERE musicbrainz_album_id = old.musicbrainz_album_id;
	END;`)
}

func InsertTrackMetadataRow(fileRowId int, metadata types.TrackMetadata) error {
	stmt, err := Db.Prepare(`INSERT INTO track_metadata (
		file_id, filename, format, duration, size, bitrate, title, artist, album,
		album_artist, genre, track_number, total_tracks, disc_number, total_discs, release_date,
		musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id, label
	) VALUES (
	  $file_id, $filename, $format, $duration, $size, $bitrate, $title, $artist, $album,
		$album_artist, $genre, $track_number, $total_tracks, $disc_number, $total_discs, $release_date,
		$musicbrainz_artist_id, $musicbrainz_album_id, $musicbrainz_track_id, $label
	 )`)
	if err != nil {
		return err
	}
	defer stmt.Finalize()

	stmt.SetInt64("$file_id", int64(fileRowId))
	stmt.SetText("$filename", metadata.Filename)
	stmt.SetText("$format", metadata.Format)
	stmt.SetText("$duration", metadata.Duration)
	stmt.SetText("$size", metadata.Size)
	stmt.SetText("$bitrate", metadata.Bitrate)
	stmt.SetText("$title", metadata.Title)
	stmt.SetText("$artist", metadata.Artist)
	stmt.SetText("$album", metadata.Album)
	stmt.SetText("$album_artist", metadata.AlbumArtist)
	stmt.SetText("$genre", metadata.Genre)
	stmt.SetText("$track_number", metadata.TrackNumber)
	stmt.SetText("$total_tracks", metadata.TotalTracks)
	stmt.SetText("$disc_number", metadata.DiscNumber)
	stmt.SetText("$total_discs", metadata.TotalDiscs)
	stmt.SetText("$release_date", metadata.ReleaseDate)
	stmt.SetText("$musicbrainz_artist_id", metadata.MusicBrainzArtistID)
	stmt.SetText("$musicbrainz_album_id", metadata.MusicBrainzAlbumID)
	stmt.SetText("$musicbrainz_track_id", metadata.MusicBrainzTrackID)
	stmt.SetText("$label", metadata.Label)

	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to insert metadata row: %v", err)
	}

	return nil
}

func DeleteMetadataByFileId(file_id int) error {
	stmt, err := Db.Prepare(`delete FROM metadata WHERE file_id = $file_id;`)
	if err != nil {
		return err
	}
	defer stmt.Finalize()
	stmt.SetInt64("$file_id", int64(file_id))
	_, err = stmt.Step()
	if err != nil {
		return fmt.Errorf("failed to delete metadata row for file_id %d: %v", file_id, err)
	}
	log.Printf("Deleted metadata row for file_id %d", file_id)
	return nil
}

func SelectAllArtists() ([]types.ArtistResponse, error) {
	stmt, err := Db.Prepare(`SELECT DISTINCT artist, musicbrainz_artist_id
		FROM track_metadata
		ORDER BY artist;`)

	var rows []types.ArtistResponse

	if err != nil {
		log.Printf("Error selecting artists from track_metadata: %v", err)
		return rows, err
	}
	defer stmt.Finalize()

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return rows, err
		} else if !hasRow {
			break
		} else {
			row := types.ArtistResponse{
				Artist:              stmt.GetText("artist"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
			}
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func SelectAllAlbums() ([]types.AlbumsResponse, error) {
	stmt, err := Db.Prepare(`SELECT DISTINCT album, musicbrainz_album_id, artist, musicbrainz_artist_id
		FROM track_metadata
		ORDER BY album;`)

	var rows []types.AlbumsResponse

	if err != nil {
		log.Printf("Error selecting albums from track_metadata: %v", err)
		return rows, err
	}
	defer stmt.Finalize()

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return rows, err
		} else if !hasRow {
			break
		} else {
			row := types.AlbumsResponse{
				Album:               stmt.GetText("album"),
				Artist:              stmt.GetText("artist"),
				MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
			}
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func SelectAllMetadata() ([]types.TrackMetadata, error) {
	stmt, err := Db.Prepare(`SELECT * FROM track_metadata ORDER BY id;`)

	var rows []types.TrackMetadata

	if err != nil {
		log.Printf("Error selecting albums from track_metadata: %v", err)
		return rows, err
	}
	defer stmt.Finalize()

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return rows, err
		} else if !hasRow {
			break
		} else {

			row := types.TrackMetadata{
				Id:                  int(stmt.GetInt64("id")),
				FileId:              int(stmt.GetInt64("file_id")),
				Filename:            stmt.GetText("filename"),
				Format:              stmt.GetText("format"),
				Duration:            stmt.GetText("duration"),
				Size:                stmt.GetText("size"),
				Bitrate:             stmt.GetText("bitrate"),
				Title:               stmt.GetText("title"),
				Artist:              stmt.GetText("artist"),
				Album:               stmt.GetText("album"),
				AlbumArtist:         stmt.GetText("album_artist"),
				Genre:               stmt.GetText("genre"),
				TrackNumber:         stmt.GetText("track_number"),
				TotalTracks:         stmt.GetText("total_tracks"),
				DiscNumber:          stmt.GetText("disc_number"),
				TotalDiscs:          stmt.GetText("total_discs"),
				ReleaseDate:         stmt.GetText("release_date"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
				MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
				MusicBrainzTrackID:  stmt.GetText("musicbrainz_track_id"),
				Label:               stmt.GetText("label"),
			}
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func SelectMetadataByAlbumID(musicbrainz_album_id string) ([]types.TrackMetadata, error) {
	stmt, err := Db.Prepare(`SELECT * FROM track_metadata where musicbrainz_album_id = $musicbrainz_album_id ORDER BY id;`)

	var rows []types.TrackMetadata

	if err != nil {
		log.Printf("Error selecting albums from track_metadata: %v", err)
		return rows, err
	}
	defer stmt.Finalize()

	stmt.SetText("$musicbrainz_album_id", musicbrainz_album_id)

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return rows, err
		} else if !hasRow {
			break
		} else {

			row := types.TrackMetadata{
				Id:                  int(stmt.GetInt64("id")),
				FileId:              int(stmt.GetInt64("file_id")),
				Filename:            stmt.GetText("filename"),
				Format:              stmt.GetText("format"),
				Duration:            stmt.GetText("duration"),
				Size:                stmt.GetText("size"),
				Bitrate:             stmt.GetText("bitrate"),
				Title:               stmt.GetText("title"),
				Artist:              stmt.GetText("artist"),
				Album:               stmt.GetText("album"),
				AlbumArtist:         stmt.GetText("album_artist"),
				Genre:               stmt.GetText("genre"),
				TrackNumber:         stmt.GetText("track_number"),
				TotalTracks:         stmt.GetText("total_tracks"),
				DiscNumber:          stmt.GetText("disc_number"),
				TotalDiscs:          stmt.GetText("total_discs"),
				ReleaseDate:         stmt.GetText("release_date"),
				MusicBrainzArtistID: stmt.GetText("musicbrainz_artist_id"),
				MusicBrainzAlbumID:  stmt.GetText("musicbrainz_album_id"),
				MusicBrainzTrackID:  stmt.GetText("musicbrainz_track_id"),
				Label:               stmt.GetText("label"),
			}
			rows = append(rows, row)
		}
	}
	return rows, nil
}
