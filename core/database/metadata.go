package database

import (
	"context"
	"fmt"
	"strings"
	"zene/core/logger"
	"zene/core/types"
)

func migrateMetadata(ctx context.Context) {
	schema := `CREATE TABLE metadata (
		file_path TEXT PRIMARY KEY,
		file_name TEXT NOT NULL,
		date_added TEXT NOT NULL,
		date_modified TEXT NOT NULL,
		format TEXT,
		duration TEXT,
		size TEXT,
		bitrate TEXT,
		title TEXT,
		artist TEXT,
		album TEXT,
		album_artist TEXT,
		genre TEXT,
		track_number INTEGER DEFAULT 0,
		total_tracks INTEGER DEFAULT 0,
		disc_number INTEGER DEFAULT 0,
		total_discs INTEGER DEFAULT 0,
		release_date TEXT,
		musicbrainz_artist_id TEXT NOT NULL,
		musicbrainz_album_id TEXT NOT NULL,
		musicbrainz_track_id TEXT NOT NULL,
		label TEXT,
		music_folder_id INTEGER DEFAULT 1,
		codec TEXT,
		bit_depth INTEGER,
		sample_rate INTEGER,
		channels INTEGER,
		FOREIGN KEY (music_folder_id) REFERENCES music_folders(id) ON DELETE CASCADE
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_metadata_track_id", "metadata", []string{"musicbrainz_track_id"}, false)
	createIndex(ctx, "idx_metadata_album_id", "metadata", []string{"musicbrainz_album_id"}, false)
	createIndex(ctx, "idx_metadata_artist_id", "metadata", []string{"musicbrainz_artist_id"}, false)
	createIndex(ctx, "idx_metadata_file_path_album_track ", "metadata", []string{"file_path", "musicbrainz_album_id", "musicbrainz_track_id"}, false)
	createIndex(ctx, "idx_metadata_artist", "metadata", []string{"artist"}, false)
	createIndex(ctx, "idx_metadata_album_artist", "metadata", []string{"album_artist"}, false)
	createIndex(ctx, "idx_metadata_artist_lower", "metadata", []string{"lower(artist)"}, false)
	createIndex(ctx, "idx_metadata_album_lower", "metadata", []string{"lower(album)"}, false)
	createIndex(ctx, "idx_metadata_title_lower", "metadata", []string{"lower(title)"}, false)
	createIndex(ctx, "idx_metadata_date_added_album_id", "metadata", []string{"date_added", "musicbrainz_album_id"}, false)
	createIndex(ctx, "idx_metadata_album_artist", "metadata", []string{"album_artist", "musicbrainz_album_id"}, false)
}

func UpsertMetadataRows(ctx context.Context, metadataSlice []types.Metadata) error {
	if len(metadataSlice) == 0 {
		return nil
	}

	const batchSize = 30
	numberOfColumns := 26 // TODO: derive this from the metadata struct rather than hardcoding it

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	for start := 0; start < len(metadataSlice); start += batchSize {
		end := start + batchSize
		if end > len(metadataSlice) {
			end = len(metadataSlice)
		}

		batch := metadataSlice[start:end]

		placeholders := make([]string, len(batch))
		args := make([]interface{}, 0, len(batch)*numberOfColumns)

		for i, m := range batch {
			placeholders[i] = "(" + strings.Repeat("?, ", numberOfColumns-1) + "?" + ")"
			args = append(args,
				m.FilePath, m.DateAdded, m.DateModified, m.FileName,
				m.Format, m.Duration, m.Size, m.Bitrate,
				m.Title, m.Artist, m.Album, m.AlbumArtist,
				m.Genre, m.TrackNumber, m.TotalTracks, m.DiscNumber,
				m.TotalDiscs, m.ReleaseDate, m.MusicBrainzArtistID,
				m.MusicBrainzAlbumID, m.MusicBrainzTrackID, m.Label,
				m.Codec, m.BitDepth, m.SampleRate, m.Channels,
			)
		}

		query := fmt.Sprintf(`
			INSERT INTO metadata (
				file_path, date_added, date_modified, file_name, format, duration, size, bitrate, title, artist, album,
				album_artist, genre, track_number, total_tracks, disc_number, total_discs, release_date,
				musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id, label, codec, bit_depth, sample_rate, channels
			) VALUES %s
			ON CONFLICT(file_path) DO UPDATE SET
				date_modified = excluded.date_modified,
				file_name = excluded.file_name,
				format = excluded.format,
				duration = excluded.duration,
				size = excluded.size,
				bitrate = excluded.bitrate,
				title = excluded.title,
				artist = excluded.artist,
				album = excluded.album,
				album_artist = excluded.album_artist,
				genre = excluded.genre,
				track_number = excluded.track_number,
				total_tracks = excluded.total_tracks,
				disc_number = excluded.disc_number,
				total_discs = excluded.total_discs,
				release_date = excluded.release_date,
				musicbrainz_artist_id = excluded.musicbrainz_artist_id,
				musicbrainz_album_id = excluded.musicbrainz_album_id,
				musicbrainz_track_id = excluded.musicbrainz_track_id,
				label = excluded.label,
				codec = excluded.codec,
				bit_depth = excluded.bit_depth,
				sample_rate = excluded.sample_rate,
				channels = excluded.channels
		`, strings.Join(placeholders, ", "))

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert batch %d-%d: %w", start, end, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func DeleteMetadataRows(ctx context.Context, filepaths []string) error {
	if len(filepaths) == 0 {
		return nil
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	const batchSize = 100
	for start := 0; start < len(filepaths); start += batchSize {
		end := start + batchSize
		if end > len(filepaths) {
			end = len(filepaths)
		}

		batch := filepaths[start:end]
		placeholders := make([]string, len(batch))
		args := make([]interface{}, len(batch))
		for i, fp := range batch {
			placeholders[i] = "?"
			args[i] = fp
		}

		query := fmt.Sprintf(
			`DELETE FROM metadata WHERE file_path IN (%s)`,
			strings.Join(placeholders, ","),
		)

		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			tx.Rollback()
			return fmt.Errorf("deleting metadata rows (batch %d-%d): %w", start, end, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	logger.Printf("Deleted %d metadata rows in batches of %d", len(filepaths), batchSize)
	return nil
}

type isValidMetadataResponse struct {
	MusicbrainzTrackId  bool `json:"track_valid"`
	MusicbrainzAlbumId  bool `json:"album_valid"`
	MusicbrainzArtistId bool `json:"artist_valid"`
}

func IsValidMetadataId(ctx context.Context, metadataId string) (bool, isValidMetadataResponse, error) {
	query := `SELECT musicbrainz_track_id, musicbrainz_album_id, musicbrainz_artist_id
		FROM metadata
		WHERE musicbrainz_track_id = ?
		OR musicbrainz_album_id = ?
		OR musicbrainz_artist_id = ?
		limit 1`
	row := DB.QueryRowContext(ctx, query, metadataId, metadataId, metadataId)

	var response isValidMetadataResponse

	var musicbrainzTrackId string
	var musicbrainzAlbumId string
	var musicbrainzArtistId string
	if err := row.Scan(&musicbrainzTrackId, &musicbrainzAlbumId, &musicbrainzArtistId); err != nil {
		return false, isValidMetadataResponse{}, fmt.Errorf("checking metadata ID validity: %v", err)
	}
	response.MusicbrainzAlbumId = musicbrainzAlbumId == metadataId
	response.MusicbrainzArtistId = musicbrainzArtistId == metadataId
	response.MusicbrainzTrackId = musicbrainzTrackId == metadataId

	isValid := response.MusicbrainzTrackId || response.MusicbrainzAlbumId || response.MusicbrainzArtistId

	return isValid, response, nil
}

type GetFileAndFolderCountsResponse struct {
	FileCount   int `json:"file_count"`
	FolderCount int `json:"folder_count"`
}

func GetFileAndFolderCounts(ctx context.Context) (GetFileAndFolderCountsResponse, error) {
	var fileCount, folderCount int

	query := `select count(distinct parent_directory) from (
		SELECT
			CASE
				WHEN instr(file_path, '\') > 0 THEN
					replace(file_path, replace(file_path, rtrim(file_path, replace(file_path, '\', char(1))), ''), '')
				ELSE
					replace(file_path, replace(file_path, rtrim(file_path, replace(file_path, '/', char(1))), ''), '')
			END AS parent_directory
		FROM metadata
	);`

	row := DB.QueryRowContext(ctx, query)
	if err := row.Scan(&folderCount); err != nil {
		return GetFileAndFolderCountsResponse{}, fmt.Errorf("getting folder count: %v", err)
	}

	query = "select count(*) from metadata"
	row = DB.QueryRowContext(ctx, query)
	if err := row.Scan(&fileCount); err != nil {
		return GetFileAndFolderCountsResponse{}, fmt.Errorf("getting file count: %v", err)
	}

	return GetFileAndFolderCountsResponse{
		FileCount:   fileCount,
		FolderCount: folderCount,
	}, nil
}
