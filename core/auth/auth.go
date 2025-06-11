package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"zene/core/database"
	"zene/core/types"
)

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := os.Getenv("ADMIN_USER")
	password := os.Getenv("ADMIN_PASSWORD")
	passedUsername := r.FormValue("username")
	passedPassword := r.FormValue("password")
	log.Println("User logging in")
	if passedUsername == username && passedPassword == password {
		token, err := generateToken()
		if err != nil {
			log.Println("Error generating token:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		database.SaveSessionToken(ctx, token, time.Hour*24*7)
		http.SetCookie(w, &http.Cookie{
			Name:     "appSession",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Path:     "/",
			MaxAge:   604800,
		})
		log.Println("Login successful")

		sessionCheck := types.SessionCheck{
			LoggedIn: true,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(sessionCheck); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

	} else {
		log.Println("Login unsuccessful, invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("appSession")
		if err != nil || !database.IsSessionValid(r.Context(), cookie.Value) {
			log.Println("Unauthorized access attempt")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("User logging out")
	cookie, err := r.Cookie("appSession")
	if err == nil {
		database.DeleteSessionToken(r.Context(), cookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:   "appSession",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		})
	}
	sessionCheck := types.SessionCheck{
		LoggedIn: false,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sessionCheck); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionCheck := types.SessionCheck{
		LoggedIn: false,
	}
	cookie, err := r.Cookie("appSession")
	if err == nil {
		if database.IsSessionValid(r.Context(), cookie.Value) {
			sessionCheck.LoggedIn = true
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sessionCheck); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
