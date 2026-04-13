package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/cyradin/fixik/internal/container"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func NewRouter(c *container.Container, allowedOriginsCORS []string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOriginsCORS,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, //nolint:mnd
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", authRoutes(c))

		// все остальные роуты требуют аутентификацию
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(c.AuthService()))

			r.Route("/users", userRoutes(c))
			r.Route("/statuses", statusRoutes(c))
			r.Route("/priorities", priorityRoutes(c))
			r.Route("/teams", teamRoutes(c))
			r.Route("/roles", roleRoutes(c))
			r.Route("/incidents", incidentRoutes(c))
		})
	})

	r.Handle("/*", http.FileServer(
		http.FS(
			newNoDirFs(getStaticFS()),
		),
	))

	return r
}

func decodeJSON(r *http.Request, v any) error {
	defer func() { _ = r.Body.Close() }()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(v)
}

func handle[Req any, Resp any](fn func(ctx context.Context, req Req) (Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req

		if _, ok := any(req).(NoBody); !ok {
			if err := decodeJSON(r, &req); err != nil {
				handleError(r.Context(), w, ErrRequestDecode(err.Error()))
				return
			}
		}

		// optional validation
		if v, ok := any(req).(Validatable); ok {
			if err := v.Validate(); err != nil {
				handleError(r.Context(), w, ErrValidation(err.Error()))
				return
			}
		}

		ctx := r.Context()

		resp, err := fn(ctx, req)
		if err != nil {
			handleError(ctx, w, err)
			return
		}

		if _, ok := any(resp).(NoBody); !ok {
			writeJSON(w, http.StatusOK, resp)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

type Validatable interface {
	Validate() error
}

type NoBody struct{}

func decodePagination(r *http.Request, minLimit int, maxLimit int) (int, int, error) {
	type pagination struct {
		Limit  int
		Offset int
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, offset := 0, 0

	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, fmt.Errorf("parse 'limit': %w", err)
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, fmt.Errorf("parse 'offset': %w", err)
		}
	}

	p := pagination{Limit: limit, Offset: offset}

	if err := validation.ValidateStruct(
		&p,
		validation.Field(&p.Limit, validation.Required, validation.Min(minLimit), validation.Max(maxLimit)),
		validation.Field(&p.Offset, validation.Min(0)),
	); err != nil {
		return 0, 0, err
	}

	return limit, offset, nil
}

// noDirFS prevents directory listing with go file server
type noDirFS struct {
	fs.FS
}

func newNoDirFs(fs fs.FS) *noDirFS {
	return &noDirFS{
		FS: fs,
	}
}

func (n noDirFS) Open(name string) (fs.File, error) {
	f, err := n.FS.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if info.Name() != "static" && info.IsDir() {
		return nil, fs.ErrNotExist
	}

	return f, nil
}
