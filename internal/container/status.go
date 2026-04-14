package container

import (
	"github.com/cyradin/fixik/internal/status"
)

func (c *Container) StatusManager() *status.StatusManager {
	if c.statusManager == nil {
		c.statusManager = status.NewStatusManager(
			c.StatusRepository(),
			c.IncidentRepository(),
			c.TxExecutor(),
		)
	}

	return c.statusManager
}
