package server

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func aboutMe(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	return c.JSON(user)
}
