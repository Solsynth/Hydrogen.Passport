package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/utils"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getNotifications(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	only_unread := !c.QueryBool("past", false)

	tx := database.C.Where(&models.Notification{RecipientID: user.ID}).Model(&models.Notification{})
	if only_unread {
		tx = tx.Where("read_at IS NULL")
	}

	var count int64
	var notifications []models.Notification
	if err := tx.
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.
		Limit(take).
		Offset(offset).
		Order("read_at desc").
		Find(&notifications).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  notifications,
	})
}

func markNotificationRead(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("notificationId", 0)

	var notify models.Notification
	if err := database.C.Where(&models.Notification{
		BaseModel:   models.BaseModel{ID: uint(id)},
		RecipientID: user.ID,
	}).First(&notify).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	notify.ReadAt = lo.ToPtr(time.Now())

	if err := database.C.Save(&notify).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func markNotificationReadBatch(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		MessageIDs []uint `json:"messages"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := database.C.Model(&models.Notification{}).
		Where("recipient_id = ? AND id IN ?", user.ID, data.MessageIDs).
		Updates(&models.Notification{ReadAt: lo.ToPtr(time.Now())}).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func addNotifySubscriber(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Provider    string `json:"provider" validate:"required"`
		DeviceToken string `json:"device_token" validate:"required"`
		DeviceID    string `json:"device_id" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	var count int64
	if err := database.C.Where(&models.NotificationSubscriber{
		DeviceID:    data.DeviceID,
		DeviceToken: data.DeviceToken,
		AccountID:   user.ID,
	}).Model(&models.NotificationSubscriber{}).Count(&count).Error; err != nil || count > 0 {
		return c.SendStatus(fiber.StatusOK)
	}

	subscriber, err := services.AddNotifySubscriber(
		user,
		data.Provider,
		data.DeviceID,
		data.DeviceToken,
		c.Get(fiber.HeaderUserAgent),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(subscriber)
}
