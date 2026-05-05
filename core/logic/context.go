package logic

import (
	"context"
	"fmt"
	"zene/core/types"
)

func GetUserIdFromContext(ctx context.Context) (int, error) {
	val := ctx.Value(types.ContextKey("userId"))
	userId, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("userId missing or invalid in context")
	}
	return userId, nil
}
