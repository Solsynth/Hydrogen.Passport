package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func getPersonalPage(c *fiber.Ctx) error {
	alias := c.Params("alias")

	var account models.Account
	if err := database.C.
		Where(&models.Account{Name: alias}).
		First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var page models.AccountPage
	if err := database.C.
		Where(&models.AccountPage{AccountID: account.ID}).
		First(&page).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(page)
}

func getOwnPersonalPage(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var page models.AccountPage
	if err := database.C.
		Where(&models.AccountPage{AccountID: user.ID}).
		FirstOrCreate(&page, &models.AccountPage{AccountID: user.ID}).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(page)
}

func editPersonalPage(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Content string                    `json:"content"`
		Links   []models.AccountPageLinks `json:"links"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	var page models.AccountPage
	if err := database.C.
		Where(&models.AccountPage{AccountID: user.ID}).
		FirstOrInit(&page).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	page.Content = data.Content
	page.Links = data.Links

	if err := database.C.Save(&page).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
