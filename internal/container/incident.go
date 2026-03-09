package container

import "github.com/cyradin/fixik/internal/incident"

func (c *Container) StatusManager() *incident.StatusManager {
	if c.statusManager == nil {
		c.statusManager = incident.NewStatusManager(
			c.StatusRepository(),
		)
	}

	return c.statusManager
}

func (c *Container) ImpactManager() *incident.ImpactManager {
	if c.impactManager == nil {
		c.impactManager = incident.NewImpactManager(
			c.ImpactRepository(),
		)
	}

	return c.impactManager
}

func (c *Container) PriorityManager() *incident.PriorityManager {
	if c.priorityManager == nil {
		c.priorityManager = incident.NewPriorityManager(
			c.PriorityRepository(),
		)
	}

	return c.priorityManager
}

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
