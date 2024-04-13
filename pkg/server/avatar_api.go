package server

import (
	"os"
	"path/filepath"

	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
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

	var previous string
	if len(user.Avatar) > 0 {
		previous = user.GetAvatarPath()
	}

	user.Avatar = uuid.NewString()

	if err := c.SaveFile(file, user.GetAvatarPath()); err != nil {
		return err
	} else {
		database.C.Save(&user)

		// Clean up
		if len(previous) > 0 {
			basepath := viper.GetString("content")
			filepath := filepath.Join(basepath, previous)
			if info, err := os.Stat(filepath); err == nil && !info.IsDir() {
				os.Remove(filepath)
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func setBanner(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	file, err := c.FormFile("banner")
	if err != nil {
		return err
	}

	var previous string
	if len(user.Banner) > 0 {
		previous = user.GetBannerPath()
	}

	user.Banner = uuid.NewString()

	if err := c.SaveFile(file, user.GetBannerPath()); err != nil {
		return err
	} else {
		database.C.Save(&user)

		// Clean up
		if len(previous) > 0 {
			basepath := viper.GetString("content")
			filepath := filepath.Join(basepath, previous)
			if info, err := os.Stat(filepath); err == nil && !info.IsDir() {
				os.Remove(filepath)
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
