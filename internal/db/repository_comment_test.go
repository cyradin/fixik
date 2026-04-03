package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type CommentRepositorySuite struct {
	tests.PostgresSuite
	repo *CommentRepository
}

func TestCommentRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CommentRepositorySuite))
}

func (s *CommentRepositorySuite) SetupTest() {
	s.repo = NewCommentRepository(s.Postgres())

	ctx := s.T().Context()
	_, err := s.Postgres().Exec(ctx, `
		TRUNCATE TABLE comments RESTART IDENTITY CASCADE;
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
		TRUNCATE TABLE statuses RESTART IDENTITY CASCADE;
		TRUNCATE TABLE priorities RESTART IDENTITY CASCADE;
		TRUNCATE TABLE incidents RESTART IDENTITY CASCADE;
		TRUNCATE TABLE teams RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *CommentRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	user := s.createUser()
	status := s.createStatus("status", "status")
	priority := s.createPriority("p1", "p1")
	incident := s.createIncident(status, priority, user)

	c := &Comment{
		AuthorID:   user.ID,
		IncidentID: incident.ID,
		Text:       "hello",
	}

	err := s.repo.Create(ctx, c)
	s.Require().NoError(err)

	s.NotZero(c.ID)
	s.WithinDuration(time.Now(), c.CreatedAt, time.Second)
	s.WithinDuration(time.Now(), c.UpdatedAt, time.Second)
}

func (s *CommentRepositorySuite) TestUpdate() {
	ctx := s.T().Context()

	user := s.createUser()
	status := s.createStatus("status", "status")
	priority := s.createPriority("p1", "p1")
	incident := s.createIncident(status, priority, user)

	c := s.createComment(user.ID, incident.ID, "old")

	oldUpdated := c.UpdatedAt

	c.Text = "new text"

	err := s.repo.Update(ctx, c)
	s.Require().NoError(err)

	fromDB := s.getComment(c.ID)

	s.Equal("new text", fromDB.Text)
	s.True(fromDB.UpdatedAt.After(oldUpdated))
}

func (s *CommentRepositorySuite) TestUpdate_NotFound() {
	ctx := s.T().Context()

	c := &Comment{
		ID:   999999,
		Text: "ghost",
	}

	err := s.repo.Update(ctx, c)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *CommentRepositorySuite) TestDelete_SoftDelete() {
	ctx := s.T().Context()

	user := s.createUser()
	status := s.createStatus("status", "status")
	priority := s.createPriority("p1", "p1")
	incident := s.createIncident(status, priority, user)

	c := s.createComment(user.ID, incident.ID, "text")

	err := s.repo.Delete(ctx, c.ID)
	s.Require().NoError(err)

	comments, err := s.repo.ListByIncident(ctx, incident.ID, 100, 0)
	s.Require().NoError(err)
	s.Len(comments, 0)
}

func (s *CommentRepositorySuite) TestDelete_NotFound() {
	ctx := s.T().Context()

	err := s.repo.Delete(ctx, 999999)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *CommentRepositorySuite) TestListByIncident() {
	ctx := s.T().Context()

	user := s.createUser()
	status := s.createStatus("status", "status")
	priority := s.createPriority("p1", "p1")
	incident := s.createIncident(status, priority, user)

	c1 := s.createComment(user.ID, incident.ID, "1")
	c2 := s.createComment(user.ID, incident.ID, "2")
	c3 := s.createComment(user.ID, incident.ID, "3")

	s.Run("list all", func() {
		list, err := s.repo.ListByIncident(ctx, incident.ID, 100, 0)
		s.Require().NoError(err)

		s.Len(list, 3)
		s.Equal(c1.ID, list[0].ID)
		s.Equal(c2.ID, list[1].ID)
		s.Equal(c3.ID, list[2].ID)
	})

	s.Run("limit", func() {
		list, err := s.repo.ListByIncident(ctx, incident.ID, 2, 0)
		s.Require().NoError(err)

		s.Len(list, 2)
		s.Equal(c1.ID, list[0].ID)
		s.Equal(c2.ID, list[1].ID)
	})

	s.Run("offset", func() {
		list, err := s.repo.ListByIncident(ctx, incident.ID, 100, 1)
		s.Require().NoError(err)

		s.Len(list, 2)
		s.Equal(c2.ID, list[0].ID)
		s.Equal(c3.ID, list[1].ID)
	})

	s.Run("deleted not returned", func() {
		err := s.repo.Delete(ctx, c2.ID)
		s.Require().NoError(err)

		list, err := s.repo.ListByIncident(ctx, incident.ID, 100, 0)
		s.Require().NoError(err)

		s.Len(list, 2)
		s.Equal(c1.ID, list[0].ID)
		s.Equal(c3.ID, list[1].ID)
	})
}

func (s *CommentRepositorySuite) createUser() User {
	ctx := s.T().Context()

	team := s.createTeam("t1", "Team")

	u := &User{
		Name:     "User",
		Username: "user",
		Email:    "user@test.com",
		Password: "pass",
		TeamID:   &team.ID,
		Role:     RoleUser,
	}

	err := s.repoUser().Create(ctx, u)
	s.Require().NoError(err)

	return *u
}

func (s *CommentRepositorySuite) createIncident(status Status, priority Priority, author User) *Incident {
	ctx := s.T().Context()

	incident := &Incident{
		Title:       "title",
		Description: "discription",
		PriorityID:  priority.ID,
		StatusID:    status.ID,
		TeamID:      nil,
		AuthorID:    &author.ID,
	}

	err := NewIncidentRepository(s.Postgres()).Create(ctx, incident)

	s.Require().NoError(err)

	return incident
}

func (s *CommentRepositorySuite) createComment(authorID, incidentID int64, text string) *Comment {
	ctx := s.T().Context()

	c := &Comment{
		AuthorID:   authorID,
		IncidentID: incidentID,
		Text:       text,
	}

	err := s.repo.Create(ctx, c)
	s.Require().NoError(err)

	return c
}

func (s *CommentRepositorySuite) getComment(id int64) Comment {
	ctx := s.T().Context()

	var c Comment

	err := s.Postgres().QueryRow(ctx, `
		SELECT id, author_id, incident_id, text, created_at, updated_at, deleted_at
		FROM comments
		WHERE id = $1
	`, id).Scan(
		&c.ID,
		&c.AuthorID,
		&c.IncidentID,
		&c.Text,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)

	s.Require().NoError(err)

	return c
}

func (s *CommentRepositorySuite) createTeam(code, name string) Team {
	ctx := s.T().Context()
	team := &Team{Code: code, Name: name}

	err := NewTeamRepository(s.Postgres()).Create(ctx, team)
	s.Require().NoError(err)

	return *team
}

func (s *CommentRepositorySuite) createStatus(code, name string) Status {
	ctx := s.T().Context()
	sr := NewStatusRepository(s.Postgres())

	status := &Status{Code: code, Name: name}
	err := sr.Create(ctx, status)
	s.Require().NoError(err)
	s.NotZero(status.ID)

	return *status
}

func (s *CommentRepositorySuite) createPriority(code, name string) Priority {
	ctx := s.T().Context()
	pr := NewPriorityRepository(s.Postgres())

	priority := &Priority{Code: code, Name: name}
	err := pr.Create(ctx, priority)
	s.Require().NoError(err)
	s.NotZero(priority.ID)

	return *priority
}

func (s *CommentRepositorySuite) repoUser() *UserRepository {
	return NewUserRepository(s.Postgres())
}
