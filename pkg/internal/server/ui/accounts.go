package ui

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"html/template"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/sujit-baniya/flash"
)

func selfUserinfoPage(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return DoAuthRedirect(c)
	}
	user := c.Locals("user").(models.Account)

	var data models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		Preload("Profile").
		Preload("PersonalPage").
		Preload("Contacts").
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	birthday := "Unknown"
	if data.Profile.Birthday != nil {
		birthday = data.Profile.Birthday.Format(time.RFC822)
	}

	doc := parser.
		NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock).
		Parse([]byte(data.PersonalPage.Content))

	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank})

	return c.Render("views/users/me", fiber.Map{
		"info":          flash.Get(c)["message"],
		"uid":           fmt.Sprintf("%08d", data.ID),
		"joined_at":     data.CreatedAt.Format(time.RFC822),
		"birthday_at":   birthday,
		"personal_page": template.HTML(markdown.Render(doc, renderer)),
		"userinfo":      data,
		"avatar":        data.GetAvatar(),
		"banner":        data.GetBanner(),
	}, "views/layouts/user-center")
}
