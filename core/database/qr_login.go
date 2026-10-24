package database

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"
	"zene/core/logic"
)

const qrLoginChallengeLifetime = 2 * time.Minute

type QRLoginChallenge struct {
	State     string
	ExpiresAt time.Time
}

type NewQRLoginChallenge struct {
	PhoneSecret    string
	ComputerSecret string
	ExpiresAt      time.Time
}

func createQRLoginChallengesTable(ctx context.Context) {
	schema := `CREATE TABLE qr_login_challenges (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		phone_secret_hash TEXT NOT NULL UNIQUE,
		computer_secret_hash TEXT NOT NULL UNIQUE,
		state TEXT NOT NULL,
		approved_user_id INTEGER,
		expires_at INTEGER NOT NULL,
		created_at INTEGER NOT NULL,
		FOREIGN KEY (approved_user_id) REFERENCES users(id) ON DELETE SET NULL
	);`
	createTable(ctx, schema)
	createIndex(ctx, "idx_qr_login_challenges_expiry", "qr_login_challenges", []string{"expires_at"}, false)
}

func CreateQRLoginChallenge(ctx context.Context) (NewQRLoginChallenge, error) {
	phoneSecret, err := newQRLoginSecret()
	if err != nil {
		return NewQRLoginChallenge{}, err
	}
	computerSecret, err := newQRLoginSecret()
	if err != nil {
		return NewQRLoginChallenge{}, err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(qrLoginChallengeLifetime)
	_, err = DB.ExecContext(ctx, `DELETE FROM qr_login_challenges WHERE expires_at < ?`, now.Add(-5*time.Minute).Unix())
	if err != nil {
		return NewQRLoginChallenge{}, fmt.Errorf("cleaning expired QR login challenges: %w", err)
	}
	_, err = DB.ExecContext(ctx, `INSERT INTO qr_login_challenges
		(phone_secret_hash, computer_secret_hash, state, expires_at, created_at)
		VALUES (?, ?, 'pending', ?, ?)`,
		hashQRLoginSecret(phoneSecret), hashQRLoginSecret(computerSecret), expiresAt.Unix(), now.Unix())
	if err != nil {
		return NewQRLoginChallenge{}, fmt.Errorf("creating QR login challenge: %w", err)
	}
	return NewQRLoginChallenge{PhoneSecret: phoneSecret, ComputerSecret: computerSecret, ExpiresAt: expiresAt}, nil
}

func ClaimQRLoginChallenge(ctx context.Context, phoneSecret string) (QRLoginChallenge, error) {
	challenge, err := updateQRLoginChallengeState(ctx, phoneSecret, 0, "scanned")
	if err != nil {
		return QRLoginChallenge{}, err
	}
	if challenge.State != "scanned" {
		return QRLoginChallenge{}, fmt.Errorf("QR login challenge is unavailable")
	}
	return challenge, nil
}

func DecideQRLoginChallenge(ctx context.Context, phoneSecret string, userID int, approve bool) (QRLoginChallenge, error) {
	state := "denied"
	if approve {
		state = "approved"
	}
	current, err := getQRLoginChallenge(ctx, "phone_secret_hash", hashQRLoginSecret(phoneSecret))
	if err != nil {
		return QRLoginChallenge{}, err
	}
	if current.State != "scanned" {
		return QRLoginChallenge{}, fmt.Errorf("QR login challenge is unavailable")
	}
	return updateQRLoginChallengeState(ctx, phoneSecret, userID, state)
}

func GetQRLoginChallengeStatus(ctx context.Context, computerSecret string) (QRLoginChallenge, error) {
	return getQRLoginChallenge(ctx, "computer_secret_hash", hashQRLoginSecret(computerSecret))
}

func CompleteQRLoginChallenge(ctx context.Context, computerSecret string) (string, error) {
	now := time.Now().UTC()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("starting QR login completion: %w", err)
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRowContext(ctx, `SELECT approved_user_id FROM qr_login_challenges
		WHERE computer_secret_hash = ? AND state = 'approved' AND expires_at > ?`,
		hashQRLoginSecret(computerSecret), now.Unix()).Scan(&userID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("QR login challenge is not approved")
	}
	if err != nil {
		return "", fmt.Errorf("looking up QR login challenge: %w", err)
	}

	result, err := tx.ExecContext(ctx, `UPDATE qr_login_challenges SET state = 'consumed'
		WHERE computer_secret_hash = ? AND state = 'approved' AND expires_at > ?`,
		hashQRLoginSecret(computerSecret), now.Unix())
	if err != nil {
		return "", fmt.Errorf("consuming QR login challenge: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("checking QR login challenge consumption: %w", err)
	}
	if affected != 1 {
		return "", fmt.Errorf("QR login challenge has already been consumed")
	}

	apiKey, err := logic.GenerateNewApiKey()
	if err != nil {
		return "", fmt.Errorf("generating QR login API key: %w", err)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO api_keys (user_id, api_key, date_created) VALUES (?, ?, ?)`,
		userID, apiKey, logic.GetCurrentTimeFormatted())
	if err != nil {
		return "", fmt.Errorf("creating QR login API key: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("committing QR login completion: %w", err)
	}
	return apiKey, nil
}

func updateQRLoginChallengeState(ctx context.Context, phoneSecret string, userID int, requestedState string) (QRLoginChallenge, error) {
	now := time.Now().UTC()
	phoneHash := hashQRLoginSecret(phoneSecret)
	var query string
	var args []any
	switch requestedState {
	case "scanned":
		query = `UPDATE qr_login_challenges SET state = 'scanned'
			WHERE phone_secret_hash = ? AND state = 'pending' AND expires_at > ?`
		args = []any{phoneHash, now.Unix()}
	case "approved", "denied":
		query = `UPDATE qr_login_challenges SET state = ?, approved_user_id = ?
			WHERE phone_secret_hash = ? AND state = 'scanned' AND expires_at > ?`
		args = []any{requestedState, userID, phoneHash, now.Unix()}
	default:
		return QRLoginChallenge{}, fmt.Errorf("invalid QR login challenge state")
	}

	result, err := DB.ExecContext(ctx, query, args...)
	if err != nil {
		return QRLoginChallenge{}, fmt.Errorf("updating QR login challenge: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return QRLoginChallenge{}, fmt.Errorf("checking QR login challenge state: %w", err)
	}
	if affected != 1 {
		return getQRLoginChallenge(ctx, "phone_secret_hash", phoneHash)
	}
	return getQRLoginChallenge(ctx, "phone_secret_hash", phoneHash)
}

func getQRLoginChallenge(ctx context.Context, column, secretHash string) (QRLoginChallenge, error) {
	if column != "phone_secret_hash" && column != "computer_secret_hash" {
		return QRLoginChallenge{}, fmt.Errorf("invalid QR login challenge lookup")
	}
	var state string
	var expiresAt int64
	err := DB.QueryRowContext(ctx, `SELECT state, expires_at FROM qr_login_challenges WHERE `+column+` = ?`, secretHash).Scan(&state, &expiresAt)
	if err == sql.ErrNoRows {
		return QRLoginChallenge{}, fmt.Errorf("QR login challenge not found")
	}
	if err != nil {
		return QRLoginChallenge{}, fmt.Errorf("getting QR login challenge: %w", err)
	}
	if time.Now().UTC().Unix() >= expiresAt && state != "consumed" {
		_, _ = DB.ExecContext(ctx, `UPDATE qr_login_challenges SET state = 'expired' WHERE `+column+` = ? AND state != 'consumed'`, secretHash)
		state = "expired"
	}
	return QRLoginChallenge{State: state, ExpiresAt: time.Unix(expiresAt, 0).UTC()}, nil
}

func newQRLoginSecret() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("generating QR login secret: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func hashQRLoginSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
