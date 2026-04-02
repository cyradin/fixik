package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cyradin/fixik/internal/container"
	"github.com/cyradin/fixik/internal/user"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	maxUsernameLen = 100
	maxEmailLen    = 255
	minPasswordLen = 6
	maxPasswordLen = 255
)

type userCreater interface {
	Create(ctx context.Context, u user.CreateUser) (user.User, error)
}

type userUpdater interface {
	Update(ctx context.Context, u user.UpdateUser) (user.User, error)
}

type userDeleter interface {
	Delete(ctx context.Context, id int64) error
}

type userLister interface {
	List(ctx context.Context, limit, offset int) ([]user.User, error)
}

type userGetter interface {
	GetByID(ctx context.Context, id int64) (user.User, error)
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"` //nolint:gosec
	TeamID   *int64 `json:"teamId"`
	Role     string `json:"role" enums:"admin,manager,user" validate:"required"`
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxUsernameLen)),
		validation.Field(&r.Username, validation.Required, validation.Length(1, maxUsernameLen)),
		validation.Field(&r.Email, validation.Required, validation.Length(1, maxEmailLen)),
		validation.Field(&r.Password, validation.Required, validation.Length(minPasswordLen, maxPasswordLen)),
		validation.Field(&r.TeamID, validation.Min(1)),
		validation.Field(&r.Role, validation.Required, validation.In(toAnySlice(user.RoleTypes())...)),
	)
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"` //nolint:gosec
	TeamID   *int64  `json:"teamId"`
	Role     *string `json:"role" enums:"admin,manager,user"`
}

func (r UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Length(1, maxUsernameLen)),
		validation.Field(&r.Username, validation.Length(1, maxUsernameLen)),
		validation.Field(&r.Email, validation.Length(1, maxEmailLen)),
		validation.Field(&r.Password, validation.Length(minPasswordLen, maxPasswordLen)),
		validation.Field(&r.TeamID, validation.Min(1)),
		validation.Field(&r.Role, validation.In(toAnySlice(user.RoleTypes())...)),
	)
}

type UserResponse struct {
	ID       int64  `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	TeamID   *int64 `json:"teamId"`
	Role     string `json:"role" enums:"admin,manager,user" validate:"required"`
}

type ListUsersResponse struct {
	Items []UserResponse `json:"items" validate:"required"`
}

func userRoutes(c *container.Container) func(r chi.Router) {
	manager := c.UserManager()

	return func(r chi.Router) {
		r.Post("/", createUser(manager))
		r.Get("/", listUsers(manager))
		r.Get("/{id}", getUser(manager))
		r.Patch("/{id}", updateUser(manager))
		r.Delete("/{id}", deleteUser(manager))
	}
}

// @Summary Create user
// @Description Create new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User data"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func createUser(manager userCreater) http.HandlerFunc {
	return handle(func(ctx context.Context, req CreateUserRequest) (UserResponse, error) {
		if err := req.Validate(); err != nil {
			return UserResponse{}, err
		}

		u, err := manager.Create(ctx, user.CreateUser{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
			TeamID:   req.TeamID,
			Role:     req.Role,
		})
		if err != nil {
			return UserResponse{}, err
		}

		return toUserResponse(u), nil
	})
}

// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func getUser(manager userGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (UserResponse, error) {
			u, err := manager.GetByID(ctx, id)
			if err != nil {
				return UserResponse{}, err
			}

			return toUserResponse(u), nil
		})(w, r)
	}
}

// @Summary Update user
// @Description Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "User data"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [patch]
func updateUser(manager userUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, req UpdateUserRequest) (UserResponse, error) {
			u, err := manager.Update(ctx, user.UpdateUser{
				ID:       id,
				Username: req.Username,
				Email:    req.Email,
				Password: req.Password,
				TeamID:   req.TeamID,
				Role:     req.Role,
			})
			if err != nil {
				return UserResponse{}, err
			}

			return toUserResponse(u), nil
		})(w, r)
	}
}

// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func deleteUser(manager userDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (NoBody, error) {
			return NoBody{}, manager.Delete(ctx, id)
		})(w, r)
	}
}

// @Summary List users
// @Description Get all users with optional limit and offset
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int true "Limit of items" default(100)
// @Param offset query int true "Offset for pagination" default(0)
// @Success 200 {object} ListUsersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func listUsers(manager userLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset, err := decodePagination(r, 1, 100) //nolint:mnd
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (ListUsersResponse, error) {
			users, err := manager.List(ctx, limit, offset)
			if err != nil {
				return ListUsersResponse{}, err
			}

			resp := make([]UserResponse, len(users))
			for i, u := range users {
				resp[i] = toUserResponse(u)
			}

			return ListUsersResponse{Items: resp}, nil
		})(w, r)
	}
}

func toUserResponse(u user.User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
		TeamID:   u.TeamID,
		Role:     u.Role,
	}
}

func toAnySlice[T any](value []T) []any {
	result := make([]any, 0, len(value))

	for _, v := range value {
		result = append(result, v)
	}

	return result
}
