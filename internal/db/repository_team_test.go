package db

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type TeamRepositorySuite struct {
	tests.PostgresSuite
	repo *TeamRepository
}

func TestTeamRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TeamRepositorySuite))
}

func (s *TeamRepositorySuite) SetupTest() {
	s.repo = NewTeamRepository(s.Postgres())

	_, err := s.Postgres().Exec(s.T().Context(), `
		TRUNCATE TABLE teams RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *TeamRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	team := &Team{
		Code: "core",
		Name: "Core Team",
	}

	err := s.repo.Create(ctx, team)
	s.Require().NoError(err)
	s.NotZero(team.ID)
}

func (s *TeamRepositorySuite) TestGetByID_Found() {
	team := s.createTeam("backend", "Backend Team")

	fromDB, err := s.repo.GetByID(s.T().Context(), team.ID)
	s.Require().NoError(err)
	s.Equal(team.ID, fromDB.ID)
	s.Equal(team.Code, fromDB.Code)
	s.Equal(team.Name, fromDB.Name)
}

func (s *TeamRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Team{}, fromDB)
}

func (s *TeamRepositorySuite) TestList() {
	t1 := s.createTeam("core", "Core Team")
	t2 := s.createTeam("platform", "Platform Team")

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)
	s.Contains(list, *t1)
	s.Contains(list, *t2)
}

func (s *TeamRepositorySuite) TestUpdate() {
	team := s.createTeam("ops", "Operations")

	team.Code = "devops"
	team.Name = "DevOps Team"

	err := s.repo.Update(s.T().Context(), team)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), team.ID)
	s.Require().NoError(err)
	s.Equal("devops", fromDB.Code)
	s.Equal("DevOps Team", fromDB.Name)
}

func (s *TeamRepositorySuite) TestDelete() {
	team := s.createTeam("old", "Old Team")

	err := s.repo.Delete(s.T().Context(), team.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), team.ID)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Team{}, fromDB)
}

func (s *TeamRepositorySuite) TestDeleteNotFound() {
	err := s.repo.Delete(s.T().Context(), 123)
	s.Require().NoError(err)
}

func (s *TeamRepositorySuite) createTeam(code, name string) *Team {
	ctx := s.T().Context()

	team := &Team{
		Code: code,
		Name: name,
	}

	err := s.repo.Create(ctx, team)
	s.Require().NoError(err)

	return team
}
