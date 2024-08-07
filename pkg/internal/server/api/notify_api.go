package api

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func notifyUser(c *fiber.Ctx) error {
	var data struct {
		ClientID     string         `json:"client_id" validate:"required"`
		ClientSecret string         `json:"client_secret" validate:"required"`
		Topic        string         `json:"type" validate:"required"`
		Title        string         `json:"subject" validate:"required,max=1024"`
		Subtitle     *string        `json:"subtitle" validate:"max=1024"`
		Body         string         `json:"content" validate:"required,max=4096"`
		Metadata     map[string]any `json:"metadata"`
		Avatar       *string        `json:"avatar"`
		Picture      *string        `json:"picture"`
		IsForcePush  bool           `json:"is_force_push"`
		IsRealtime   bool           `json:"is_realtime"`
		UserID       uint           `json:"user_id" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
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
		Topic:       data.Topic,
		Subtitle:    data.Subtitle,
		Title:       data.Title,
		Body:        data.Body,
		Metadata:    data.Metadata,
		Avatar:      data.Avatar,
		Picture:     data.Picture,
		IsRealtime:  data.IsRealtime,
		IsForcePush: data.IsForcePush,
		AccountID:   user.ID,
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
