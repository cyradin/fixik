package user

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = db.ErrNotFound

type userRepo interface {
	Create(ctx context.Context, u *db.User) error
	GetByID(ctx context.Context, id int64) (db.User, error)
	GetByUsername(ctx context.Context, username string) (db.User, error)
	GetByIDMany(ctx context.Context, ids []int64) ([]db.User, error)
	Update(ctx context.Context, u *db.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]db.User, error)
}

type UserManager struct {
	repo userRepo
}

func NewUserManager(repo userRepo) *UserManager {
	return &UserManager{
		repo: repo,
	}
}

func (m *UserManager) Create(ctx context.Context, u CreateUser) (User, error) {
	hashed, err := hashPassword(u.Password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	user := db.User{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		Password: hashed,
		TeamID:   u.TeamID,
		Role:     u.Role,
	}

	if err := m.repo.Create(ctx, &user); err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return m.GetByID(ctx, user.ID)
}

func (m *UserManager) Update(ctx context.Context, u UpdateUser) (User, error) {
	user, err := m.repo.GetByID(ctx, u.ID)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	if u.Username != nil {
		user.Username = *u.Username
	}

	if u.Name != nil {
		user.Name = *u.Name
	}

	if u.Email != nil {
		user.Email = *u.Email
	}

	if u.TeamID != nil {
		if *u.TeamID == 0 {
			user.TeamID = nil
		} else {
			user.TeamID = u.TeamID
		}
	}

	if u.Password != nil {
		hashed, err := hashPassword(*u.Password)
		if err != nil {
			return User{}, fmt.Errorf("hash password: %w", err)
		}

		user.Password = hashed
	}

	if u.Role != nil {
		user.Role = *u.Role
	}

	if err := m.repo.Update(ctx, &user); err != nil {
		return User{}, fmt.Errorf("update user: %w", err)
	}

	return m.GetByID(ctx, user.ID)
}

func (m *UserManager) Delete(ctx context.Context, id int64) error {
	return m.repo.Delete(ctx, id)
}

func (m *UserManager) GetByID(ctx context.Context, id int64) (User, error) {
	userDB, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return m.fromDB(userDB), nil
}

func (m *UserManager) GetByUsername(ctx context.Context, username string) (User, error) {
	userDB, err := m.repo.GetByUsername(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return m.fromDB(userDB), nil
}

func (m *UserManager) GetByIDMany(ctx context.Context, ids []int64) ([]User, error) {
	usersDB, err := m.repo.GetByIDMany(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	users := make([]User, 0, len(usersDB))

	for _, u := range usersDB {
		users = append(users, m.fromDB(u))
	}

	return users, nil
}

func (m *UserManager) List(ctx context.Context, limit, offset int) ([]User, error) {
	usersDB, err := m.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	users := make([]User, 0, len(usersDB))

	for _, u := range usersDB {
		users = append(users, m.fromDB(u))
	}

	return users, nil
}

func (m *UserManager) fromDB(u db.User) User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		TeamID:    u.TeamID,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}
