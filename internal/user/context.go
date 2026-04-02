package user

import "context"

type contextKey string

const userCtxKey = contextKey("user")

// WithContext возвращает новый контекст с пользователем
func (u User) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

// FromContext достаёт пользователя из контекста
func FromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(userCtxKey).(User)

	return u, ok
}
