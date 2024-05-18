package ui

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/viper"
	"github.com/sujit-baniya/flash"
	"html/template"
	"time"
)

func otherUserinfoPage(c *fiber.Ctx) error {
	name := c.Params("account")

	var data models.Account
	if err := database.C.
		Where(&models.Account{Name: name}).
		Preload("Profile").
		Preload("PersonalPage").
		Preload("Contacts").
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var birthday = "Unknown"
	if data.Profile.Birthday != nil {
		birthday = data.Profile.Birthday.Format(time.RFC822)
	}

	doc := parser.
		NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock).
		Parse([]byte(data.PersonalPage.Content))

	renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank})

	return c.Render("views/users/directory/userinfo", fiber.Map{
		"info":          flash.Get(c)["message"],
		"uid":           fmt.Sprintf("%08d", data.ID),
		"joined_at":     data.CreatedAt.Format(time.RFC822),
		"birthday_at":   birthday,
		"personal_page": template.HTML(markdown.Render(doc, renderer)),
		"userinfo":      data,
		"avatar":        fmt.Sprintf("%s/api/attachments/%s", viper.GetString("paperclip.endpoint"), data.Avatar),
		"banner":        fmt.Sprintf("%s/api/attachments/%s", viper.GetString("paperclip.endpoint"), data.Banner),
	}, "views/layouts/user-center")
}
