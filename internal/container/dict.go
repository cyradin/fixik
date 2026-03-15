package container

import (
	"github.com/cyradin/fixik/internal/dict"
)

func (c *Container) StatusManager() *dict.EntityManager {
	if c.statusManager == nil {
		c.statusManager = dict.NewStatusManager(
			c.StatusRepository(),
		)
	}

	return c.statusManager
}

func (c *Container) PriorityManager() *dict.EntityManager {
	if c.priorityManager == nil {
		c.priorityManager = dict.NewPriorityManager(
			c.PriorityRepository(),
		)
	}

	return c.priorityManager
}

func (c *Container) TeamManager() *dict.EntityManager {
	if c.teamManager == nil {
		c.teamManager = dict.NewTeamManager(
			c.TeamRepository(),
		)
	}

	return c.teamManager
}
