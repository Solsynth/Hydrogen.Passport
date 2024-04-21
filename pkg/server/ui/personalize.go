package ui

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/sujit-baniya/flash"
	"strings"
	"time"
)

func personalizePage(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	localizer := c.Locals("localizer").(*i18n.Localizer)

	var data models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		Preload("Profile").
		Preload("PersonalPage").
		Preload("Contacts").
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var birthday any
	if data.Profile.Birthday != nil {
		birthday = strings.SplitN(data.Profile.Birthday.Format(time.RFC3339), "T", 1)[0]
	}

	apply, _ := localizer.LocalizeMessage(&i18n.Message{ID: "apply"})
	back, _ := localizer.LocalizeMessage(&i18n.Message{ID: "back"})

	return c.Render("views/users/personalize", fiber.Map{
		"info":        flash.Get(c)["message"],
		"birthday_at": birthday,
		"userinfo":    data,
		"i18n": fiber.Map{
			"apply": apply,
			"back":  back,
		},
	}, "views/layouts/user-center")
}

func personalizeAction(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Nick        string `form:"nick" validate:"required,min=4,max=24"`
		Description string `form:"description"`
		FirstName   string `form:"first_name"`
		LastName    string `form:"last_name"`
		Birthday    string `form:"birthday"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/users/me/personalize")
	}

	var account models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		Preload("Profile").
		First(&account).Error; err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable to get your userinfo: %v", err),
		}).Redirect("/users/me/personalize")
	}

	account.Nick = data.Nick
	account.Description = data.Description
	account.Profile.FirstName = data.FirstName
	account.Profile.LastName = data.LastName

	if birthday, err := time.Parse(time.DateOnly, data.Birthday); err == nil {
		account.Profile.Birthday = lo.ToPtr(birthday)
	}

	if err := database.C.Save(&account).Error; err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable to personalize your account: %v", err),
		}).Redirect("/users/me/personalize")
	} else if err := database.C.Save(&account.Profile).Error; err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable to personalize your profile: %v", err),
		}).Redirect("/users/me/personalize")
	}

	return flash.WithInfo(c, fiber.Map{
		"message": "your account has been personalized",
	}).Redirect("/users/me")
}
