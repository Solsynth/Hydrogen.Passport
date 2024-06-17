package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/gofiber/fiber/v2"
)

func getOtherUserinfo(c *fiber.Ctx) error {
	alias := c.Params("alias")

	var account models.Account
	if err := database.C.
		Where(&models.Account{Name: alias}).
		Omit("tickets", "challenges", "factors", "events", "clients", "notifications", "notify_subscribers").
		Preload("Profile").
		Preload("Badges").
		First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(account)
}
