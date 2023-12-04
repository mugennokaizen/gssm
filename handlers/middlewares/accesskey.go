package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"gssm/db"
	"gssm/utils"
	"slices"
)

var accessKeyDefaultInjectorConfig = InjectorConfig{
	Filter:   nil,
	Injector: nil,
}

func accessKeyInjectorConfigDefault(config ...InjectorConfig) InjectorConfig {
	if len(config) < 1 {
		return accessKeyDefaultInjectorConfig
	}

	cfg := config[0]

	if cfg.Injector == nil {
		panic("Injector can't be nil")
	}

	if cfg.Filter == nil {
		cfg.Filter = accessKeyDefaultInjectorConfig.Filter
	}

	return cfg
}

func NewAccessKey(config InjectorConfig) fiber.Handler {
	cfg := accessKeyInjectorConfigDefault(config)

	aks := do.MustInvoke[*db.AccessKeySource](cfg.Injector)

	return func(c *fiber.Ctx) error {

		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}

		ctx := c.UserContext()

		key := c.Get("GSSM-Access-Key", "")

		if key == "" {
			return fiber.ErrUnauthorized
		}

		secret := c.Get("GSSM-Secret-Key", "")

		if secret == "" {
			return fiber.ErrUnauthorized
		}

		accessKey, err := aks.Get(ctx, key)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		signature, err := utils.GetSignature(accessKey.UserId, accessKey.ProjectId, secret)
		if err != nil {
			return fiber.ErrNotFound
		}

		if !slices.Equal(accessKey.Signature, signature) {
			return fiber.ErrNotFound
		}

		c.Locals("userId", accessKey.UserId)
		c.Locals("projectId", accessKey.ProjectId)

		return c.Next()
	}
}
