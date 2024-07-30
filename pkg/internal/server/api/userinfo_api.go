package api

import (
	"fmt"
	"strings"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getOtherUserinfo(c *fiber.Ctx) error {
	alias := c.Params("alias")

	var account models.Account
	if err := database.C.
		Where(&models.Account{Name: alias}).
		Preload("Profile").
		Preload("Badges").
		First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	groups, err := services.GetUserAccountGroup(account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("unable to get account groups: %v", err))
	}
	for _, group := range groups {
		for k, v := range group.PermNodes {
			if _, ok := account.PermNodes[k]; !ok {
				account.PermNodes[k] = v
			}
		}
	}

	return c.JSON(account)
}

func getOtherUserinfoBatch(c *fiber.Ctx) error {
	idSet := strings.Split(c.Query("id"), ",")
	if len(idSet) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "id list is required")
	}

	var accounts []models.Account
	if err := database.C.
		Where("id IN ?", idSet).
		Find(&accounts).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(accounts)
}
