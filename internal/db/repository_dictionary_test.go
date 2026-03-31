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
		TRUNCATE TABLE priorities RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *DictRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	e := &DictEntity{
		Code:        "high",
		Name:        "High",
		Description: "High priority level",
		Sort:        10,
	}

	err := s.repo.Create(ctx, e)
	s.Require().NoError(err)
	s.NotZero(e.ID)
	s.NotZero(e.CreatedAt)
	s.NotZero(e.UpdatedAt)
	s.Nil(e.DeletedAt)
	s.Equal(10, e.Sort)
}

func (s *DictRepositorySuite) TestGetByID_Found() {
	e := s.createEntity("medium", "Medium", "Medium priority level", 20)

	fromDB, err := s.repo.GetByID(s.T().Context(), e.ID)
	s.Require().NoError(err)
	s.Equal(e.ID, fromDB.ID)
	s.Equal(e.Code, fromDB.Code)
	s.Equal(e.Name, fromDB.Name)
	s.Equal(e.Description, fromDB.Description)
	s.Equal(e.Sort, fromDB.Sort)
}

func (s *DictRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(DictEntity{}, fromDB)
}

func (s *DictRepositorySuite) TestList() {
	e1 := s.createEntity("low", "Low", "Low priority", 5)
	e2 := s.createEntity("critical", "Critical", "Critical priority", 1)

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)

	s.Equal(e2.ID, list[0].ID)
	s.Equal(e1.ID, list[1].ID)
}

func (s *DictRepositorySuite) TestUpdate() {
	e := s.createEntity("pending", "Pending", "Pending priority", 15)

	e.Code = "waiting"
	e.Name = "Waiting"
	e.Description = "Waiting priority"
	e.Sort = 25

	err := s.repo.Update(s.T().Context(), e)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), e.ID)
	s.Require().NoError(err)
	s.Equal("waiting", fromDB.Code)
	s.Equal("Waiting", fromDB.Name)
	s.Equal(e.Description, fromDB.Description)
	s.Equal(25, fromDB.Sort)
}

func (s *DictRepositorySuite) TestDelete() {
	e := s.createEntity("obsolete", "Obsolete", "Obsolete priority", 30)

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

func (s *DictRepositorySuite) createEntity(code, name, description string, sort int) *DictEntity {
	ctx := s.T().Context()
	e := &DictEntity{
		Code:        code,
		Name:        name,
		Description: description,
		Sort:        sort,
	}

	err := s.repo.Create(ctx, e)
	s.Require().NoError(err)

	return e
}
