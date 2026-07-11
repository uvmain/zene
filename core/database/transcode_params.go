package database

import (
	"context"
	"fmt"
	"time"
	"zene/core/config"
	"zene/core/logic"
	"zene/core/types"
)

func createTranscodeParamsTable(ctx context.Context) {
	schema := `CREATE TABLE transcode_params (
		paramString TEXT PRIMARY KEY,
		media_id TEXT NOT NULL,
		media_type TEXT NOT NULL,
		container TEXT NOT NULL,
		audio_codec TEXT NOT NULL,
		protocol TEXT NOT NULL,
		target_format TEXT NOT NULL,
		bitrate INTEGER NOT NULL,
		audio_channels INTEGER NOT NULL,
		created_at TEXT NOT NULL
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_transcode_params_media", "transcode_params", []string{"media_id", "media_type"}, false)
}

func UpsertTranscodeParams(ctx context.Context, row types.TranscodeParamsRow) error {
	query := `INSERT INTO transcode_params (paramString, media_id, media_type, container, audio_codec, protocol, target_format, bitrate, audio_channels, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(paramString) DO UPDATE SET
			media_id=excluded.media_id,
			media_type=excluded.media_type,
			container=excluded.container,
			audio_codec=excluded.audio_codec,
			protocol=excluded.protocol,
			target_format=excluded.target_format,
			bitrate=excluded.bitrate,
			audio_channels=excluded.audio_channels,
			created_at=excluded.created_at`

	createdAt := logic.GetCurrentTimeFormatted()

	_, err := DB.ExecContext(ctx, query, row.ParamString, row.MediaId, row.MediaType, row.Container, row.AudioCodec, row.Protocol, row.TargetFormat, row.Bitrate, row.AudioChannels, createdAt)
	if err != nil {
		return err
	}
	return nil
}

func GetTranscodeParams(ctx context.Context, paramString string) (types.TranscodeParamsRow, error) {
	query := `SELECT media_id, media_type, container, audio_codec, protocol, target_format, bitrate, audio_channels, created_at FROM transcode_params WHERE paramString = ?`
	var mediaId, mediaType, container, audioCodec, protocol, targetFormat, createdAt string
	var bitrate, audioChannels int

	err := DB.QueryRowContext(ctx, query, paramString).Scan(&mediaId, &mediaType, &container, &audioCodec, &protocol, &targetFormat, &bitrate, &audioChannels, &createdAt)
	if err != nil {
		return types.TranscodeParamsRow{}, err
	}

	createdAtTime := logic.GetTimeFromString(createdAt)
	createdAtPlusTtl := createdAtTime.Add(time.Duration(config.TranscodeParamsTtlSeconds) * time.Second)
	if time.Now().UTC().After(createdAtPlusTtl) {
		_ = DeleteTranscodeParams(ctx, paramString)
		return types.TranscodeParamsRow{}, fmt.Errorf("transcode params expired")
	}

	return types.TranscodeParamsRow{
		ParamString:   paramString,
		MediaId:       mediaId,
		MediaType:     mediaType,
		Container:     container,
		AudioCodec:    audioCodec,
		Protocol:      protocol,
		TargetFormat:  targetFormat,
		Bitrate:       bitrate,
		AudioChannels: audioChannels,
	}, nil
}

func DeleteTranscodeParams(ctx context.Context, paramString string) error {
	query := `DELETE FROM transcode_params WHERE paramString = ?`
	_, err := DB.ExecContext(ctx, query, paramString)
	if err != nil {
		return err
	}
	return nil
}

func ClearOldTranscodeParams(ctx context.Context) error {
	query := `DELETE FROM transcode_params WHERE created_at < ?`
	_, err := DB.ExecContext(ctx, query, logic.GetTimeMinusSecondsFormatted(config.TranscodeParamsTtlSeconds))
	if err != nil {
		return err
	}
	return nil
}
