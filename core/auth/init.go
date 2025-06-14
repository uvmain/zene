package auth

import (
	"context"
	"log"
	"zene/core/database"
)

var FirstUserAllowed bool

func Initialise(ctx context.Context) {
	exists, err := database.AnyUsersExist(ctx)
	if err != nil {
		log.Fatalf("Failed to check if users exist: %v", err)
	}

	FirstUserAllowed = !exists
	if FirstUserAllowed {
		log.Println("No users found, first login will create admin user")
	}
}
