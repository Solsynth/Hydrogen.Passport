package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/gofiber/fiber/v2"
)

func getTickets(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var count int64
	var tickets []models.AuthTicket
	if err := database.C.
		Where(&models.AuthTicket{AccountID: user.ID}).
		Model(&models.AuthTicket{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.C.
		Order("created_at desc").
		Where(&models.AuthTicket{AccountID: user.ID}).
		Limit(take).
		Offset(offset).
		Find(&tickets).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  tickets,
	})
}
