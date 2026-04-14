package container

import "github.com/cyradin/fixik/internal/team"

func (c *Container) TeamManager() *team.TeamManager {
	if c.teamManager == nil {
		c.teamManager = team.NewTeamManager(
			c.TeamRepository(),
			c.IncidentRepository(),
			c.UserRepository(),
			c.TxExecutor(),
		)
	}

	return c.teamManager
}
