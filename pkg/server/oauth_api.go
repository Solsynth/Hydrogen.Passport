package server

import (
	"strings"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/security"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func preConnect(c *fiber.Ctx) error {
	id := c.Query("client_id")
	redirect := c.Query("redirect_uri")

	if len(id) <= 0 || len(redirect) <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request, missing query parameters")
	}

	var client models.ThirdClient
	if err := database.C.Where(&models.ThirdClient{Alias: id}).First(&client).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if !client.IsDraft && !lo.Contains(client.Callbacks, strings.Split(redirect, "?")[0]) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid callback url")
	}

	user := c.Locals("principal").(models.Account)

	var session models.AuthSession
	if err := database.C.Where(&models.AuthSession{
		AccountID: user.ID,
		ClientID:  &client.ID,
	}).Where("last_grant_at IS NULL").First(&session).Error; err == nil {
		if session.ExpiredAt != nil && session.ExpiredAt.Unix() < time.Now().Unix() {
			return c.JSON(fiber.Map{
				"client":  client,
				"session": nil,
			})
		} else {
			session, err = security.RegenSession(session)
		}

		return c.JSON(fiber.Map{
			"client":  client,
			"session": session,
		})
	}

	return c.JSON(fiber.Map{
		"client":  client,
		"session": nil,
	})
}

func doConnect(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id := c.Query("client_id")
	response := c.Query("response_type")
	redirect := c.Query("redirect_uri")
	scope := c.Query("scope")
	if len(scope) <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request params")
	}

	var client models.ThirdClient
	if err := database.C.Where(&models.ThirdClient{Alias: id}).First(&client).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch response {
	case "code":
		// OAuth Authorization Mode
		session, err := security.GrantOauthSession(
			user,
			client,
			strings.Split(scope, " "),
			[]string{"passport", client.Alias},
			nil,
			lo.ToPtr(time.Now()),
			c.IP(),
			c.Get(fiber.HeaderUserAgent),
		)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			services.AddEvent(user, "oauth.connect", client.Alias, c.IP(), c.Get(fiber.HeaderUserAgent))
			return c.JSON(fiber.Map{
				"session":      session,
				"redirect_uri": redirect,
			})
		}
	case "token":
		// OAuth Implicit Mode
		session, err := security.GrantOauthSession(
			user,
			client,
			strings.Split(scope, " "),
			[]string{"passport", client.Alias},
			nil,
			lo.ToPtr(time.Now()),
			c.IP(),
			c.Get(fiber.HeaderUserAgent),
		)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else if access, refresh, err := security.GetToken(session); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			services.AddEvent(user, "oauth.connect", client.Alias, c.IP(), c.Get(fiber.HeaderUserAgent))
			return c.JSON(fiber.Map{
				"access_token":  access,
				"refresh_token": refresh,
				"redirect_uri":  redirect,
				"session":       session,
			})
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported response type")
	}
}
