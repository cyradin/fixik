package incident

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/cyradin/fixik/pkg/tests"
)

type IncidentRepositorySuite struct {
	tests.PostgresSuite
	repo *Repository
}

func TestIncidentRepositorySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(IncidentRepositorySuite))
}

func (s *IncidentRepositorySuite) SetupTest() {
	s.repo = NewRepository(s.Postgres())

	_, err := s.Postgres().Exec(s.T().Context(), `
		TRUNCATE TABLE incidents RESTART IDENTITY CASCADE;
		TRUNCATE TABLE incident_statuses RESTART IDENTITY CASCADE;
		TRUNCATE TABLE incident_priorities RESTART IDENTITY CASCADE;
		TRUNCATE TABLE incident_impacts RESTART IDENTITY CASCADE;
	`)
	s.Require().NoError(err)
}

func (s *IncidentRepositorySuite) TestCreate() {
	ctx := s.T().Context()

	impact := s.createImpact("high", "High")
	priority := s.createPriority("critical", "Critical")
	status := s.createStatus("open", "Open")

	inc := &Incident{
		Title:       "DB down",
		Description: "Database unavailable",
		ImpactID:    impact.ID,
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
	impact := s.createImpact("medium", "Medium")
	priority := s.createPriority("medium", "Medium")
	status := s.createStatus("open", "Open")

	inc := s.loadIncidents([]Incident{
		{Title: "Service down", Description: "Manual insert", ImpactID: impact.ID, PriorityID: priority.ID, StatusID: status.ID},
	})[0]

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Equal(inc.ID, fromDB.ID)
	s.Equal(inc.Title, fromDB.Title)
	s.Equal(inc.Description, fromDB.Description)
	s.Equal(inc.ImpactID, fromDB.ImpactID)
	s.Equal(inc.PriorityID, fromDB.PriorityID)
	s.Equal(inc.StatusID, fromDB.StatusID)
}

func (s *IncidentRepositorySuite) TestUpdate_Found() {
	impact := s.createImpact("medium", "Medium")
	priority := s.createPriority("medium", "Medium")
	status := s.createStatus("open", "Open")
	newStatus := s.createStatus("in_progress", "In Progress")

	inc := s.loadIncidents([]Incident{
		{Title: "Service unavailable", Description: "Initial description", ImpactID: impact.ID, PriorityID: priority.ID, StatusID: status.ID},
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
	impact := s.createImpact("low", "Low")
	priority := s.createPriority("low", "Low")
	status := s.createStatus("closed", "Closed")

	inc := &Incident{
		ID:         999999,
		Title:      "Nonexistent",
		ImpactID:   impact.ID,
		PriorityID: priority.ID,
		StatusID:   status.ID,
	}
	err := s.repo.Update(s.T().Context(), inc)
	s.Require().Error(err)
}

func (s *IncidentRepositorySuite) TestDelete_Found() {
	impact := s.createImpact("low", "Low")
	priority := s.createPriority("low", "Low")
	status := s.createStatus("open", "Open")

	inc := s.loadIncidents([]Incident{
		{Title: "Temp incident", Description: "To be deleted", ImpactID: impact.ID, PriorityID: priority.ID, StatusID: status.ID},
	})[0]

	err := s.repo.Delete(s.T().Context(), inc.ID)
	s.Require().NoError(err)

	fromDB, err := s.repo.GetByID(s.T().Context(), inc.ID)
	s.Require().NoError(err)
	s.Require().Empty(fromDB)
}

func (s *IncidentRepositorySuite) TestDelete_NotFound() {
	err := s.repo.Delete(s.T().Context(), 999999)
	s.Require().NoError(err)
}

func (s *IncidentRepositorySuite) loadIncidents(fixtures []Incident) []Incident {
	s.T().Helper()
	ctx := s.T().Context()

	inserted := make([]Incident, len(fixtures))
	for i, inc := range fixtures {
		row := s.Postgres().QueryRow(ctx,
			`INSERT INTO incidents (title, description, impact_id, priority_id, status_id)
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id, created_at, updated_at`,
			inc.Title, inc.Description, inc.ImpactID, inc.PriorityID, inc.StatusID,
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

func (s *IncidentRepositorySuite) createImpact(code, name string) *Impact {
	ctx := s.T().Context()
	ir := NewImpactRepository(s.Postgres())

	impact := &Impact{
		Code: code,
		Name: name,
	}
	err := ir.Create(ctx, impact)
	s.Require().NoError(err)
	s.NotZero(impact.ID)

	return impact
}
