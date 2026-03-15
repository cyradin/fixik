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
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
		TRUNCATE TABLE teams RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *UserRepositorySuite) TestGetByID() {
	team := s.createTeam("team1", "Team One")

	user := s.createUser("Алексей", "alex", "alex@example.com", "passhash", team.ID, RoleUser)

	fromDB, err := s.repo.GetByID(s.T().Context(), user.ID)
	s.Require().NoError(err)

	s.Equal(user.ID, fromDB.ID)
	s.Equal(user.Name, fromDB.Name)
	s.Equal(user.Username, fromDB.Username)
	s.Equal(user.Email, fromDB.Email)
	s.Equal(user.Password, fromDB.Password)
	s.Equal(user.TeamID, fromDB.TeamID)
	s.Equal(user.Role, fromDB.Role)
	s.WithinDuration(time.Now(), fromDB.CreatedAt, time.Second)
	s.WithinDuration(time.Now(), fromDB.UpdatedAt, time.Second)
	s.Nil(fromDB.DeletedAt)
}

func (s *UserRepositorySuite) TestList() {
	ctx := s.T().Context()

	team := s.createTeam("team1", "Team One")

	u1 := s.createUser("Александр", "u1", "u1@example.com", "pass", team.ID, RoleAdmin)
	u2 := s.createUser("Мария", "u2", "u2@example.com", "pass", team.ID, RoleUser)
	u3 := s.createUser("Иван", "u3", "u3@example.com", "pass", team.ID, RoleManager)

	s.Run("list all", func() {
		users, err := s.repo.List(ctx, 100, 0)
		s.Require().NoError(err)

		s.Len(users, 3)
		s.Equal(u1.ID, users[0].ID)
		s.Equal(u2.ID, users[1].ID)
		s.Equal(u3.ID, users[2].ID)
	})

	s.Run("limit", func() {
		users, err := s.repo.List(ctx, 2, 0)
		s.Require().NoError(err)

		s.Len(users, 2)
		s.Equal(u1.ID, users[0].ID)
		s.Equal(u2.ID, users[1].ID)
	})

	s.Run("offset", func() {
		users, err := s.repo.List(ctx, 100, 1)
		s.Require().NoError(err)

		s.Len(users, 2)
		s.Equal(u2.ID, users[0].ID)
		s.Equal(u3.ID, users[1].ID)
	})

	s.Run("limit with offset", func() {
		users, err := s.repo.List(ctx, 1, 1)
		s.Require().NoError(err)

		s.Len(users, 1)
		s.Equal(u2.ID, users[0].ID)
	})

	s.Run("deleted users not returned", func() {
		err := s.repo.Delete(ctx, u2.ID)
		s.Require().NoError(err)

		users, err := s.repo.List(ctx, 100, 0)
		s.Require().NoError(err)

		s.Len(users, 2)
		s.Equal(u1.ID, users[0].ID)
		s.Equal(u3.ID, users[1].ID)
	})
}

func (s *UserRepositorySuite) TestUpdate() {
	team1 := s.createTeam("team1", "Team One")
	team2 := s.createTeam("team2", "Team Two")

	user := s.createUser("Боб", "bob", "bob@example.com", "pass", team1.ID, RoleUser)

	user.Name = "Боб2"
	user.Username = "bob2"
	user.Email = "bob2@example.com"
	user.Password = "newpass"
	user.TeamID = team2.ID
	user.Role = RoleAdmin

	oldUpdated := user.UpdatedAt

	err := s.repo.Update(s.T().Context(), user)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), user.ID)
	s.Require().NoError(err)
	s.Equal("Боб2", fromDB.Name)
	s.Equal("bob2", fromDB.Username)
	s.Equal("bob2@example.com", fromDB.Email)
	s.Equal("newpass", fromDB.Password)
	s.Equal(team2.ID, fromDB.TeamID)
	s.Equal(user.Role, fromDB.Role)
	s.True(fromDB.UpdatedAt.After(oldUpdated))
}

func (s *UserRepositorySuite) TestDelete_SoftDelete() {
	team := s.createTeam("team1", "Team One")
	user := s.createUser("Каролина", "carol", "carol@example.com", "pass", team.ID, RoleUser)

	err := s.repo.Delete(s.T().Context(), user.ID)
	s.Require().NoError(err)

	_, err = s.repo.GetByID(s.T().Context(), user.ID)
	s.Require().ErrorIs(err, ErrNotFound)
}

func (s *UserRepositorySuite) TestUpdate_NotFound() {
	team := s.createTeam("team1", "Team One")
	user := &User{
		ID:       999999,
		Name:     "Призрак",
		Username: "ghost",
		Email:    "ghost@example.com",
		Password: "pass",
		TeamID:   team.ID,
		Role:     RoleUser,
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

func (s *UserRepositorySuite) createUser(name, username, email, password string, teamID int64, role Role) *User {
	ctx := s.T().Context()
	user := &User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
		TeamID:   teamID,
		Role:     role,
	}

	err := s.repo.Create(ctx, user)
	s.Require().NoError(err)
	s.NotZero(user.ID)

	return user
}
