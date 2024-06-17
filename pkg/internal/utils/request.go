package utils

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/sujit-baniya/flash"
)

var validation = validator.New(validator.WithRequiredStructEnabled())

func BindAndValidate(c *fiber.Ctx, out any) error {
	if err := c.BodyParser(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if err := validation.Struct(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}

func GetPermissions(c *fiber.Ctx) map[string]any {
	return c.Locals("permissions").(map[string]any)
}

func CheckPermissions(c *fiber.Ctx, key string, val any) error {
	if !services.HasPermNode(GetPermissions(c), key, val) {
		return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("requires permission: %s = %v", key, val))
	}
	return nil
}

func GetRedirectUri(c *fiber.Ctx, fallback ...string) *string {
	if len(c.Query("redirect_uri")) > 0 {
		return lo.ToPtr(c.Query("redirect_uri"))
	} else if val, ok := flash.Get(c)["redirect_uri"].(*string); ok {
		return val
	} else if len(fallback) > 0 {
		return &fallback[0]
	} else {
		return nil
	}
}
