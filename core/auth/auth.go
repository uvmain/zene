package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"zene/core/database"
	"zene/core/types"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	passedUsername := r.FormValue("username")
	passedPassword := r.FormValue("password")
	var user types.User

	log.Println("User logging in")

	// Check if users exist
	usersExist, err := database.AnyUsersExist(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !usersExist {
		// create first admin user
		hashedPassword, err := hashPassword(passedPassword)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = database.UpsertUser(ctx, types.User{
			Username:     passedUsername,
			PasswordHash: hashedPassword,
			IsAdmin:      true,
		})
		if err != nil {
			log.Printf("Error creating initial admin user: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Println("Initial admin user created")
		user, err = database.GetUserByUsername(ctx, passedUsername)
		if err != nil {
			log.Printf("Error fetching new admin user details from database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// normal login flow: verify credentials
		user, err = database.GetUserByUsername(ctx, passedUsername)
		if err != nil || !checkPasswordHash(passedPassword, user.PasswordHash) {
			log.Println("Login unsuccessful, invalid credentials")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	// Generate token, save session, set cookie as before
	token, err := generateToken()
	if err != nil {
		log.Println("Error generating token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = database.SaveSessionToken(ctx, user.Id, token, time.Hour*24*7)
	if err != nil {
		log.Println("Error saving session token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "appSession",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   604800,
	})

	sessionCheck := types.SessionCheck{
		LoggedIn: true,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sessionCheck); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("appSession")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId, valid, err := database.GetUserIDFromSession(r.Context(), cookie.Value)
		if err != nil || !valid {
			log.Println("Unauthorized access attempt:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("User logging out")
	cookie, err := r.Cookie("appSession")
	if err == nil {
		err := database.DeleteSessionToken(r.Context(), cookie.Value)
		if err != nil {
			log.Printf("Error deleting session token: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// expire cookie client-side
		http.SetCookie(w, &http.Cookie{
			Name:     "appSession",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
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
		_, isValid, err := database.GetUserIDFromSession(r.Context(), cookie.Value)
		if err == nil && isValid {
			sessionCheck.LoggedIn = true
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sessionCheck); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func GetUserIdFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value("userId").(int)
	return id, ok
}
