package container

import (
	"github.com/cyradin/fixik/internal/user"
)

func (c *Container) UserManager() *user.UserManager {
	if c.userManager == nil {
		c.userManager = user.NewUserManager(
			c.UserRepository(),
			c.IncidentRepository(),
			c.TxExecutor(),
		)
	}

	return c.userManager
}

func (c *Container) JWTManager() *user.JWTManager {
	if c.jwtManager == nil {
		c.jwtManager = user.NewJWTManager(
			c.cfg.Auth.Secret,
			c.cfg.Auth.AccessTokenTTL,
			c.cfg.Auth.RefreshTokenTTL,
		)
	}

	return c.jwtManager
}

func (c *Container) AuthService() *user.AuthService {
	if c.authService == nil {
		c.authService = user.NewAuthService(
			c.UserRepository(),
			c.UserRepository(),
			c.JWTManager(),
		)
	}

	return c.authService
}
