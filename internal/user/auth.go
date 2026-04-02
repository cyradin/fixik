package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = fmt.Errorf("user not found")

type tokenManager interface {
	GenerateAccessToken(userID int64) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
}

type userProvider interface {
	GetByUsername(ctx context.Context, username string) (db.User, error)
}

//nolint:gosec
type LoginResult struct {
	User         User
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userProvider userProvider
	tokens       tokenManager
}

func NewAuthService(
	userProvider userProvider,
	tokenManager tokenManager,
) *AuthService {
	return &AuthService{
		userProvider: userProvider,
		tokens:       tokenManager,
	}
}

func (a *AuthService) Login(ctx context.Context, username string, password string) (LoginResult, error) {
	user, err := a.userProvider.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return LoginResult{}, ErrUserNotFound
		}

		return LoginResult{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return LoginResult{}, ErrUserNotFound
	}

	accessToken, err := a.tokens.GenerateAccessToken(user.ID)
	if err != nil {
		return LoginResult{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := a.tokens.GenerateRefreshToken(user.ID)
	if err != nil {
		return LoginResult{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return LoginResult{
		User:         NewFromDB(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
