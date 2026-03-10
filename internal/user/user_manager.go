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

func (m *UserManager) Create(ctx context.Context, cmd CreateUser) (User, error) {
	hashed, err := hashPassword(cmd.Password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	user := db.User{
		Username: cmd.Username,
		Email:    cmd.Email,
		Password: hashed,
		TeamID:   cmd.TeamID,
		RoleIDs:  cmd.RoleIDs,
	}

	if err := m.repo.Create(ctx, &user); err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return m.GetByID(ctx, user.ID)
}

func (m *UserManager) Update(ctx context.Context, cmd UpdateUser) (User, error) {
	user, err := m.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	if cmd.Username != nil {
		user.Username = *cmd.Username
	}

	if cmd.Email != nil {
		user.Email = *cmd.Email
	}

	if cmd.TeamID != nil {
		user.TeamID = *cmd.TeamID
	}

	if cmd.Password != nil {
		hashed, err := hashPassword(*cmd.Password)
		if err != nil {
			return User{}, fmt.Errorf("hash password: %w", err)
		}

		user.Password = hashed
	}

	if cmd.RoleIDs != nil {
		user.RoleIDs = *cmd.RoleIDs
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
		RoleIDs:   roles,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}
