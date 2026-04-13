package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cyradin/fixik/internal/role"
	"github.com/cyradin/fixik/internal/user"
	"github.com/cyradin/fixik/pkg/logger"
)

var (
	ErrUnauthorized = &UserMessageError{
		Err:    errors.New(http.StatusText(http.StatusUnauthorized)),
		Msg:    "Вы не авторизованы",
		Status: http.StatusUnauthorized,
	}

	ErrInvalidLoginPassword = &UserMessageError{
		Err:    errors.New("invalid Password"),
		Msg:    "Неверный логин или пароль",
		Status: http.StatusUnauthorized,
	}

	ErrInvalidPassword = &UserMessageError{
		Err:    errors.New("invalid Password"),
		Msg:    "Неверный пароль",
		Status: http.StatusUnauthorized,
	}

	ErrRequestPathEntityID = &UserMessageError{
		Err:    errors.New("invalid entity id provided in request path"),
		Status: http.StatusBadRequest,
	}

	ErrValidation = func(msg string) *UserMessageError {
		return &UserMessageError{
			Err:    errors.New("validation error"),
			Status: http.StatusBadRequest,
			Msg:    msg,
		}
	}

	ErrRequestDecode = func(msg string) *UserMessageError {
		return &UserMessageError{
			Err:    errors.New("request decode error"),
			Status: http.StatusBadRequest,
			Msg:    msg,
		}
	}

	ErrForbidden = &UserMessageError{
		Err:    errors.New(http.StatusText(http.StatusForbidden)),
		Status: http.StatusForbidden,
		Msg:    "Доступ запрещен",
	}
)

type userMessager interface {
	UserMessage() string
}

type httpStatuser interface {
	HTTPStatus() int
}

type UserMessageError struct {
	Status int
	Err    error
	Msg    string
}

func (e *UserMessageError) Error() string {
	return e.Err.Error()
}

func (e *UserMessageError) UserMessage() string {
	return e.Msg
}

func (e *UserMessageError) HTTPStatus() int {
	return e.Status
}

type ErrorResponse struct {
	UserMessage string `json:"userMessage" validate:"required"` // сообщение, которое нужно показать пользователю
	Error       string `json:"error" validate:"required"`       // причина ошибки
}

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if e, ok := err.(httpStatuser); ok {
		status = e.HTTPStatus()
	}

	userMsg := ""
	if e, ok := err.(userMessager); ok {
		userMsg = e.UserMessage()
	}

	if status >= http.StatusInternalServerError {
		logger.FromContext(ctx).Error("request error", logger.Error(err))
	}

	writeJSON(w, status, ErrorResponse{
		Error:       err.Error(),
		UserMessage: userMsg,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(v)
}

func checkPermissions(ctx context.Context, perms role.Permission) bool {
	u, ok := user.FromContext(ctx)
	if !ok {
		return false
	}

	return role.Can(u.Role, perms)
}
