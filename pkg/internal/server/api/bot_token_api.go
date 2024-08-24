package api

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func listBotKeys(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var keys []models.ApiKey
	if err := database.C.Where("account_id = ?", user.ID).Find(&keys).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(keys)
}

func getBotKey(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	id, _ := c.ParamsInt("id", 0)

	var key models.ApiKey
	if err := database.C.Where("id = ? AND account_id = ?", id, user.ID).First(&key).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(key)
}

func createBotKey(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Lifecycle   *int64   `json:"lifecycle"`
		Claims      []string `json:"claims"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	key, err := services.NewApiKey(user, models.ApiKey{
		Name:        data.Name,
		Description: data.Description,
		Lifecycle:   data.Lifecycle,
	}, c.IP(), c.Get(fiber.HeaderUserAgent), data.Claims)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(key)
}

func editBotKey(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Lifecycle   *int64 `json:"lifecycle"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	id, _ := c.ParamsInt("id", 0)

	var key models.ApiKey
	if err := database.C.Where("id = ? AND account_id = ?", id, user.ID).First(&key).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	key.Name = data.Name
	key.Description = data.Description
	key.Lifecycle = data.Lifecycle

	if err := database.C.Save(&key).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(key)
}

func rollBotKey(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	id, _ := c.ParamsInt("id", 0)

	var key models.ApiKey
	if err := database.C.Where("id = ? AND account_id = ?", id, user.ID).First(&key).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if key, err := services.RollApiKey(key); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(key)
	}
}

func revokeBotKey(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	id, _ := c.ParamsInt("id", 0)

	var key models.ApiKey
	if err := database.C.Where("id = ? AND account_id = ?", id, user.ID).First(&key).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := database.C.Delete(&key).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(key)
}
