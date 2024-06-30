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
		Type        string                    `json:"type" validate:"required"`
		Subject     string                    `json:"subject" validate:"required,max=1024"`
		Content     string                    `json:"content" validate:"required,max=4096"`
		Metadata    map[string]any            `json:"metadata"`
		Links       []models.NotificationLink `json:"links"`
		IsForcePush bool                      `json:"is_force_push"`
		IsRealtime  bool                      `json:"is_realtime"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := exts.EnsureGrantedPerm(c, "AdminNotifyAll", true); err != nil {
		return err
	}

	var users []models.Account
	if err := database.C.Find(&users).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	go func() {
		for _, user := range users {
			notification := models.Notification{
				Type:        data.Type,
				Subject:     data.Subject,
				Content:     data.Content,
				Links:       data.Links,
				IsRealtime:  data.IsRealtime,
				IsForcePush: data.IsForcePush,
				RecipientID: user.ID,
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
