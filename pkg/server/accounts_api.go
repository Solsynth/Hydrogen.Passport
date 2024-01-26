package server

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"code.smartsheep.studio/hydrogen/passport/pkg/security"
	"github.com/gofiber/fiber/v2"
)

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

	user := models.Account{
		Name:  data.Name,
		Nick:  data.Nick,
		State: models.PendingAccountState,
		Factors: []models.AuthFactor{
			{
				Type:   models.PasswordAuthFactor,
				Secret: security.HashPassword(data.Password),
			},
		},
		Contacts: []models.AccountContact{
			{
				Type:       models.EmailAccountContact,
				Content:    data.Email,
				VerifiedAt: nil,
			},
		},
	}

	if err := database.C.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(user)
}
