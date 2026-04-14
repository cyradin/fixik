package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type IncidentRepositorySuite struct {
	tests.PostgresSuite
	repo *IncidentRepository
}

func TestIncidentRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IncidentRepositorySuite))
}

func (s *IncidentRepositorySuite) SetupTest() {
	s.repo = NewIncidentRepository(s.Postgres())

	_, err := s.Postgres().Exec(s.T().Context(), `
		TRUNCATE TABLE incidents RESTART IDENTITY CASCADE;
		TRUNCATE TABLE statuses RESTART IDENTITY CASCADE;
		TRUNCATE TABLE priorities RESTART IDENTITY CASCADE;
		TRUNCATE TABLE teams RESTART IDENTITY CASCADE;
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *IncidentRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	priority := s.createPriority("critical", "Critical")
	status := s.createStatus("open", "Open")
	team := s.createTeam("dev", "Dev Team")
	user := s.createUser("Alice", "alice", "alice@test.com", team.ID)
	author := s.createUser("Bob", "bob", "bob@test.com", team.ID)

	inc := &Incident{
		Title:       "DB down",
		Description: "Database unavailable",
		PriorityID:  priority.ID,
		StatusID:    status.ID,
		TeamID:      &team.ID,
		UserID:      &user.ID,
		AuthorID:    &author.ID,
	}

	err := s.repo.Create(ctx, inc)
	s.Require().NoError(err)
	s.NotZero(inc.ID)
	s.NotZero(inc.CreatedAt)
	s.NotZero(inc.UpdatedAt)
	s.Equal(&team.ID, inc.TeamID)
	s.Equal(&user.ID, inc.UserID)
	s.Equal(&author.ID, inc.AuthorID)
}

func (s *IncidentRepositorySuite) TestGetByID_Found() {
	priority := s.createPriority("medium", "Medium")
	status := s.createStatus("open", "Open")
	team := s.createTeam("ops", "Ops Team")
	user := s.createUser("Alice", "alice", "alice@test.com", team.ID)
	author := s.createUser("Bob", "bob", "bob@test.com", team.ID)

	inc := s.loadIncidents([]Incident{
		{
			Title:       "Service down",
			Description: "Manual insert",
			PriorityID:  priority.ID,
			StatusID:    status.ID,
			TeamID:      &team.ID,
			UserID:      &user.ID,
			AuthorID:    &author.ID,
		},
	})[0]

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Equal(inc.ID, fromDB.ID)
	s.Equal(inc.TeamID, fromDB.TeamID)
	s.Equal(inc.UserID, fromDB.UserID)
	s.Equal(inc.AuthorID, fromDB.AuthorID)
}

func (s *IncidentRepositorySuite) TestGetByID_NotFound() {
	fromDB, err := s.repo.GetByID(s.T().Context(), 999999)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Equal(Incident{}, fromDB)
}

func (s *IncidentRepositorySuite) TestUpdate_Found() {
	priority := s.createPriority("medium", "Medium")
	status := s.createStatus("open", "Open")
	newStatus := s.createStatus("in_progress", "In Progress")
	team := s.createTeam("team1", "Team 1")
	user := s.createUser("Charlie", "charlie", "charlie@test.com", team.ID)
	author := s.createUser("Bob", "bob", "bob@test.com", team.ID)

	inc := s.loadIncidents([]Incident{
		{
			Title:       "Service unavailable",
			Description: "Initial description",
			PriorityID:  priority.ID,
			StatusID:    status.ID,
			TeamID:      &team.ID,
			UserID:      &user.ID,
		},
	})[0]

	oldUpdatedAt := inc.UpdatedAt
	inc.Title = "Updated title"
	inc.StatusID = newStatus.ID

	newTeam := s.createTeam("team2", "Team 2")
	inc.TeamID = &newTeam.ID
	inc.UserID = nil
	inc.AuthorID = &author.ID

	err := s.repo.Update(s.T().Context(), &inc)
	s.Require().NoError(err)
	s.True(inc.UpdatedAt.After(oldUpdatedAt))

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Equal("Updated title", fromDB.Title)
	s.Equal(newStatus.ID, fromDB.StatusID)
	s.Equal(&newTeam.ID, fromDB.TeamID)
	s.Equal(&author.ID, fromDB.AuthorID)
	s.Nil(fromDB.UserID)
}

func (s *IncidentRepositorySuite) TestUpdate_NotFound() {
	priority := s.createPriority("low", "Low")
	status := s.createStatus("closed", "Closed")

	inc := &Incident{
		ID:         999999,
		Title:      "Nonexistent",
		PriorityID: priority.ID,
		StatusID:   status.ID,
	}
	err := s.repo.Update(s.T().Context(), inc)
	s.Require().Error(err)
}

func (s *IncidentRepositorySuite) TestDelete_Found() {
	priority := s.createPriority("low", "Low")
	status := s.createStatus("open", "Open")
	team := s.createTeam("team3", "Team 3")
	user := s.createUser("Dave", "dave", "dave@test.com", team.ID)

	inc := s.loadIncidents([]Incident{
		{Title: "Temp incident", Description: "To be deleted", PriorityID: priority.ID, StatusID: status.ID, TeamID: &team.ID, UserID: &user.ID},
	})[0]

	err := s.repo.Delete(s.T().Context(), inc.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().ErrorIs(err, ErrNotFound)
	s.Require().Empty(fromDB)
}

func (s *IncidentRepositorySuite) TestDelete_NotFound() {
	err := s.repo.Delete(s.T().Context(), 999999)
	s.Require().NoError(err)
}

func (s *IncidentRepositorySuite) TestList() {
	ctx := s.T().Context()

	priority := s.createPriority("critical", "Critical")
	status := s.createStatus("open", "Open")

	incidents := s.loadIncidents([]Incident{
		{Title: "Incident 1", Description: "Desc 1", PriorityID: priority.ID, StatusID: status.ID},
		{Title: "Incident 2", Description: "Desc 2", PriorityID: priority.ID, StatusID: status.ID},
		{Title: "Incident 3", Description: "Desc 3", PriorityID: priority.ID, StatusID: status.ID},
	})

	s.Run("all", func() {
		res, err := s.repo.List(ctx, 10, 0)
		s.Require().NoError(err)

		s.Len(res.Items, 3)
		s.Equal(3, res.Total)

		s.Equal(incidents[2].Title, res.Items[0].Title)
		s.Equal(incidents[1].Title, res.Items[1].Title)
		s.Equal(incidents[0].Title, res.Items[2].Title)
	})

	s.Run("limit", func() {
		res, err := s.repo.List(ctx, 2, 0)
		s.Require().NoError(err)

		s.Len(res.Items, 2)
		s.Equal(3, res.Total)

		s.Equal(incidents[2].Title, res.Items[0].Title)
		s.Equal(incidents[1].Title, res.Items[1].Title)
	})

	s.Run("offset", func() {
		res, err := s.repo.List(ctx, 2, 1)
		s.Require().NoError(err)

		s.Len(res.Items, 2)
		s.Equal(3, res.Total)

		s.Equal(incidents[1].Title, res.Items[0].Title)
		s.Equal(incidents[0].Title, res.Items[1].Title)
	})

	s.Run("deleted filtered", func() {
		err := s.repo.Delete(ctx, incidents[1].ID)
		s.Require().NoError(err)

		res, err := s.repo.List(ctx, 10, 0)
		s.Require().NoError(err)

		s.Len(res.Items, 2)
		s.Equal(2, res.Total)

		for _, i := range res.Items {
			s.NotEqual(incidents[1].ID, i.ID)
		}
	})
}

func (s *IncidentRepositorySuite) loadIncidents(fixtures []Incident) []Incident {
	s.T().Helper()
	ctx := s.T().Context()

	inserted := make([]Incident, len(fixtures))
	for i, inc := range fixtures {
		if inc.TeamID == nil {
			team := s.createTeam("team"+strconv.Itoa(i), "Team "+strconv.Itoa(i))
			inc.TeamID = &team.ID
		}

		if inc.UserID == nil {
			user := s.createUser("User"+strconv.Itoa(i), "user"+strconv.Itoa(i), fmt.Sprintf("user%d@test.com", i), *inc.TeamID)
			inc.UserID = &user.ID
		}

		row := s.Postgres().QueryRow(ctx,
			`INSERT INTO incidents (title, description, priority_id, status_id, team_id, user_id, author_id)
             VALUES ($1, $2, $3, $4, $5, $6, $7)
             RETURNING id, created_at, updated_at`,
			inc.Title, inc.Description, inc.PriorityID, inc.StatusID, inc.TeamID, inc.UserID, inc.AuthorID,
		)
		err := row.Scan(&inc.ID, &inc.CreatedAt, &inc.UpdatedAt)
		s.Require().NoError(err)

		inserted[i] = inc
	}

	return inserted
}

func (s *IncidentRepositorySuite) createStatus(code, name string) *Status {
	ctx := s.T().Context()
	sr := NewStatusRepository(s.Postgres())

	status := &Status{Code: code, Name: name}
	err := sr.Create(ctx, status)
	s.Require().NoError(err)
	s.NotZero(status.ID)

	return status
}

func (s *IncidentRepositorySuite) createPriority(code, name string) *Priority {
	ctx := s.T().Context()
	pr := NewPriorityRepository(s.Postgres())

	priority := &Priority{Code: code, Name: name}
	err := pr.Create(ctx, priority)
	s.Require().NoError(err)
	s.NotZero(priority.ID)

	return priority
}

func (s *IncidentRepositorySuite) createTeam(code, name string) *Team {
	ctx := s.T().Context()
	tr := NewTeamRepository(s.Postgres())

	team := &Team{Code: code, Name: name}
	err := tr.Create(ctx, team)
	s.Require().NoError(err)
	s.NotZero(team.ID)

	return team
}

func (s *IncidentRepositorySuite) createUser(name, username, email string, teamID int64) *User {
	ctx := s.T().Context()
	ur := NewUserRepository(s.Postgres())

	user := &User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: "password123",
		TeamID:   &teamID,
		Role:     RoleUser,
	}
	err := ur.Create(ctx, user)
	s.Require().NoError(err)
	s.NotZero(user.ID)

	return user
}

func (s *IncidentRepositorySuite) TestCountByStatus() {
	ctx := s.T().Context()

	priority := s.createPriority("p1", "P1")
	status1 := s.createStatus("open", "Open")
	status2 := s.createStatus("closed", "Closed")

	incidents := s.loadIncidents([]Incident{
		{Title: "i1", Description: "d1", PriorityID: priority.ID, StatusID: status1.ID},
		{Title: "i2", Description: "d2", PriorityID: priority.ID, StatusID: status1.ID},
		{Title: "i3", Description: "d3", PriorityID: priority.ID, StatusID: status2.ID},
	})

	s.Run("count correct", func() {
		count, err := s.repo.CountByStatus(ctx, status1.ID)
		s.Require().NoError(err)
		s.Equal(2, count)
	})

	s.Run("deleted not counted", func() {
		err := s.repo.Delete(ctx, incidents[0].ID)
		s.Require().NoError(err)

		count, err := s.repo.CountByStatus(ctx, status1.ID)
		s.Require().NoError(err)
		s.Equal(1, count)
	})

	s.Run("zero", func() {
		count, err := s.repo.CountByStatus(ctx, 999999)
		s.Require().NoError(err)
		s.Equal(0, count)
	})
}

func (s *IncidentRepositorySuite) TestCountByPriority() {
	ctx := s.T().Context()

	priority1 := s.createPriority("p1", "P1")
	priority2 := s.createPriority("p2", "P2")
	status := s.createStatus("open", "Open")

	s.loadIncidents([]Incident{
		{Title: "i1", Description: "d1", PriorityID: priority1.ID, StatusID: status.ID},
		{Title: "i2", Description: "d2", PriorityID: priority1.ID, StatusID: status.ID},
		{Title: "i3", Description: "d3", PriorityID: priority2.ID, StatusID: status.ID},
	})

	count, err := s.repo.CountByPriority(ctx, priority1.ID)
	s.Require().NoError(err)
	s.Equal(2, count)
}

func (s *IncidentRepositorySuite) TestCountByTeam() {
	ctx := s.T().Context()

	priority := s.createPriority("p1", "P1")
	status := s.createStatus("open", "Open")

	team1 := s.createTeam("t1", "Team 1")
	team2 := s.createTeam("t2", "Team 2")

	user1 := s.createUser("u1", "u1", "u1@test.com", team1.ID)
	user2 := s.createUser("u2", "u2", "u2@test.com", team2.ID)

	s.loadIncidents([]Incident{
		{Title: "i1", Description: "d1", PriorityID: priority.ID, StatusID: status.ID, TeamID: &team1.ID, UserID: &user1.ID},
		{Title: "i2", Description: "d2", PriorityID: priority.ID, StatusID: status.ID, TeamID: &team1.ID, UserID: &user1.ID},
		{Title: "i3", Description: "d3", PriorityID: priority.ID, StatusID: status.ID, TeamID: &team2.ID, UserID: &user2.ID},
	})

	count, err := s.repo.CountByTeam(ctx, team1.ID)
	s.Require().NoError(err)
	s.Equal(2, count)
}

func (s *IncidentRepositorySuite) TestCountByUser() {
	ctx := s.T().Context()

	priority := s.createPriority("p1", "P1")
	status := s.createStatus("open", "Open")
	team := s.createTeam("t1", "Team 1")

	user1 := s.createUser("u1", "u1", "u1@test.com", team.ID)
	user2 := s.createUser("u2", "u2", "u2@test.com", team.ID)

	s.loadIncidents([]Incident{
		{
			Title:       "i1",
			Description: "d1",
			PriorityID:  priority.ID,
			StatusID:    status.ID,
			UserID:      &user1.ID,
		},
		{
			Title:       "i2",
			Description: "d2",
			PriorityID:  priority.ID,
			StatusID:    status.ID,
			UserID:      &user1.ID,
		},

		{
			Title:       "i3",
			Description: "d3",
			PriorityID:  priority.ID,
			StatusID:    status.ID,
			UserID:      &user2.ID,
		},

		{
			Title:       "i4",
			Description: "d4",
			PriorityID:  priority.ID,
			StatusID:    status.ID,
			AuthorID:    &user1.ID,
		},
	})

	count, err := s.repo.CountByUser(ctx, user1.ID)
	s.Require().NoError(err)

	s.Equal(3, count)
}
