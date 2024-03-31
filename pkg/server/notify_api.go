package server

import (
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"git.solsynth.dev/hydrogen/identity/pkg/services"
	"github.com/gofiber/fiber/v2"
)

func notifyUser(c *fiber.Ctx) error {
	var data struct {
		ClientID     string                    `json:"client_id" validate:"required"`
		ClientSecret string                    `json:"client_secret" validate:"required"`
		Subject      string                    `json:"subject" validate:"required,max=1024"`
		Content      string                    `json:"content" validate:"required,max=3072"`
		Links        []models.NotificationLink `json:"links"`
		IsImportant  bool                      `json:"is_important"`
		IsRealtime   bool                      `json:"is_realtime"`
		UserID       uint                      `json:"user_id" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	client, err := services.GetThirdClientWithSecret(data.ClientID, data.ClientSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	var user models.Account
	if user, err = services.GetAccount(data.UserID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	notification := models.Notification{
		Subject:     data.Subject,
		Content:     data.Content,
		Links:       data.Links,
		IsImportant: data.IsImportant,
		IsRealtime:  data.IsRealtime,
		ReadAt:      nil,
		RecipientID: user.ID,
		SenderID:    &client.ID,
	}

	if data.IsRealtime {
		if err := services.PushNotification(notification); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else {
		if err := services.NewNotification(notification); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
