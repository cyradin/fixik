package db

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
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
		TRUNCATE TABLE incident_statuses RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *StatusRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	status := &Status{
		Code: "open",
		Name: "Open",
	}

	err := s.repo.Create(ctx, status)
	s.Require().NoError(err)
	s.NotZero(status.ID)
}

func (s *StatusRepositorySuite) TestGetByID_Found() {
	status := s.createStatus("in_progress", "In Progress")

	fromDB, err := s.repo.GetByID(s.T().Context(), status.ID)
	s.Require().NoError(err)
	s.Equal(status.ID, fromDB.ID)
	s.Equal(status.Code, fromDB.Code)
	s.Equal(status.Name, fromDB.Name)
}

func (s *StatusRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Status{}, fromDB)
}

func (s *StatusRepositorySuite) TestList() {
	s1 := s.createStatus("open", "Open")
	s2 := s.createStatus("closed", "Closed")

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)
	s.Contains(list, *s1)
	s.Contains(list, *s2)
}

func (s *StatusRepositorySuite) TestUpdate() {
	status := s.createStatus("pending", "Pending")

	status.Code = "waiting"
	status.Name = "Waiting"
	err := s.repo.Update(s.T().Context(), status)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), status.ID)
	s.Require().NoError(err)
	s.Equal("waiting", fromDB.Code)
	s.Equal("Waiting", fromDB.Name)
}

func (s *StatusRepositorySuite) TestDelete() {
	status := s.createStatus("obsolete", "Obsolete")

	err := s.repo.Delete(s.T().Context(), status.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), status.ID)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Status{}, fromDB)
}

func (s *StatusRepositorySuite) TestDeleteNotFound() {
	err := s.repo.Delete(s.T().Context(), 123)
	s.Require().NoError(err)
}

func (s *StatusRepositorySuite) createStatus(code, name string) *Status {
	ctx := s.T().Context()
	status := &Status{Code: code, Name: name}
	err := s.repo.Create(ctx, status)
	s.Require().NoError(err)

	return status
}
