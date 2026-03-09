package db

import (
	"testing"

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
}

func (s *UserRepositorySuite) TestGetByID_NotFound() {
	_, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *UserRepositorySuite) TestList() {
	team1 := s.createTeam("team1", "Team One")
	team2 := s.createTeam("team2", "Team Two")
	role := s.createRole("user", "User")

	u1 := s.createUser("alice", "alice@example.com", "pass1", team1.ID, []int64{role.ID})
	u2 := s.createUser("bob", "bob@example.com", "pass2", team2.ID, []int64{})

	list, err := s.repo.List(s.T().Context())
	s.Require().NoError(err)
	s.Len(list, 2)

	ids := make([]int64, 0, 2)
	for _, u := range list {
		ids = append(ids, u.ID)
	}

	s.Contains(ids, u1.ID)
	s.Contains(ids, u2.ID)
}

func (s *UserRepositorySuite) TestUpdate() {
	team := s.createTeam("team1", "Team One")
	team2 := s.createTeam("team2", "Team Two")
	role := s.createRole("editor", "Editor")
	user := s.createUser("charlie", "charlie@example.com", "pass", team.ID, []int64{})

	user.Username = "charlie2"
	user.Email = "charlie2@example.com"
	user.Password = "newpass"
	user.TeamID = team2.ID
	user.RoleIDs = []int64{role.ID}

	err := s.repo.Update(s.T().Context(), user)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), user.ID)
	s.Require().NoError(err)
	s.Equal("charlie2", fromDB.Username)
	s.Equal("charlie2@example.com", fromDB.Email)
	s.Equal("newpass", fromDB.Password)
	s.Equal(team2.ID, fromDB.TeamID)
	s.ElementsMatch([]int64{role.ID}, fromDB.RoleIDs)
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

func (s *UserRepositorySuite) TestDelete() {
	team := s.createTeam("team1", "Team One")
	role := s.createRole("guest", "Guest")
	user := s.createUser("dave", "dave@example.com", "pass", team.ID, []int64{role.ID})

	err := s.repo.Delete(s.T().Context(), user.ID)
	s.Require().NoError(err)

	_, err = s.repo.GetByID(s.T().Context(), user.ID)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *UserRepositorySuite) TestDelete_NotFound() {
	err := s.repo.Delete(s.T().Context(), 999999)
	s.Require().NoError(err)
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
