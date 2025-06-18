package logic

import (
	"context"
	"fmt"
)

type ContextKey string

const UserIdKey ContextKey = "userId"
const UsernameKey ContextKey = "username"

func GetUserIdFromContext(ctx context.Context) (int64, error) {
	val := ctx.Value(UserIdKey)
	userId, ok := val.(int64)
	if !ok {
		return 0, fmt.Errorf("userId missing or invalid in context")
	}
	return userId, nil
}

func GetUsernameFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(UsernameKey)
	username, ok := val.(string)
	if !ok || username == "" {
		return "", fmt.Errorf("username missing or invalid in context")
	}
	return username, nil
}
