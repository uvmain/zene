package database

import (
	"log"

	"zombiezen.com/go/sqlite"
)

var err error
var stmtSelectAlbumArtByMusicBrainzAlbumId *sqlite.Stmt
var stmtInsertAlbumArtRow *sqlite.Stmt
var stmtInsertIntoFiles *sqlite.Stmt
var stmtDeleteFileById *sqlite.Stmt
var stmtSelectFileByFilePath *sqlite.Stmt
var stmtSelectFileByFilename *sqlite.Stmt
var stmtSelectAllFiles *sqlite.Stmt
var stmtInsertTrackMetadataRow *sqlite.Stmt
var stmtDeleteMetadataByFileId *sqlite.Stmt
var stmtSelectAllArtists *sqlite.Stmt
var stmtSelectAllAlbums *sqlite.Stmt
var stmtSelectAllMetadata *sqlite.Stmt
var stmtSelectMetadataByAlbumID *sqlite.Stmt
var stmtSelectLastScan *sqlite.Stmt
var stmtInsertScanRow *sqlite.Stmt

func prepareStatements() {
	log.Println("Preparing SQL statements")
	stmtSelectAlbumArtByMusicBrainzAlbumId, err = Db.Prepare(`SELECT musicbrainz_album_id, date_modified FROM album_art WHERE musicbrainz_album_id = $musicbrainz_album_id;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtInsertAlbumArtRow, err = Db.Prepare(`INSERT INTO album_art (musicbrainz_album_id, date_modified)
		VALUES ($musicbrainz_album_id, $date_modified)
		ON CONFLICT(musicbrainz_album_id) DO UPDATE SET date_modified=excluded.date_modified
	 	WHERE excluded.date_modified>album_art.date_modified;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtInsertIntoFiles, err = Db.Prepare(`INSERT INTO files (dir_path, filename, date_added, date_modified)
		VALUES ($dir_path, $filename, $date_added, $date_modified);`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtDeleteFileById, err = Db.Prepare(`delete FROM files WHERE id = $id;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtSelectFileByFilePath, err = Db.Prepare(`SELECT id, dir_path, filename, date_added, date_modified FROM files WHERE dir_path = $dir_path and filename = $filename;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtSelectFileByFilename, err = Db.Prepare(`SELECT id, dir_path, filename, date_added, date_modified FROM files WHERE filename = $filename;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtSelectAllFiles, err = Db.Prepare(`SELECT id, dir_path, filename, date_added, date_modified FROM files;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtInsertTrackMetadataRow, err = Db.Prepare(`INSERT INTO track_metadata (
		file_id, filename, format, duration, size, bitrate, title, artist, album,
		album_artist, genre, track_number, total_tracks, disc_number, total_discs, release_date,
		musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id, label
	) VALUES (
	  $file_id, $filename, $format, $duration, $size, $bitrate, $title, $artist, $album,
		$album_artist, $genre, $track_number, $total_tracks, $disc_number, $total_discs, $release_date,
		$musicbrainz_artist_id, $musicbrainz_album_id, $musicbrainz_track_id, $label
	 )`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtDeleteMetadataByFileId, err = Db.Prepare(`delete FROM track_metadata WHERE file_id = $file_id;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtSelectAllArtists, err = Db.Prepare(`SELECT DISTINCT artist, musicbrainz_artist_id FROM track_metadata	ORDER BY artist;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtSelectAllAlbums, err = Db.Prepare(`SELECT DISTINCT album, musicbrainz_album_id, artist, musicbrainz_artist_id FROM track_metadata ORDER BY album;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtSelectAllMetadata, err = Db.Prepare(`SELECT * FROM track_metadata ORDER BY id;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtSelectMetadataByAlbumID, err = Db.Prepare(`SELECT * FROM track_metadata where musicbrainz_album_id = $musicbrainz_album_id ORDER BY id;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtSelectLastScan, err = Db.Prepare(`SELECT id, scan_date, file_count, date_modified from scans order by id desc limit 1;`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}

	stmtInsertScanRow, err = Db.Prepare(`INSERT INTO scans (scan_date, file_count, date_modified)
		VALUES ($scan_date, $file_count, $date_modified);`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
}
