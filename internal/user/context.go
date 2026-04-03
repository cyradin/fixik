package user

import "context"

type contextKey string

const userCtxKey = contextKey("user")

func WithContext(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

func FromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(userCtxKey).(User)

	return u, ok
}
