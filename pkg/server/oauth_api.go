package server

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"code.smartsheep.studio/hydrogen/passport/pkg/security"
	"code.smartsheep.studio/hydrogen/passport/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"strings"
	"time"
)

func preConnect(c *fiber.Ctx) error {
	id := c.Query("client_id")
	redirect := c.Query("redirect_uri")

	var client models.ThirdClient
	if err := database.C.Where(&models.ThirdClient{Alias: id}).First(&client).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if !client.IsDraft && !lo.Contains(client.Callbacks, strings.Split(redirect, "?")[0]) {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request url")
	}

	user := c.Locals("principal").(models.Account)

	var session models.AuthSession
	if err := database.C.Where(&models.AuthSession{
		AccountID: user.ID,
		ClientID:  &client.ID,
	}).First(&session).Error; err == nil {
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
			[]string{"Hydrogen.Passport", client.Alias},
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
			[]string{"Hydrogen.Passport", client.Alias},
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
