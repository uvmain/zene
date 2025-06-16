package auth

import (
	"context"
	"log"
	"zene/core/database"
	"zene/core/logger"
)

var FirstUserAllowed bool

func Initialise(ctx context.Context) {
	exists, err := database.AnyUsersExist(ctx)
	if err != nil {
		log.Fatalf("Failed to check if users exist: %v", err)
	}

	FirstUserAllowed = !exists
	if FirstUserAllowed {
		logger.Println("No users found, first login will create admin user")
	}
}
