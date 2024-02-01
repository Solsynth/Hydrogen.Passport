package server

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func getChallenges(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var count int64
	var challenges []models.AuthChallenge
	if err := database.C.
		Where(&models.AuthChallenge{AccountID: user.ID}).
		Model(&models.AuthChallenge{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.C.
		Where(&models.AuthChallenge{AccountID: user.ID}).
		Limit(take).
		Offset(offset).
		Find(&challenges).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  challenges,
	})
}

func getSessions(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var count int64
	var sessions []models.AuthSession
	if err := database.C.
		Where(&models.AuthSession{AccountID: user.ID}).
		Model(&models.AuthSession{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.C.
		Where(&models.AuthSession{AccountID: user.ID}).
		Limit(take).
		Offset(offset).
		Find(&sessions).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  sessions,
	})
}
