package db

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type DictRepositorySuite struct {
	tests.PostgresSuite
	repo *DictRepository
}

func TestDictRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DictRepositorySuite))
}

func (s *DictRepositorySuite) SetupTest() {
	s.repo = NewPriorityRepository(s.Postgres())

	_, err := s.Postgres().Exec(s.T().Context(), `
		TRUNCATE TABLE incident_priorities RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *DictRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	e := &DictEntity{
		Code:        "high",
		Name:        "High",
		Description: "High priority level",
	}

	err := s.repo.Create(ctx, e)
	s.Require().NoError(err)
	s.NotZero(e.ID)
	s.NotZero(e.CreatedAt)
	s.NotZero(e.UpdatedAt)
	s.Nil(e.DeletedAt)
}

func (s *DictRepositorySuite) TestGetByID_Found() {
	e := s.createEntity("medium", "Medium", "Medium priority level")

	fromDB, err := s.repo.GetByID(s.T().Context(), e.ID)
	s.Require().NoError(err)
	s.Equal(e.ID, fromDB.ID)
	s.Equal(e.Code, fromDB.Code)
	s.Equal(e.Name, fromDB.Name)
	s.Equal(e.Description, fromDB.Description)
}

func (s *DictRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(DictEntity{}, fromDB)
}

func (s *DictRepositorySuite) TestList() {
	e1 := s.createEntity("low", "Low", "Low priority")
	e2 := s.createEntity("critical", "Critical", "Critical priority")

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)
	s.Contains(list, *e1)
	s.Contains(list, *e2)
}

func (s *DictRepositorySuite) TestUpdate() {
	e := s.createEntity("pending", "Pending", "Pending priority")

	e.Code = "waiting"
	e.Name = "Waiting"
	e.Description = "Waiting priority"
	err := s.repo.Update(s.T().Context(), e)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), e.ID)
	s.Require().NoError(err)
	s.Equal("waiting", fromDB.Code)
	s.Equal("Waiting", fromDB.Name)
	s.Equal(e.Description, fromDB.Description)
}

func (s *DictRepositorySuite) TestDelete() {
	e := s.createEntity("obsolete", "Obsolete", "Obsolete priority")

	err := s.repo.Delete(s.T().Context(), e.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), e.ID)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(DictEntity{}, fromDB)
}

func (s *DictRepositorySuite) TestDeleteNotFound() {
	err := s.repo.Delete(s.T().Context(), 123)
	s.Require().NoError(err)
}

func (s *DictRepositorySuite) createEntity(code, name, description string) *DictEntity {
	ctx := s.T().Context()
	e := &DictEntity{
		Code:        code,
		Name:        name,
		Description: description,
	}

	err := s.repo.Create(ctx, e)
	s.Require().NoError(err)

	return e
}
