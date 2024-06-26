package api

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"time"
)

func getStatus(c *fiber.Ctx) error {
	alias := c.Params("alias")

	user, err := services.GetAccountWithName(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("account not found: %s", alias))
	}

	status, err := services.GetStatus(user.ID)
	disturbable := services.GetStatusDisturbable(user.ID) == nil
	online := services.GetStatusOnline(user.ID) == nil

	return c.JSON(fiber.Map{
		"status":         lo.Ternary(err == nil, &status, nil),
		"is_disturbable": disturbable,
		"is_online":      online,
	})
}

func setStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(models.Account)
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}

	var req struct {
		Type        string     `json:"type" validate:"required"`
		Label       string     `json:"label" validate:"required"`
		Attitude    uint       `json:"attitude" validate:"required"`
		IsNoDisturb bool       `json:"is_no_disturb"`
		IsInvisible bool       `json:"is_invisible"`
		ClearAt     *time.Time `json:"clear_at"`
	}

	if err := exts.BindAndValidate(c, &req); err != nil {
		return err
	}

	status := models.Status{
		Type:        req.Type,
		Label:       req.Label,
		Attitude:    models.StatusAttitude(req.Attitude),
		IsNoDisturb: req.IsNoDisturb,
		IsInvisible: req.IsInvisible,
		ClearAt:     req.ClearAt,
		AccountID:   user.ID,
	}

	if status, err := services.NewStatus(user, status); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(status)
	}
}
