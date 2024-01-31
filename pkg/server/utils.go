package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
