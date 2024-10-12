package api

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func getAuthConfig(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	return c.JSON(user.AuthConfig)
}

func updateAuthConfig(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data models.AuthConfig
	if err := exts.BindAndValidate(c, &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user.AuthConfig = datatypes.NewJSONType(data)

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.InvalidAuthCacheWithUser(user.ID)
	}

	return c.JSON(user.AuthConfig)
}

func getNotificationPreference(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	notification, err := services.GetNotificationPreference(user)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(notification)
}

func updateNotificationPreference(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Config map[string]bool `json:"config"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	notification, err := services.UpdateNotificationPreference(user, data.Config)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(notification)
}
