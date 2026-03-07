package incident

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type ImpactRepositorySuite struct {
	tests.PostgresSuite
	repo *ImpactRepository
}

func TestImpactRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ImpactRepositorySuite))
}

func (s *ImpactRepositorySuite) SetupTest() {
	s.repo = NewImpactRepository(s.Postgres())

	_, err := s.Postgres().Exec(s.T().Context(), `
		TRUNCATE TABLE incident_impacts RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *ImpactRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	im := &Impact{
		Code: "high",
		Name: "High",
	}

	err := s.repo.Create(ctx, im)
	s.Require().NoError(err)
	s.NotZero(im.ID)
}

func (s *ImpactRepositorySuite) TestGetByID_Found() {
	im := s.createImpact("medium", "Medium")

	fromDB, err := s.repo.GetByID(s.T().Context(), im.ID)
	s.Require().NoError(err)
	s.Equal(im.ID, fromDB.ID)
	s.Equal(im.Code, fromDB.Code)
	s.Equal(im.Name, fromDB.Name)
}

func (s *ImpactRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().NoError(err)
	s.Equal(Impact{}, fromDB)
}

func (s *ImpactRepositorySuite) TestList() {
	im1 := s.createImpact("low", "Low")
	im2 := s.createImpact("critical", "Critical")

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)
	s.Contains(list, *im1)
	s.Contains(list, *im2)
}

func (s *ImpactRepositorySuite) TestUpdate() {
	im := s.createImpact("pending", "Pending")

	im.Code = "waiting"
	im.Name = "Waiting"
	err := s.repo.Update(s.T().Context(), im)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), im.ID)
	s.Require().NoError(err)
	s.Equal("waiting", fromDB.Code)
	s.Equal("Waiting", fromDB.Name)
}

func (s *ImpactRepositorySuite) TestDelete() {
	im := s.createImpact("obsolete", "Obsolete")

	err := s.repo.Delete(s.T().Context(), im.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), im.ID)
	s.Require().NoError(err)
	s.Equal(Impact{}, fromDB)
}

func (s *ImpactRepositorySuite) createImpact(code, name string) *Impact {
	ctx := s.T().Context()
	im := &Impact{Code: code, Name: name}
	err := s.repo.Create(ctx, im)
	s.Require().NoError(err)

	return im
}
