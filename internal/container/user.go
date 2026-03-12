package container

import (
	"github.com/cyradin/fixik/internal/user"
)

func (c *Container) UserManager() *user.UserManager {
	if c.userManager == nil {
		c.userManager = user.NewUserManager(
			c.UserRepository(),
		)
	}

	return c.userManager
}
