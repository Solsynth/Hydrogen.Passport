package server

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"code.smartsheep.studio/hydrogen/passport/pkg/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func getPrincipal(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data models.Account
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: user.ID},
	}).Preload(clause.Associations).First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(data)
}

func doRegister(c *fiber.Ctx) error {
	var data struct {
		Name     string `json:"name"`
		Nick     string `json:"nick"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	if user, err := services.CreateAccount(
		data.Name,
		data.Nick,
		data.Email,
		data.Password,
	); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(user)
	}
}

func doRegisterConfirm(c *fiber.Ctx) error {
	var data struct {
		Code string `json:"code"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := services.ConfirmAccount(data.Code); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
