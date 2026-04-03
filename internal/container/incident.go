package container

import "github.com/cyradin/fixik/internal/incident"

func (c *Container) IncidentManager() *incident.IncidentManager {
	if c.incidentManager == nil {
		c.incidentManager = incident.NewIncidentManager(
			c.IncidentRepository(),
			c.StatusManager(),
			c.PriorityManager(),
			c.TeamManager(),
			c.UserManager(),
		)
	}

	return c.incidentManager
}

func (c *Container) CommentManager() *incident.CommentManager {
	if c.commentManager == nil {
		c.commentManager = incident.NewCommentManager(
			c.CommentRepository(),
			c.UserManager(),
		)
	}

	return c.commentManager
}
