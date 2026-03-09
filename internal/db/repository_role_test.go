package db

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type RoleRepositorySuite struct {
	tests.PostgresSuite
	repo *RoleRepository
}

func TestRoleRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RoleRepositorySuite))
}

func (s *RoleRepositorySuite) SetupTest() {
	s.repo = NewRoleRepository(s.Postgres())

	_, err := s.Postgres().Exec(s.T().Context(), `
		TRUNCATE TABLE roles RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *RoleRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	role := &Role{
		Code: "admin",
		Name: "Administrator",
	}

	err := s.repo.Create(ctx, role)
	s.Require().NoError(err)
	s.NotZero(role.ID)
}

func (s *RoleRepositorySuite) TestGetByID_Found() {
	role := s.createRole("viewer", "Viewer")

	fromDB, err := s.repo.GetByID(s.T().Context(), role.ID)
	s.Require().NoError(err)
	s.Equal(role.ID, fromDB.ID)
	s.Equal(role.Code, fromDB.Code)
	s.Equal(role.Name, fromDB.Name)
}

func (s *RoleRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Role{}, fromDB)
}

func (s *RoleRepositorySuite) TestList() {
	r1 := s.createRole("admin", "Administrator")
	r2 := s.createRole("viewer", "Viewer")

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)
	s.Contains(list, *r1)
	s.Contains(list, *r2)
}

func (s *RoleRepositorySuite) TestUpdate() {
	role := s.createRole("ops", "Operations")

	role.Code = "devops"
	role.Name = "DevOps"

	err := s.repo.Update(s.T().Context(), role)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), role.ID)
	s.Require().NoError(err)
	s.Equal("devops", fromDB.Code)
	s.Equal("DevOps", fromDB.Name)
}

func (s *RoleRepositorySuite) TestDelete() {
	role := s.createRole("old_role", "Old Role")

	err := s.repo.Delete(s.T().Context(), role.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), role.ID)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Role{}, fromDB)
}

func (s *RoleRepositorySuite) TestDeleteNotFound() {
	err := s.repo.Delete(s.T().Context(), 123)
	s.Require().NoError(err)
}

func (s *RoleRepositorySuite) createRole(code, name string) *Role {
	ctx := s.T().Context()

	role := &Role{
		Code: code,
		Name: name,
	}

	err := s.repo.Create(ctx, role)
	s.Require().NoError(err)

	return role
}
