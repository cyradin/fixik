package incident

import (
	"context"
	"fmt"

	"github.com/cyradin/fixik/internal/db"
	"github.com/cyradin/fixik/internal/user"
)

var ErrUnauthorized = fmt.Errorf("unauthorized")

type commentRepo interface {
	Create(ctx context.Context, c *db.Comment) error
	Update(ctx context.Context, c *db.Comment) error
	Delete(ctx context.Context, id int64) error
	ListByIncident(ctx context.Context, incidentID int64, limit, offset int) (db.CommentListResult, error)
}

type CommentManager struct {
	repo         commentRepo
	userProvider userProvider
}

func NewCommentManager(repo commentRepo, userProvider userProvider) *CommentManager {
	return &CommentManager{
		repo:         repo,
		userProvider: userProvider,
	}
}

func (m *CommentManager) Create(ctx context.Context, incidentID int64, text string) (Comment, error) {
	u, ok := user.FromContext(ctx)
	if !ok {
		return Comment{}, ErrUnauthorized
	}

	dbComment := db.Comment{
		IncidentID: incidentID,
		AuthorID:   u.ID,
		Text:       text,
	}

	if err := m.repo.Create(ctx, &dbComment); err != nil {
		return Comment{}, fmt.Errorf("create comment: %w", err)
	}

	return m.fromDB(dbComment, u), nil
}

func (m *CommentManager) ListByIncident(ctx context.Context, incidentID int64, limit, offset int) (CommentList, error) {
	res, err := m.repo.ListByIncident(ctx, incidentID, limit, offset)
	if err != nil {
		return CommentList{}, fmt.Errorf("list comments: %w", err)
	}

	if len(res.Items) == 0 {
		return CommentList{
			Items: []Comment{},
			Total: res.Total,
		}, nil
	}

	ids := make([]int64, 0, len(res.Items))
	seen := make(map[int64]struct{})

	for _, c := range res.Items {
		if _, ok := seen[c.AuthorID]; !ok {
			ids = append(ids, c.AuthorID)
			seen[c.AuthorID] = struct{}{}
		}
	}

	users, err := m.userProvider.GetByIDMany(ctx, ids)
	if err != nil {
		return CommentList{}, fmt.Errorf("get users: %w", err)
	}

	userMap := make(map[int64]user.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	items := make([]Comment, 0, len(res.Items))
	for _, c := range res.Items {
		u, ok := userMap[c.AuthorID]
		if !ok {
			return CommentList{}, fmt.Errorf("user %d not found", c.AuthorID)
		}

		items = append(items, m.fromDB(c, u))
	}

	return CommentList{
		Items: items,
		Total: res.Total,
	}, nil
}

func (m *CommentManager) fromDB(c db.Comment, u user.User) Comment {
	return Comment{
		ID:         c.ID,
		IncidentID: c.IncidentID,
		Author:     u,
		Text:       c.Text,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}
