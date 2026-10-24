package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
	"zene/core/database"
	znet "zene/core/net"
	"zene/core/types"
)

type qrLoginRateWindow struct {
	StartedAt time.Time
	Requests  int
}

var qrLoginRateLimits = struct {
	sync.Mutex
	byOperation map[string]qrLoginRateWindow
}{byOperation: make(map[string]qrLoginRateWindow)}

func HandleCreateQRLoginChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeQRLoginError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !allowQRLoginRequest(r, "create", 10, time.Minute) {
		writeQRLoginError(w, http.StatusTooManyRequests, "too many login requests")
		return
	}
	challenge, err := database.CreateQRLoginChallenge(r.Context())
	if err != nil {
		writeQRLoginError(w, http.StatusInternalServerError, "could not create login request")
		return
	}
	writeQRLoginJSON(w, http.StatusCreated, map[string]any{
		"phoneSecret":    challenge.PhoneSecret,
		"computerSecret": challenge.ComputerSecret,
		"expiresAt":      challenge.ExpiresAt.Format(time.RFC3339),
	})
}

func HandleClaimQRLoginChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeQRLoginError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !allowQRLoginRequest(r, "claim", 30, time.Minute) {
		writeQRLoginError(w, http.StatusTooManyRequests, "too many login requests")
		return
	}
	secret := znet.NormalisedForm(r, w)["phonesecret"]
	challenge, err := database.ClaimQRLoginChallenge(r.Context(), secret)
	if err != nil {
		writeQRLoginError(w, http.StatusNotFound, "login request is unavailable")
		return
	}
	writeQRLoginJSON(w, http.StatusOK, qrLoginChallengeResponse(challenge))
}

func HandleDecideQRLoginChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeQRLoginError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !allowQRLoginRequest(r, "decide", 20, time.Minute) {
		writeQRLoginError(w, http.StatusTooManyRequests, "too many login requests")
		return
	}
	form := znet.NormalisedForm(r, w)
	decision := form["decision"]
	if decision != "approve" && decision != "deny" {
		writeQRLoginError(w, http.StatusBadRequest, "a valid decision is required")
		return
	}
	userID, ok := r.Context().Value(types.ContextKey("userId")).(int)
	if !ok || userID <= 0 {
		writeQRLoginError(w, http.StatusUnauthorized, "authentication is required")
		return
	}
	challenge, err := database.DecideQRLoginChallenge(r.Context(), form["phonesecret"], userID, decision == "approve")
	if err != nil {
		writeQRLoginError(w, http.StatusConflict, "login request is unavailable")
		return
	}
	writeQRLoginJSON(w, http.StatusOK, qrLoginChallengeResponse(challenge))
}

func HandleQRLoginChallengeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeQRLoginError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !allowQRLoginRequest(r, "status", 90, time.Minute) {
		writeQRLoginError(w, http.StatusTooManyRequests, "too many login requests")
		return
	}
	secret := znet.NormalisedForm(r, w)["computersecret"]
	challenge, err := database.GetQRLoginChallengeStatus(r.Context(), secret)
	if err != nil {
		writeQRLoginError(w, http.StatusNotFound, "login request is unavailable")
		return
	}
	writeQRLoginJSON(w, http.StatusOK, qrLoginChallengeResponse(challenge))
}

func HandleCompleteQRLoginChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeQRLoginError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !allowQRLoginRequest(r, "complete", 20, time.Minute) {
		writeQRLoginError(w, http.StatusTooManyRequests, "too many login requests")
		return
	}
	apiKey, err := database.CompleteQRLoginChallenge(r.Context(), znet.NormalisedForm(r, w)["computersecret"])
	if err != nil {
		writeQRLoginError(w, http.StatusConflict, "login request is unavailable")
		return
	}
	writeQRLoginJSON(w, http.StatusOK, map[string]string{"apiKey": apiKey})
}

func qrLoginChallengeResponse(challenge database.QRLoginChallenge) map[string]string {
	return map[string]string{
		"state":     challenge.State,
		"expiresAt": challenge.ExpiresAt.Format(time.RFC3339),
	}
}

func writeQRLoginJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeQRLoginError(w http.ResponseWriter, status int, message string) {
	writeQRLoginJSON(w, status, map[string]string{"error": message})
}

func allowQRLoginRequest(r *http.Request, operation string, limit int, interval time.Duration) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	key := operation + ":" + host
	now := time.Now()
	qrLoginRateLimits.Lock()
	defer qrLoginRateLimits.Unlock()
	window := qrLoginRateLimits.byOperation[key]
	if now.Sub(window.StartedAt) >= interval {
		window = qrLoginRateWindow{StartedAt: now}
	}
	if window.Requests >= limit {
		qrLoginRateLimits.byOperation[key] = window
		return false
	}
	window.Requests++
	qrLoginRateLimits.byOperation[key] = window
	return true
}
