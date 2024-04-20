package ui

import "github.com/gofiber/fiber/v2"

func selfUserinfoPage(c *fiber.Ctx) error {
	return c.Render("views/users/me/index", fiber.Map{})
}
