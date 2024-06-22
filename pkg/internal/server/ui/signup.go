package ui

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/sujit-baniya/flash"
)

func signupPage(c *fiber.Ctx) error {
	localizer := c.Locals("localizer").(*i18n.Localizer)

	next, _ := localizer.LocalizeMessage(&i18n.Message{ID: "next"})
	email, _ := localizer.LocalizeMessage(&i18n.Message{ID: "email"})
	nickname, _ := localizer.LocalizeMessage(&i18n.Message{ID: "nickname"})
	username, _ := localizer.LocalizeMessage(&i18n.Message{ID: "username"})
	password, _ := localizer.LocalizeMessage(&i18n.Message{ID: "password"})
	magicToken, _ := localizer.LocalizeMessage(&i18n.Message{ID: "magicToken"})
	signin, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signinTitle"})
	title, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signupTitle"})
	caption, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signupCaption"})

	return c.Render("views/signup", fiber.Map{
		"info":            flash.Get(c)["message"],
		"use_magic_token": viper.GetBool("use_registration_magic_token"),
		"i18n": fiber.Map{
			"next":        next,
			"email":       email,
			"username":    username,
			"nickname":    nickname,
			"password":    password,
			"magic_token": magicToken,
			"signin":      signin,
			"title":       title,
			"caption":     caption,
		},
	}, "views/layouts/auth")
}

func signupAction(c *fiber.Ctx) error {
	var data struct {
		Name       string `form:"name" validate:"required,lowercase,alphanum,min=4,max=16"`
		Nick       string `form:"nick" validate:"required,min=4,max=24"`
		Email      string `form:"email" validate:"required,email"`
		Password   string `form:"password" validate:"required,min=4,max=32"`
		MagicToken string `form:"magic_token"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/sign-up")
	} else if viper.GetBool("use_registration_magic_token") && len(data.MagicToken) <= 0 {
		return flash.WithInfo(c, fiber.Map{
			"message": "magic token was required",
		}).Redirect("/sign-up")
	} else if viper.GetBool("use_registration_magic_token") {
		if tk, err := services.ValidateMagicToken(data.MagicToken, models.RegistrationMagicToken); err != nil {
			return flash.WithInfo(c, fiber.Map{
				"message": fmt.Sprintf("magic token was invalid: %v", err.Error()),
			}).Redirect("/sign-up")
		} else {
			database.C.Delete(&tk)
		}
	}

	if _, err := services.CreateAccount(
		data.Name,
		data.Nick,
		data.Email,
		data.Password,
	); err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/sign-up")
	} else {
		return flash.WithInfo(c, fiber.Map{
			"message": "account has been created. now you can sign in!",
		}).Redirect(lo.FromPtr(exts.GetRedirectUri(c, "/sign-in")))
	}
}
