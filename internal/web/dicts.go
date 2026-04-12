package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cyradin/fixik/internal/dict"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	maxNameLen = 100
	maxCodeLen = 50
)

type entityManager interface {
	Create(ctx context.Context, e dict.Entity) (dict.Entity, error)
	GetByID(ctx context.Context, id int64) (dict.Entity, error)
	List(ctx context.Context) ([]dict.Entity, error)
	Update(ctx context.Context, s dict.Entity) (dict.Entity, error)
	Delete(ctx context.Context, id int64) error
}

type CreateDictEntityRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

func (r CreateDictEntityRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxNameLen)),
		validation.Field(&r.Code, validation.Required, validation.Length(1, maxCodeLen)),
		validation.Field(&r.Sort, validation.Required),
	)
}

type UpdateDictEntityRequest struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

func (r UpdateDictEntityRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, maxNameLen)),
		validation.Field(&r.Code, validation.Required, validation.Length(1, maxCodeLen)),
		validation.Field(&r.Sort, validation.Required),
	)
}

type ListDictEntitiesResponse struct {
	Items []DictEntity `json:"items" validate:"required"`
}

type DictEntity struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Sort        int    `json:"sort" validate:"required"`
}

func createDictEntity(manager entityManager) http.HandlerFunc {
	return handle(func(ctx context.Context, req CreateDictEntityRequest) (DictEntity, error) {
		entity := dict.Entity{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			Sort:        req.Sort,
		}

		result, err := manager.Create(ctx, entity)
		if err != nil {
			return DictEntity{}, err
		}

		return toDictEntity(result), nil
	})
}

func getDictEntity(manager entityManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (DictEntity, error) {
			result, err := manager.GetByID(ctx, id)
			if err != nil {
				return DictEntity{}, err
			}

			return toDictEntity(result), nil
		})(w, r)
	}
}

func updateDictEntity(manager entityManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, req UpdateDictEntityRequest) (DictEntity, error) {
			entity := dict.Entity{
				ID:          id,
				Name:        req.Name,
				Code:        req.Code,
				Description: req.Description,
				Sort:        req.Sort,
			}

			result, err := manager.Update(ctx, entity)
			if err != nil {
				return DictEntity{}, err
			}

			return toDictEntity(result), nil
		})(w, r)
	}
}

func deleteDictEntity(manager entityManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		handle(func(ctx context.Context, _ NoBody) (NoBody, error) {
			if err := manager.Delete(ctx, id); err != nil {
				return NoBody{}, err
			}

			return NoBody{}, nil
		})(w, r)
	}
}

func listDictEntities(manager entityManager) http.HandlerFunc {
	return handle(func(ctx context.Context, _ NoBody) (ListDictEntitiesResponse, error) {
		items, err := manager.List(ctx)
		if err != nil {
			return ListDictEntitiesResponse{}, err
		}

		resp := make([]DictEntity, 0, len(items))
		for _, item := range items {
			resp = append(resp, toDictEntity(item))
		}

		return ListDictEntitiesResponse{Items: resp}, nil
	})
}

func toDictEntity(item dict.Entity) DictEntity {
	return DictEntity{
		ID:          item.ID,
		Name:        item.Name,
		Code:        item.Code,
		Description: item.Description,
		Sort:        item.Sort,
	}
}
