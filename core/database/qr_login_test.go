package database

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"
)

func setupQRLoginTestDB(t *testing.T) {
	t.Helper()
	previousDB := DB
	db, err := sql.Open("sqlite3", "file:qr-login-test?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("opening test database: %v", err)
	}
	db.SetMaxOpenConns(1)
	DB = db
	t.Cleanup(func() {
		DB = previousDB
		_ = db.Close()
	})

	_, err = DB.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY);
		CREATE TABLE api_keys (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, api_key TEXT NOT NULL, date_created TEXT NOT NULL, last_used TEXT);
		INSERT INTO users (id) VALUES (1);
		INSERT INTO users (id) VALUES (2);
		INSERT INTO api_keys (user_id, api_key, date_created) VALUES (2, 'phone-key', 'now');`)
	if err != nil {
		t.Fatalf("creating test schema: %v", err)
	}
	createQRLoginChallengesTable(context.Background())
}

func TestQRLoginChallengeRequiresApprovalAndIssuesIndependentAPIKey(t *testing.T) {
	setupQRLoginTestDB(t)
	ctx := context.Background()

	challenge, err := CreateQRLoginChallenge(ctx)
	if err != nil {
		t.Fatalf("creating challenge: %v", err)
	}
	if challenge.PhoneSecret == challenge.ComputerSecret {
		t.Fatal("phone and computer secrets must be independent")
	}
	if _, err := CompleteQRLoginChallenge(ctx, challenge.ComputerSecret); err == nil {
		t.Fatal("completion before approval succeeded")
	}
	if _, err := ClaimQRLoginChallenge(ctx, "invalid"); err == nil {
		t.Fatal("invalid phone secret was accepted")
	}
	if _, err := ClaimQRLoginChallenge(ctx, challenge.PhoneSecret); err != nil {
		t.Fatalf("claiming challenge: %v", err)
	}
	if _, err := DecideQRLoginChallenge(ctx, challenge.PhoneSecret, 2, true); err != nil {
		t.Fatalf("approving challenge: %v", err)
	}

	apiKey, err := CompleteQRLoginChallenge(ctx, challenge.ComputerSecret)
	if err != nil {
		t.Fatalf("completing challenge: %v", err)
	}
	var userID int
	if err := DB.QueryRow(`SELECT user_id FROM api_keys WHERE api_key = ?`, apiKey).Scan(&userID); err != nil {
		t.Fatalf("looking up issued API key: %v", err)
	}
	if userID != 2 {
		t.Fatalf("issued API key belongs to user %d, want approving user 2", userID)
	}
	if apiKey == "phone-key" {
		t.Fatal("completion reused the phone API key")
	}
	var phoneKeyCount int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM api_keys WHERE api_key = 'phone-key'`).Scan(&phoneKeyCount); err != nil {
		t.Fatalf("checking existing phone API key: %v", err)
	}
	if phoneKeyCount != 1 {
		t.Fatal("completion invalidated the existing phone API key")
	}
	if _, err := CompleteQRLoginChallenge(ctx, challenge.ComputerSecret); err == nil {
		t.Fatal("challenge was consumed more than once")
	}
}

func TestQRLoginChallengeExpires(t *testing.T) {
	setupQRLoginTestDB(t)
	ctx := context.Background()
	challenge, err := CreateQRLoginChallenge(ctx)
	if err != nil {
		t.Fatalf("creating challenge: %v", err)
	}
	if _, err := DB.Exec(`UPDATE qr_login_challenges SET expires_at = ?`, time.Now().Add(-time.Second).Unix()); err != nil {
		t.Fatalf("expiring challenge: %v", err)
	}
	status, err := GetQRLoginChallengeStatus(ctx, challenge.ComputerSecret)
	if err != nil {
		t.Fatalf("getting expired challenge: %v", err)
	}
	if status.State != "expired" {
		t.Fatalf("challenge state = %q, want expired", status.State)
	}
	if _, err := ClaimQRLoginChallenge(ctx, challenge.PhoneSecret); err == nil {
		t.Fatal("expired challenge was claimed")
	}
}

func TestQRLoginChallengeConcurrentCompletionOnlyIssuesOneAPIKey(t *testing.T) {
	setupQRLoginTestDB(t)
	ctx := context.Background()
	challenge, err := CreateQRLoginChallenge(ctx)
	if err != nil {
		t.Fatalf("creating challenge: %v", err)
	}
	if _, err := ClaimQRLoginChallenge(ctx, challenge.PhoneSecret); err != nil {
		t.Fatalf("claiming challenge: %v", err)
	}
	if _, err := DecideQRLoginChallenge(ctx, challenge.PhoneSecret, 1, true); err != nil {
		t.Fatalf("approving challenge: %v", err)
	}

	var successes int
	var lock sync.Mutex
	var wait sync.WaitGroup
	for range 2 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			if _, err := CompleteQRLoginChallenge(ctx, challenge.ComputerSecret); err == nil {
				lock.Lock()
				successes++
				lock.Unlock()
			}
		}()
	}
	wait.Wait()
	if successes != 1 {
		t.Fatalf("completion successes = %d, want 1", successes)
	}
}
