package admin

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func notifyAllUser(c *fiber.Ctx) error {
	var data struct {
		Topic       string         `json:"type" validate:"required"`
		Title       string         `json:"subject" validate:"required,max=1024"`
		Subtitle    *string        `json:"subtitle" validate:"max=1024"`
		Body        string         `json:"content" validate:"required,max=4096"`
		Metadata    map[string]any `json:"metadata"`
		Avatar      *string        `json:"avatar"`
		Picture     *string        `json:"picture"`
		IsForcePush bool           `json:"is_force_push"`
		IsRealtime  bool           `json:"is_realtime"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := exts.EnsureGrantedPerm(c, "AdminNotifyAll", true); err != nil {
		return err
	}
	operator := c.Locals("user").(models.Account)

	var users []models.Account
	if err := database.C.Find(&users).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddAuditRecord(operator, "notify.all", c.IP(), c.Get(fiber.HeaderUserAgent), map[string]any{
			"payload": data,
		})
	}

	go func() {
		for _, user := range users {
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
				Account:     user,
				AccountID:   user.ID,
			}

			if data.IsRealtime {
				if err := services.PushNotification(notification); err != nil {
					log.Error().Err(err).Uint("user", user.ID).Msg("Failed to push notification...")
				}
			} else {
				if err := services.NewNotification(notification); err != nil {
					log.Error().Err(err).Uint("user", user.ID).Msg("Failed to create notification...")
				}
			}
		}
	}()

	return c.SendStatus(fiber.StatusOK)
}

func notifyOneUser(c *fiber.Ctx) error {
	var data struct {
		Topic       string         `json:"type" validate:"required"`
		Title       string         `json:"subject" validate:"required,max=1024"`
		Subtitle    *string        `json:"subtitle" validate:"max=1024"`
		Body        string         `json:"content" validate:"required,max=4096"`
		Metadata    map[string]any `json:"metadata"`
		IsForcePush bool           `json:"is_force_push"`
		IsRealtime  bool           `json:"is_realtime"`
		UserID      uint           `json:"user_id" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := exts.EnsureGrantedPerm(c, "AdminNotifyAll", true); err != nil {
		return err
	}
	operator := c.Locals("user").(models.Account)

	var user models.Account
	if err := database.C.Where("id = ?", data.UserID).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.AddAuditRecord(operator, "notify.one", c.IP(), c.Get(fiber.HeaderUserAgent), map[string]any{
			"user_id": user.ID,
			"payload": data,
		})
	}

	notification := models.Notification{
		Topic:       data.Topic,
		Subtitle:    data.Subtitle,
		Title:       data.Title,
		Body:        data.Body,
		IsRealtime:  data.IsRealtime,
		IsForcePush: data.IsForcePush,
		AccountID:   user.ID,
	}

	if data.IsRealtime {
		if err := services.PushNotification(notification); err != nil {
			log.Error().Err(err).Uint("user", user.ID).Msg("Failed to push notification...")
		}
	} else {
		if err := services.NewNotification(notification); err != nil {
			log.Error().Err(err).Uint("user", user.ID).Msg("Failed to create notification...")
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
