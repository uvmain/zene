package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"zene/core/database"
	"zene/core/logger"
	"zene/core/logic"
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

	logger.Println("User logging in")

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
			logger.Printf("Error hashing password: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		newId, err := database.UpsertUser(ctx, types.User{
			Username:     passedUsername,
			PasswordHash: hashedPassword,
			IsAdmin:      true,
		})
		if err != nil {
			logger.Printf("Error creating initial admin user: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		logger.Printf("Initial admin user created: ID %d", newId)
		user, err = database.GetUserByUsername(ctx, passedUsername)
		if err != nil {
			logger.Printf("Error fetching new admin user details from database: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// normal login flow: verify credentials
		user, err = database.GetUserByUsername(ctx, passedUsername)
		if err != nil || !checkPasswordHash(passedPassword, user.PasswordHash) {
			logger.Println("Login unsuccessful, invalid credentials")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	// Generate token, save session, set cookie as before
	token, err := generateToken()
	if err != nil {
		logger.Println("Error generating token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = database.SaveSessionToken(ctx, user.Id, token, time.Hour*24*7)
	if err != nil {
		logger.Println("Error saving session token:", err)
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
		logger.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		ctx := r.Context()

		if token != "" {
			isValidToken, err := ValidateToken(ctx, token)
			logger.Printf("%s accessed with temporary token", r.RequestURI)
			if err != nil || !isValidToken {
				logger.Printf("Unauthorized access attempt, invalid token: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		} else {
			user, isValidSession, err := GetUserFromRequest(r)
			if err != nil || !isValidSession {
				logger.Printf("Unauthorized access attempt, invalid session: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, logic.UserIdKey, user.Id)
			ctx = context.WithValue(ctx, logic.UsernameKey, user.Username)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	})
}

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, isValidSession, err := GetUserFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !user.IsAdmin {
			logger.Printf("Unauthorized access attempt, user is not an admin: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !isValidSession {
			logger.Printf("Unauthorized access attempt, invalid session: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, logic.UserIdKey, user.Id)
		ctx = context.WithValue(ctx, logic.UsernameKey, user.Username)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("User logging out")
	cookie, err := r.Cookie("appSession")
	if err == nil {
		err := database.DeleteSessionToken(r.Context(), cookie.Value)
		if err != nil {
			logger.Printf("Error deleting session token: %v", err)
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
		logger.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionCheck := types.SessionCheck{
		LoggedIn: false,
	}

	cookie, err := r.Cookie("appSession")
	if err == nil {
		_, isValid, err := database.GetUserIdFromSession(r.Context(), cookie.Value)
		if err == nil && isValid {
			sessionCheck.LoggedIn = true
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sessionCheck); err != nil {
		logger.Println("Error encoding database response:", err)
		http.Error(w, "Error encoding database response", http.StatusInternalServerError)
		return
	}
}

// returns types.User, isValidSession bool, error
func GetUserFromRequest(r *http.Request) (types.User, bool, error) {
	cookie, err := r.Cookie("appSession")
	if err == nil {
		id, isValid, err := database.GetUserIdFromSession(r.Context(), cookie.Value)
		if err == nil && isValid {
			user, err := database.GetUserById(r.Context(), id)
			if err == nil && user.Username != "" {
				return user, isValid, nil
			}
		}
	}
	return types.User{}, false, fmt.Errorf("Failed to get user from request: %v", err)
}
