package utils

import "github.com/gofiber/fiber/v2"

type AuthFunc func(c *fiber.Ctx, overrides ...string) error
