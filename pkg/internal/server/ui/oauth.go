package ui

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/sujit-baniya/flash"
	"html/template"
	"strings"
	"time"
)

func authorizePage(c *fiber.Ctx) error {
	localizer := c.Locals("localizer").(*i18n.Localizer)
	user := c.Locals("principal").(models.Account)

	id := c.Query("client_id")
	redirect := c.Query("redirect_uri")

	var message string
	if len(id) <= 0 || len(redirect) <= 0 {
		message = "invalid request, missing query parameters"
	}

	var client models.ThirdClient
	if err := database.C.Where(&models.ThirdClient{Alias: id}).First(&client).Error; err != nil {
		message = fmt.Sprintf("unable to find client: %v", err)
	} else if !client.IsDraft && !lo.Contains(client.Callbacks, strings.Split(redirect, "?")[0]) {
		message = "invalid callback url"
	}

	var ticket models.AuthTicket
	if err := database.C.Where(&models.AuthTicket{
		AccountID: user.ID,
		ClientID:  &client.ID,
	}).Where("last_grant_at IS NULL").First(&ticket).Error; err == nil {
		if !(ticket.ExpiredAt != nil && ticket.ExpiredAt.Unix() < time.Now().Unix()) {
			ticket, err = services.RegenSession(ticket)
			if c.Query("response_type") == "code" {
				return c.Redirect(fmt.Sprintf(
					"%s?code=%s&state=%s",
					redirect,
					*ticket.GrantToken,
					c.Query("state"),
				))
			} else if c.Query("response_type") == "token" {
				if access, refresh, err := services.GetToken(ticket); err == nil {
					return c.Redirect(fmt.Sprintf("%s?access_token=%s&refresh_token=%s&state=%s",
						redirect,
						access,
						refresh, c.Query("state"),
					))
				}
			}
		}
	}

	decline, _ := localizer.LocalizeMessage(&i18n.Message{ID: "decline"})
	approve, _ := localizer.LocalizeMessage(&i18n.Message{ID: "approve"})
	title, _ := localizer.LocalizeMessage(&i18n.Message{ID: "authorizeTitle"})
	caption, _ := localizer.LocalizeMessage(&i18n.Message{ID: "authorizeCaption"})

	qs := "/authorize?" + string(c.Request().URI().QueryString())

	return c.Render("views/authorize", fiber.Map{
		"info":       lo.Ternary[any](len(message) > 0, message, flash.Get(c)["message"]),
		"client":     client,
		"scopes":     strings.Split(c.Query("scope"), " "),
		"action_url": template.URL(qs),
		"i18n": fiber.Map{
			"approve": approve,
			"decline": decline,
			"title":   title,
			"caption": caption,
		},
	}, "views/layouts/auth")
}

func authorizeAction(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id := c.Query("client_id")
	response := c.Query("response_type")
	redirect := c.Query("redirect_uri")
	scope := c.Query("scope")

	redirectBackUri := "/authorize?" + string(c.Request().URI().QueryString())

	if len(scope) <= 0 {
		return flash.WithInfo(c, fiber.Map{
			"message": "invalid request parameters",
		}).Redirect(redirectBackUri)
	}

	var client models.ThirdClient
	if err := database.C.Where(&models.ThirdClient{Alias: id}).First(&client).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch response {
	case "code":
		// OAuth Authorization Mode
		ticket, err := services.NewOauthTicket(
			user,
			client,
			strings.Split(scope, " "),
			[]string{"passport", client.Alias},
			c.IP(),
			c.Get(fiber.HeaderUserAgent),
		)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			services.AddEvent(user, "oauth.connect", client.Alias, c.IP(), c.Get(fiber.HeaderUserAgent))
			return c.Redirect(fmt.Sprintf(
				"%s?code=%s&state=%s",
				redirect,
				*ticket.GrantToken,
				c.Query("state"),
			))
		}
	case "token":
		// OAuth Implicit Mode
		ticket, err := services.NewOauthTicket(
			user,
			client,
			strings.Split(scope, " "),
			[]string{"passport", client.Alias},
			c.IP(),
			c.Get(fiber.HeaderUserAgent),
		)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else if access, refresh, err := services.GetToken(ticket); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			services.AddEvent(user, "oauth.connect", client.Alias, c.IP(), c.Get(fiber.HeaderUserAgent))
			return c.Redirect(fmt.Sprintf("%s?access_token=%s&refresh_token=%s&state=%s",
				redirect,
				access,
				refresh, c.Query("state"),
			))
		}
	default:
		return flash.WithInfo(c, fiber.Map{
			"message": "unsupported response type",
		}).Redirect(redirectBackUri)
	}
}
