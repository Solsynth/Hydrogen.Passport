package ui

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func signinPage(c *fiber.Ctx) error {
	localizer := c.Locals("localizer").(*i18n.Localizer)

	next, _ := localizer.LocalizeMessage(&i18n.Message{ID: "next"})
	username, _ := localizer.LocalizeMessage(&i18n.Message{ID: "username"})
	password, _ := localizer.LocalizeMessage(&i18n.Message{ID: "password"})
	title, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signinTitle"})
	caption, _ := localizer.LocalizeMessage(&i18n.Message{ID: "signinCaption"})

	return c.Render("views/signin", fiber.Map{
		"i18n": fiber.Map{
			"next":     next,
			"username": username,
			"password": password,
			"title":    title,
			"caption":  caption,
		},
	}, "views/layouts/auth")
}
