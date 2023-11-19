package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Filter func(c *fiber.Ctx) bool // Required
}
