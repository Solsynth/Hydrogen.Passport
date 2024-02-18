package server

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/database"
	"code.smartsheep.studio/hydrogen/identity/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"path/filepath"
)

func getAvatar(c *fiber.Ctx) error {
	id := c.Params("avatarId")
	basepath := viper.GetString("content")

	return c.SendFile(filepath.Join(basepath, id))
}

func setAvatar(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	file, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	user.Avatar = uuid.NewString()

	if err := c.SaveFile(file, user.GetAvatarPath()); err != nil {
		return err
	} else {
		database.C.Save(&user)
	}

	return c.SendStatus(fiber.StatusOK)
}
