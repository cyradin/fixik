package container

import (
	"github.com/cyradin/fixik/internal/priority"
)

func (c *Container) PriorityManager() *priority.PriorityManager {
	if c.priorityManager == nil {
		c.priorityManager = priority.NewPriorityManager(
			c.PriorityRepository(),
			c.IncidentRepository(),
			c.TxExecutor(),
		)
	}

	return c.priorityManager
}
