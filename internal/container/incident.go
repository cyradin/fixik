package container

import "github.com/cyradin/fixik/internal/incident"


func (c *Container) IncidentManager() *incident.IncidentManager {
	if c.incidentManager == nil {
		c.incidentManager = incident.NewIncidentManager(
			c.IncidentRepository(),
			c.ImpactManager(),
			c.StatusManager(),
			c.PriorityManager(),
		)
	}

	return c.incidentManager
}
