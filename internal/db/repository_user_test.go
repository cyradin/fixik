package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type UserRepositorySuite struct {
	tests.PostgresSuite
	repo *UserRepository
}

func TestUserRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserRepositorySuite))
}

func (s *UserRepositorySuite) SetupTest() {
	s.repo = NewUserRepository(s.Postgres())

	ctx := s.T().Context()
	_, err := s.Postgres().Exec(ctx, `
		TRUNCATE TABLE user_roles RESTART IDENTITY CASCADE;
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
		TRUNCATE TABLE roles RESTART IDENTITY CASCADE;
		TRUNCATE TABLE teams RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *UserRepositorySuite) TestCreateAndGetByID() {
	team := s.createTeam("team1", "Team One")
	role1 := s.createRole("admin", "Admin")
	role2 := s.createRole("user", "User")

	user := s.createUser("alice", "alice@example.com", "passhash", team.ID, []int64{role1.ID, role2.ID})

	fromDB, err := s.repo.GetByID(s.T().Context(), user.ID)
	s.Require().NoError(err)

	s.Equal(user.ID, fromDB.ID)
	s.Equal(user.Username, fromDB.Username)
	s.Equal(user.Email, fromDB.Email)
	s.Equal(user.Password, fromDB.Password)
	s.Equal(user.TeamID, fromDB.TeamID)
	s.ElementsMatch(user.RoleIDs, fromDB.RoleIDs)
	s.WithinDuration(time.Now(), fromDB.CreatedAt, time.Second)
	s.WithinDuration(time.Now(), fromDB.UpdatedAt, time.Second)
	s.Nil(fromDB.DeletedAt)
}

func (s *UserRepositorySuite) TestUpdate() {
	team1 := s.createTeam("team1", "Team One")
	team2 := s.createTeam("team2", "Team Two")
	role := s.createRole("editor", "Editor")

	user := s.createUser("bob", "bob@example.com", "pass", team1.ID, []int64{})

	user.Username = "bob2"
	user.Email = "bob2@example.com"
	user.Password = "newpass"
	user.TeamID = team2.ID
	user.RoleIDs = []int64{role.ID}

	oldUpdated := user.UpdatedAt

	err := s.repo.Update(s.T().Context(), user)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), user.ID)
	s.Require().NoError(err)
	s.Equal("bob2", fromDB.Username)
	s.Equal("bob2@example.com", fromDB.Email)
	s.Equal("newpass", fromDB.Password)
	s.Equal(team2.ID, fromDB.TeamID)
	s.ElementsMatch([]int64{role.ID}, fromDB.RoleIDs)
	s.True(fromDB.UpdatedAt.After(oldUpdated))
}

func (s *UserRepositorySuite) TestDelete_SoftDelete() {
	team := s.createTeam("team1", "Team One")
	user := s.createUser("carol", "carol@example.com", "pass", team.ID, []int64{})

	err := s.repo.Delete(s.T().Context(), user.ID)
	s.Require().NoError(err)

	_, err = s.repo.GetByID(s.T().Context(), user.ID)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *UserRepositorySuite) TestUpdate_NotFound() {
	team := s.createTeam("team1", "Team One")
	user := &User{
		ID:       999999,
		Username: "ghost",
		Email:    "ghost@example.com",
		Password: "pass",
		TeamID:   team.ID,
		RoleIDs:  []int64{},
	}

	err := s.repo.Update(s.T().Context(), user)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *UserRepositorySuite) createTeam(code, name string) *Team {
	ctx := s.T().Context()
	team := &Team{Code: code, Name: name}

	err := s.Postgres().QueryRow(ctx,
		`INSERT INTO teams (code, name) VALUES ($1, $2) RETURNING id`,
		team.Code, team.Name,
	).Scan(&team.ID)
	s.Require().NoError(err)

	return team
}

func (s *UserRepositorySuite) createRole(code, name string) *Role {
	ctx := s.T().Context()
	role := &Role{Code: code, Name: name}

	err := s.Postgres().QueryRow(ctx,
		`INSERT INTO roles (code, name) VALUES ($1, $2) RETURNING id`,
		role.Code, role.Name,
	).Scan(&role.ID)
	s.Require().NoError(err)

	return role
}

func (s *UserRepositorySuite) createUser(username, email, password string, teamID int64, roleIDs []int64) *User {
	ctx := s.T().Context()
	user := &User{
		Username: username,
		Email:    email,
		Password: password,
		TeamID:   teamID,
		RoleIDs:  roleIDs,
	}

	err := s.repo.Create(ctx, user)
	s.Require().NoError(err)
	s.NotZero(user.ID)

	return user
}
