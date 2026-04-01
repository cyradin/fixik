package db

import (
	"testing"

	"github.com/cyradin/fixik/pkg/tests"
	"github.com/stretchr/testify/suite"
)

type StatusRepositorySuite struct {
	tests.PostgresSuite
	repo *StatusRepository
}

func TestStatusRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StatusRepositorySuite))
}

func (s *StatusRepositorySuite) SetupTest() {
	s.repo = NewStatusRepository(s.Postgres())

	_, err := s.Postgres().Exec(s.T().Context(), `
		TRUNCATE TABLE statuses RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *StatusRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	e := &Status{
		Code:    "done",
		Name:    "Done",
		Sort:    100,
		IsFinal: true,
	}

	err := s.repo.Create(ctx, e)
	s.Require().NoError(err)
	s.NotZero(e.ID)
	s.NotZero(e.CreatedAt)
	s.NotZero(e.UpdatedAt)
	s.Nil(e.DeletedAt)
	s.Equal(true, e.IsFinal)
	s.Equal(100, e.Sort)
}

func (s *StatusRepositorySuite) TestGetByID_Found() {
	e := s.createStatus("in_progress", "In Progress", 50, false)

	fromDB, err := s.repo.GetByID(s.T().Context(), e.ID)
	s.Require().NoError(err)
	s.Equal(e.ID, fromDB.ID)
	s.Equal(e.Code, fromDB.Code)
	s.Equal(e.Name, fromDB.Name)
	s.Equal(e.Sort, fromDB.Sort)
	s.Equal(e.IsFinal, fromDB.IsFinal)
}

func (s *StatusRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Status{}, fromDB)
}

func (s *StatusRepositorySuite) TestList() {
	e1 := s.createStatus("backlog", "Backlog", 10, false)
	e2 := s.createStatus("closed", "Closed", 100, true)

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)

	s.Equal(e1.ID, list[0].ID)
	s.Equal(e2.ID, list[1].ID)
}

func (s *StatusRepositorySuite) TestUpdate() {
	e := s.createStatus("pending", "Pending", 20, false)

	e.Code = "waiting"
	e.Name = "Waiting"
	e.Sort = 25
	e.IsFinal = true

	err := s.repo.Update(s.T().Context(), e)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), e.ID)
	s.Require().NoError(err)
	s.Equal("waiting", fromDB.Code)
	s.Equal("Waiting", fromDB.Name)
	s.Equal(25, fromDB.Sort)
	s.Equal(true, fromDB.IsFinal)
}

func (s *StatusRepositorySuite) TestDelete() {
	e := s.createStatus("obsolete", "Obsolete", 30, false)

	err := s.repo.Delete(s.T().Context(), e.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), e.ID)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Status{}, fromDB)
}

func (s *StatusRepositorySuite) TestDeleteNotFound() {
	err := s.repo.Delete(s.T().Context(), 123)
	s.Require().NoError(err)
}

func (s *StatusRepositorySuite) createStatus(code, name string, sort int, isFinal bool) *Status {
	ctx := s.T().Context()
	e := &Status{
		Code:    code,
		Name:    name,
		Sort:    sort,
		IsFinal: isFinal,
	}

	err := s.repo.Create(ctx, e)
	s.Require().NoError(err)

	return e
}
