package db

import (
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
		TRUNCATE TABLE incident_statuses RESTART IDENTITY CASCADE;
		TRUNCATE TABLE incident_priorities RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *IncidentRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	priority := s.createPriority("critical", "Critical")
	status := s.createStatus("open", "Open")

	inc := &Incident{
		Title:       "DB down",
		Description: "Database unavailable",
		PriorityID:  priority.ID,
		StatusID:    status.ID,
	}

	err := s.repo.Create(ctx, inc)
	s.Require().NoError(err)
	s.NotZero(inc.ID)
	s.NotZero(inc.CreatedAt)
	s.NotZero(inc.UpdatedAt)
}

func (s *IncidentRepositorySuite) TestGetByID_Found() {
	priority := s.createPriority("medium", "Medium")
	status := s.createStatus("open", "Open")

	inc := s.loadIncidents([]Incident{
		{Title: "Service down", Description: "Manual insert", PriorityID: priority.ID, StatusID: status.ID},
	})[0]

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Equal(inc.ID, fromDB.ID)
	s.Equal(inc.Title, fromDB.Title)
	s.Equal(inc.Description, fromDB.Description)
	s.Equal(inc.PriorityID, fromDB.PriorityID)
	s.Equal(inc.StatusID, fromDB.StatusID)
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

	inc := s.loadIncidents([]Incident{
		{Title: "Service unavailable", Description: "Initial description", PriorityID: priority.ID, StatusID: status.ID},
	})[0]

	oldUpdatedAt := inc.UpdatedAt
	inc.Title = "Updated title"
	inc.StatusID = newStatus.ID

	err := s.repo.Update(s.T().Context(), &inc)
	s.Require().NoError(err)
	s.True(inc.UpdatedAt.After(oldUpdatedAt))

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Equal("Updated title", fromDB.Title)
	s.Equal(newStatus.ID, fromDB.StatusID)
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

	inc := s.loadIncidents([]Incident{
		{Title: "Temp incident", Description: "To be deleted", PriorityID: priority.ID, StatusID: status.ID},
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
		list, err := s.repo.List(ctx, 10, 0)
		s.Require().NoError(err)

		s.Len(list, 3)
		s.Equal(incidents[2].Title, list[0].Title)
		s.Equal(incidents[1].Title, list[1].Title)
		s.Equal(incidents[0].Title, list[2].Title)
	})

	s.Run("limit", func() {
		list, err := s.repo.List(ctx, 2, 0)
		s.Require().NoError(err)

		s.Len(list, 2)
		s.Equal(incidents[2].Title, list[0].Title)
		s.Equal(incidents[1].Title, list[1].Title)
	})

	s.Run("offset", func() {
		list, err := s.repo.List(ctx, 2, 1)
		s.Require().NoError(err)

		s.Len(list, 2)
		s.Equal(incidents[1].Title, list[0].Title)
		s.Equal(incidents[0].Title, list[1].Title)
	})

	s.Run("deleted filtered", func() {
		err := s.repo.Delete(ctx, incidents[1].ID)
		s.Require().NoError(err)

		list, err := s.repo.List(ctx, 10, 0)
		s.Require().NoError(err)

		s.Len(list, 2)

		for _, i := range list {
			s.NotEqual(incidents[1].ID, i.ID)
		}
	})
}

func (s *IncidentRepositorySuite) loadIncidents(fixtures []Incident) []Incident {
	s.T().Helper()
	ctx := s.T().Context()

	inserted := make([]Incident, len(fixtures))
	for i, inc := range fixtures {
		row := s.Postgres().QueryRow(ctx,
			`INSERT INTO incidents (title, description, priority_id, status_id)
             VALUES ($1, $2, $3, $4)
             RETURNING id, created_at, updated_at`,
			inc.Title, inc.Description, inc.PriorityID, inc.StatusID,
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

	status := &Status{
		Code: code,
		Name: name,
	}
	err := sr.Create(ctx, status)
	s.Require().NoError(err)
	s.NotZero(status.ID)

	return status
}

func (s *IncidentRepositorySuite) createPriority(code, name string) *Priority {
	ctx := s.T().Context()
	pr := NewPriorityRepository(s.Postgres())

	priority := &Priority{
		Code: code,
		Name: name,
	}
	err := pr.Create(ctx, priority)
	s.Require().NoError(err)
	s.NotZero(priority.ID)

	return priority
}
