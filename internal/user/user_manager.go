package user

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/dict"
	"golang.org/x/crypto/bcrypt"
)

type entityProvider interface {
	GetByID(ctx context.Context, id dict.EntityID) (dict.Entity, error)
	List(ctx context.Context) ([]dict.Entity, error)
}

type userRepo interface {
	Create(ctx context.Context, u *db.User) error
	GetByID(ctx context.Context, id int64) (db.User, error)
	Update(ctx context.Context, u *db.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]db.User, error)
}

type UserManager struct {
	repo         userRepo
	roleProvider entityProvider
}

func NewUserManager(
	repo userRepo,
	roleProvider entityProvider,
) *UserManager {
	return &UserManager{
		repo:         repo,
		roleProvider: roleProvider,
	}
}

func (m *UserManager) Create(ctx context.Context, u CreateUser) (User, error) {
	hashed, err := hashPassword(u.Password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	user := db.User{
		Username: u.Username,
		Email:    u.Email,
		Password: hashed,
		TeamID:   u.TeamID,
		RoleIDs:  u.RoleIDs,
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

	if u.Email != nil {
		user.Email = *u.Email
	}

	if u.TeamID != nil {
		user.TeamID = *u.TeamID
	}

	if u.Password != nil {
		hashed, err := hashPassword(*u.Password)
		if err != nil {
			return User{}, fmt.Errorf("hash password: %w", err)
		}

		user.Password = hashed
	}

	if u.RoleIDs != nil {
		user.RoleIDs = *u.RoleIDs
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

	roles := make([]dict.Entity, 0, len(userDB.RoleIDs))
	for _, roleID := range userDB.RoleIDs {
		role, err := m.roleProvider.GetByID(ctx, dict.EntityID(roleID))
		if err != nil {
			return User{}, fmt.Errorf("get role %d: %w", roleID, err)
		}

		roles = append(roles, role)
	}

	return User{
		ID:        userDB.ID,
		Username:  userDB.Username,
		Email:     userDB.Email,
		Password:  userDB.Password,
		TeamID:    userDB.TeamID,
		Roles:     roles,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
	}, nil
}

func (m *UserManager) List(ctx context.Context, limit, offset int) ([]User, error) {
	usersDB, err := m.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	roles, err := m.roleProvider.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	roleMap := make(map[int64]dict.Entity, len(roles))
	for _, r := range roles {
		roleMap[int64(r.ID)] = r
	}

	users := make([]User, 0, len(usersDB))

	for _, u := range usersDB {
		userRoles := make([]dict.Entity, 0, len(u.RoleIDs))

		for _, roleID := range u.RoleIDs {
			role, ok := roleMap[roleID]
			if !ok {
				return nil, fmt.Errorf("role %d not found", roleID)
			}

			userRoles = append(userRoles, role)
		}

		users = append(users, User{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Password:  u.Password,
			TeamID:    u.TeamID,
			Roles:     userRoles,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return users, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}
