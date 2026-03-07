package incident

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type PriorityRepositorySuite struct {
	tests.PostgresSuite
	repo *PriorityRepository
}

func TestPriorityRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PriorityRepositorySuite))
}

func (s *PriorityRepositorySuite) SetupTest() {
	s.repo = NewPriorityRepository(s.Postgres())

	_, err := s.Postgres().Exec(s.T().Context(), `
		TRUNCATE TABLE incident_priorities RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *PriorityRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	p := &Priority{
		Code: "critical",
		Name: "Critical",
	}

	err := s.repo.Create(ctx, p)
	s.Require().NoError(err)
	s.NotZero(p.ID)
}

func (s *PriorityRepositorySuite) TestGetByID_Found() {
	p := s.createPriority("high", "High")

	fromDB, err := s.repo.GetByID(s.T().Context(), p.ID)
	s.Require().NoError(err)
	s.Equal(p.ID, fromDB.ID)
	s.Equal(p.Code, fromDB.Code)
	s.Equal(p.Name, fromDB.Name)
}

func (s *PriorityRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().NoError(err)
	s.Equal(Priority{}, fromDB)
}

func (s *PriorityRepositorySuite) TestList() {
	p1 := s.createPriority("low", "Low")
	p2 := s.createPriority("medium", "Medium")

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)
	s.Contains(list, *p1)
	s.Contains(list, *p2)
}

func (s *PriorityRepositorySuite) TestUpdate() {
	p := s.createPriority("pending", "Pending")

	p.Code = "waiting"
	p.Name = "Waiting"
	err := s.repo.Update(s.T().Context(), p)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), p.ID)
	s.Require().NoError(err)
	s.Equal("waiting", fromDB.Code)
	s.Equal("Waiting", fromDB.Name)
}

func (s *PriorityRepositorySuite) TestDelete() {
	p := s.createPriority("obsolete", "Obsolete")

	err := s.repo.Delete(s.T().Context(), p.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), p.ID)
	s.Require().NoError(err)
	s.Equal(Priority{}, fromDB)
}

func (s *PriorityRepositorySuite) createPriority(code, name string) *Priority {
	ctx := s.T().Context()
	p := &Priority{Code: code, Name: name}
	err := s.repo.Create(ctx, p)
	s.Require().NoError(err)

	return p
}
