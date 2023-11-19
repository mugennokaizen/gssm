package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gssm/types"
	"time"
)

type JwtConfig struct {
	Filter            func(c *fiber.Ctx) bool
	RefreshCookieName string
	AccessCookieName  string
	SecretKey         string
}

var cookieDefaultConfig = JwtConfig{
	Filter: nil,
}

func cookieConfigDefault(config ...JwtConfig) JwtConfig {
	// Return default config if nothing provided
	if len(config) < 1 {
		return cookieDefaultConfig
	}

	// Override default config
	cfg := config[0]

	// Set default values if not passed
	if cfg.Filter == nil {
		cfg.Filter = cookieDefaultConfig.Filter
	}

	return cfg
}

func NewJwt(config JwtConfig) fiber.Handler {
	cfg := cookieConfigDefault(config)

	return func(c *fiber.Ctx) error {

		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		refreshToken := c.Cookies(cfg.RefreshCookieName)
		if len(refreshToken) == 0 {
			return fiber.ErrUnauthorized
		}

		accessToken := c.Cookies(cfg.AccessCookieName)

		if len(accessToken) == 0 {
			return fiber.ErrUnauthorized
		}

		var claims types.Claims
		jwtKey := []byte(cfg.SecretKey)

		tkn, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid || claims.ExpiresAt.Unix() < time.Now().Unix() {
			return fiber.ErrUnauthorized
		}
		c.Locals("claims", claims)

		return c.Next()
	}
}
