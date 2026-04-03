package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound        = fmt.Errorf("user not found")
	ErrUserInvalidPassword = fmt.Errorf("invalid password")
)

type tokenManager interface {
	GenerateAccessToken(userID int64) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	ValidateAccessToken(tokenStr string) (int64, error)
	ValidateRefreshToken(tokenStr string) (int64, error)
}

type userProvider interface {
	GetByUsername(ctx context.Context, username string) (db.User, error)
	GetByID(ctx context.Context, id int64) (db.User, error)
}

type userUpdater interface {
	Update(ctx context.Context, u *db.User) error
}

//nolint:gosec
type LoginResult struct {
	User         User
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userProvider userProvider
	userUpdater  userUpdater
	tokens       tokenManager
}

func NewAuthService(
	userProvider userProvider,
	userUpdater userUpdater,
	tokenManager tokenManager,
) *AuthService {
	return &AuthService{
		userProvider: userProvider,
		userUpdater:  userUpdater,
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

	accessToken, refreshToken, err := a.generateTokens(user.ID)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		User:         NewFromDB(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthService) Refresh(ctx context.Context, oldRefreshToken string) (LoginResult, error) {
	userID, err := a.tokens.ValidateRefreshToken(oldRefreshToken)
	if err != nil {
		return LoginResult{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := a.userProvider.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return LoginResult{}, ErrUserNotFound
		}

		return LoginResult{}, err
	}

	accessToken, refreshToken, err := a.generateTokens(user.ID)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		User:         NewFromDB(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthService) GetUserFromAccessToken(ctx context.Context, token string) (User, error) {
	userID, err := a.tokens.ValidateAccessToken(token)
	if err != nil {
		return User{}, fmt.Errorf("invalid access token: %w", err)
	}

	user, err := a.userProvider.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return User{}, ErrUserNotFound
		}

		return User{}, err
	}

	return NewFromDB(user), nil
}

func (m *AuthService) ChangePassword(ctx context.Context, userID int64, current, new string) error {
	u, err := m.userProvider.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return ErrUserNotFound
		}

		return fmt.Errorf("get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(current)); err != nil {
		return ErrUserInvalidPassword
	}

	hashed, err := hashPassword(new)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	u.Password = hashed

	if err := m.userUpdater.Update(ctx, &u); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}

func (a *AuthService) generateTokens(userID int64) (string, string, error) {
	accessToken, err := a.tokens.GenerateAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := a.tokens.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
