package database

import (
	"log"

	"zombiezen.com/go/sqlite"
)

var err error

var stmtDoesTableExist *sqlite.Stmt
var stmtCreateTriggerIfNotExists *sqlite.Stmt

var stmtInsertIntoFiles *sqlite.Stmt
var stmtDeleteFileById *sqlite.Stmt
var stmtSelectFileByFilePath *sqlite.Stmt
var stmtSelectFileByFileId *sqlite.Stmt
var stmtSelectAllFiles *sqlite.Stmt

var stmtInsertTrackMetadataRow *sqlite.Stmt
var stmtDeleteMetadataByFileId *sqlite.Stmt
var stmtSelectMetadataRecentlyAddedWithLimit *sqlite.Stmt
var stmtSelectMetadataRecentlyAdded *sqlite.Stmt
var stmtSelectMetadataRandomisedWithLimit *sqlite.Stmt
var stmtSelectMetadataRandomised *sqlite.Stmt
var stmtSelectMetadataWithLimit *sqlite.Stmt
var stmtSelectMetadata *sqlite.Stmt
var stmtSelectMetadataByAlbumID *sqlite.Stmt
var stmtSelectDistinctGenres *sqlite.Stmt

var stmtSelectAllArtists *sqlite.Stmt
var stmtSelectAllAlbumArtists *sqlite.Stmt
var stmtSelectArtist *sqlite.Stmt

var stmtSelectAlbumArtByMusicBrainzAlbumId *sqlite.Stmt
var stmtSelectArtistArtByMusicBrainzArtistId *sqlite.Stmt
var stmtSelectArtistSubDirectories *sqlite.Stmt
var stmtInsertAlbumArtRow *sqlite.Stmt
var stmtInsertArtistArtRow *sqlite.Stmt
var stmtSelectAllAlbums *sqlite.Stmt
var stmtSelectRandomizedAlbumsWithLimit *sqlite.Stmt
var stmtSelectRandomizedAlbums *sqlite.Stmt
var stmtSelectAlbumsWithLimit *sqlite.Stmt
var stmtSelectAlbumsRecentlyAdded *sqlite.Stmt
var stmtSelectAlbumsRecentlyAddedWithLimit *sqlite.Stmt

var stmtSelectLastScan *sqlite.Stmt
var stmtInsertScanRow *sqlite.Stmt

var stmtSelectFullTextSearchFromMetadata *sqlite.Stmt
var stmtSelectFtsGenre *sqlite.Stmt
var stmtSelectFtsArtist *sqlite.Stmt

func prepareInitStatements() {
	log.Println("Preparing SQL init statements")

	stmtDoesTableExist, err = DbRO.Prepare(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = $table_name;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtDoesTableExist: %v", err)
	}

	stmtCreateTriggerIfNotExists, err = DbRO.Prepare("SELECT name FROM sqlite_master WHERE type='trigger' AND name=$triggername")
	if err != nil {
		log.Fatalf("Failed to prepare stmtCreateTriggerIfNotExists: %v", err)
	}
}

func prepareStatements() {
	log.Println("Preparing SQL statements")

	stmtSelectAlbumArtByMusicBrainzAlbumId, err = DbRO.Prepare(`SELECT musicbrainz_album_id, date_modified FROM album_art WHERE musicbrainz_album_id = $musicbrainz_album_id;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectAlbumArtByMusicBrainzAlbumId: %v", err)
	}

	stmtSelectArtistArtByMusicBrainzArtistId, err = DbRO.Prepare(`SELECT musicbrainz_artist_id, date_modified FROM artist_art WHERE musicbrainz_artist_id = $musicbrainz_artist_id;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectArtistArtByMusicBrainzArtistId: %v", err)
	}

	stmtSelectArtistSubDirectories, err = DbRO.Prepare(`SELECT DISTINCT f.dir_path FROM track_metadata m JOIN files f ON f.id = m.file_id WHERE m.musicbrainz_artist_id = $musicbrainz_artist_id;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectArtistSubDirectories: %v", err)
	}

	stmtInsertAlbumArtRow, err = DbRW.Prepare(`INSERT INTO album_art (musicbrainz_album_id, date_modified)
		VALUES ($musicbrainz_album_id, $date_modified)
		ON CONFLICT(musicbrainz_album_id) DO UPDATE SET date_modified=excluded.date_modified
	 	WHERE excluded.date_modified>album_art.date_modified;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtInsertAlbumArtRow: %v", err)
	}

	stmtInsertArtistArtRow, err = DbRW.Prepare(`INSERT INTO artist_art (musicbrainz_artist_id, date_modified)
	VALUES ($musicbrainz_artist_id, $date_modified)
	ON CONFLICT(musicbrainz_artist_id) DO UPDATE SET date_modified=excluded.date_modified
	 WHERE excluded.date_modified>artist_art.date_modified;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtInsertArtistArtRow: %v", err)
	}

	stmtInsertIntoFiles, err = DbRW.Prepare(`INSERT INTO files (dir_path, file_path, filename, date_added, date_modified)
		VALUES ($dir_path, $file_path, $filename, $date_added, $date_modified)
		ON CONFLICT(file_path) DO UPDATE SET date_modified=excluded.date_modified
	 	WHERE excluded.date_modified>files.date_modified;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtInsertIntoFiles: %v", err)
	}

	stmtDeleteFileById, err = DbRW.Prepare(`delete FROM files WHERE id = $id;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtDeleteFileById: %v", err)
	}

	stmtSelectFileByFilePath, err = DbRO.Prepare(`SELECT id, dir_path, file_path, filename, date_added, date_modified FROM files WHERE file_path = $file_path;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectFileByFilePath: %v", err)
	}

	stmtSelectFileByFileId, err = DbRO.Prepare(`SELECT id, dir_path, filename, file_path, date_added, date_modified FROM files WHERE id = $fileid;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectFileByFileId: %v", err)
	}

	stmtSelectAllFiles, err = DbRO.Prepare(`SELECT id, dir_path, filename, file_path, date_added, date_modified FROM files;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectAllFiles: %v", err)
	}

	stmtInsertTrackMetadataRow, err = DbRW.Prepare(`INSERT INTO track_metadata (
		file_id, filename, format, duration, size, bitrate, title, artist, album,
		album_artist, genre, track_number, total_tracks, disc_number, total_discs, release_date,
		musicbrainz_artist_id, musicbrainz_album_id, musicbrainz_track_id, label
	) VALUES (
	  $file_id, $filename, $format, $duration, $size, $bitrate, $title, $artist, $album,
		$album_artist, $genre, $track_number, $total_tracks, $disc_number, $total_discs, $release_date,
		$musicbrainz_artist_id, $musicbrainz_album_id, $musicbrainz_track_id, $label
	 )`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtInsertTrackMetadataRow: %v", err)
	}

	stmtDeleteMetadataByFileId, err = DbRW.Prepare(`delete FROM track_metadata WHERE file_id = $file_id;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtDeleteMetadataByFileId: %v", err)
	}

	stmtSelectAllArtists, err = DbRO.Prepare(`SELECT DISTINCT artist, musicbrainz_artist_id FROM track_metadata	ORDER BY artist;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectAllArtists: %v", err)
	}

	stmtSelectAllAlbumArtists, err = DbRO.Prepare(`SELECT DISTINCT album_artist, musicbrainz_artist_id FROM track_metadata where artist = album_artist ORDER BY artist;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectAllAlbumArtists: %v", err)
	}

	stmtSelectArtist, err = DbRO.Prepare(`SELECT DISTINCT artist, musicbrainz_artist_id FROM track_metadata	where musicbrainz_artist_id = $musicbrainz_artist_id;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectArtist: %v", err)
	}

	stmtSelectAllAlbums, err = DbRO.Prepare(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM track_metadata group by album ORDER BY album;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectAllAlbums: %v", err)
	}

	stmtSelectRandomizedAlbumsWithLimit, err = DbRO.Prepare(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM track_metadata group by album ORDER BY random() limit $limit;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectRandomizedAlbumsWithLimit: %v", err)
	}

	stmtSelectAlbumsWithLimit, err = DbRO.Prepare(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM track_metadata group by album ORDER BY album limit $limit;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectAlbumsWithLimit: %v", err)
	}

	stmtSelectRandomizedAlbums, err = DbRO.Prepare(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date FROM track_metadata group by album ORDER BY random();`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectRandomizedAlbums: %v", err)
	}

	stmtSelectAlbumsRecentlyAdded, err = DbRO.Prepare(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM track_metadata m join files f on m.file_id = f.id group by album ORDER BY f.date_added desc;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectAlbumsRecentlyAdded: %v", err)
	}

	stmtSelectAlbumsRecentlyAddedWithLimit, err = DbRO.Prepare(`SELECT DISTINCT album, musicbrainz_album_id, album_artist, musicbrainz_artist_id, genre, release_date, date_added FROM track_metadata m join files f on m.file_id = f.id group by album ORDER BY f.date_added desc limit $limit;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectAlbumsRecentlyAddedWithLimit: %v", err)
	}

	stmtSelectMetadataWithLimit, err = DbRO.Prepare(`SELECT * FROM track_metadata ORDER BY id limit $limit;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectMetadataWithLimit: %v", err)
	}

	stmtSelectMetadata, err = DbRO.Prepare(`SELECT * FROM track_metadata ORDER BY id;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectMetadata: %v", err)
	}

	stmtSelectMetadataRandomisedWithLimit, err = DbRO.Prepare(`SELECT * FROM track_metadata order by random() limit $limit;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectMetadataRandomisedWithLimit: %v", err)
	}

	stmtSelectMetadataRandomised, err = DbRO.Prepare(`SELECT * FROM track_metadata order by random();`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectMetadataRandomised: %v", err)
	}

	stmtSelectMetadataRecentlyAddedWithLimit, err = DbRO.Prepare(`SELECT * FROM track_metadata m join files f on m.file_id = f.id ORDER BY f.date_added desc limit $limit;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectMetadataRecentlyAddedWithLimit: %v", err)
	}

	stmtSelectMetadataRecentlyAdded, err = DbRO.Prepare(`SELECT * FROM track_metadata m join files f on m.file_id = f.id ORDER BY f.date_added desc;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectMetadataRecentlyAdded: %v", err)
	}

	stmtSelectMetadataByAlbumID, err = DbRO.Prepare(`SELECT * FROM track_metadata where musicbrainz_album_id = $musicbrainz_album_id ORDER BY id;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectMetadataByAlbumID: %v", err)
	}

	stmtSelectLastScan, err = DbRO.Prepare(`SELECT id, scan_date, file_count, date_modified from scans order by id desc limit 1;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectLastScan: %v", err)
	}

	stmtInsertScanRow, err = DbRW.Prepare(`INSERT INTO scans (scan_date, file_count, date_modified)
		VALUES ($scan_date, $file_count, $date_modified);`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtInsertScanRow: %v", err)
	}

	stmtSelectDistinctGenres, err = DbRO.Prepare(`SELECT DISTINCT genre FROM track_metadata;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectDistinctGenres: %v", err)
	}

	stmtSelectFullTextSearchFromMetadata, err = DbRO.Prepare(`select distinct m.file_id, m.filename, m.format, m.duration, m.size, m.bitrate, m.title, m.artist, m.album,
		m.album_artist, m.genre, m.track_number, m.total_tracks, m.disc_number, m.total_discs, m.release_date,
		m.musicbrainz_artist_id, m.musicbrainz_album_id, m.musicbrainz_track_id, m.label
		FROM track_metadata m JOIN track_metadata_fts f ON m.file_id = f.file_id
		WHERE track_metadata_fts MATCH $searchQuery
		ORDER BY m.file_id DESC;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectFullTextSearchFromMetadata: %v", err)
	}

	stmtSelectFtsArtist, err = DbRO.Prepare(`select distinct m.artist, m.musicbrainz_artist_id
		FROM track_metadata m JOIN artists_fts f ON m.file_id = f.file_id
		WHERE artists_fts MATCH $searchQuery
		ORDER BY m.file_id DESC;`)
	if err != nil {
		log.Fatalf("Failed to prepare stmtSelectFtsArtist: %v", err)
	}
}
