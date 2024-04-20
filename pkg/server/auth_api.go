package server

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"

	"git.solsynth.dev/hydrogen/passport/pkg/services"
)

func doAuthenticate(c *fiber.Ctx) error {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	user, err := services.LookupAccount(data.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account was not found: %v", err.Error()))
	}

	ticket, err := services.NewTicket(user, c.IP(), c.Get(fiber.HeaderUserAgent))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable setup ticket: %v", err.Error()))
	}

	ticket, err = services.ActiveTicketWithPassword(ticket, data.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid password: %v", err.Error()))
	}

	return c.JSON(fiber.Map{
		"is_finished": ticket.IsAvailable(),
		"ticket":      ticket,
	})
}

func doMultiFactorAuthenticate(c *fiber.Ctx) error {
	var data struct {
		TicketID uint   `json:"ticket_id" validate:"required"`
		FactorID uint   `json:"factor_id" validate:"required"`
		Code     string `json:"code" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	ticket, err := services.GetTicket(data.TicketID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("ticket was not found: %v", err.Error()))
	}

	factor, err := services.GetFactor(data.FactorID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("factor was not found: %v", err.Error()))
	}

	ticket, err = services.ActiveTicketWithMFA(ticket, factor, data.Code)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid code: %v", err.Error()))
	}

	return c.JSON(fiber.Map{
		"is_finished": ticket.IsAvailable(),
		"ticket":      ticket,
	})
}

func getToken(c *fiber.Ctx) error {
	var data struct {
		Code         string `json:"code" form:"code"`
		RefreshToken string `json:"refresh_token" form:"refresh_token"`
		ClientID     string `json:"client_id" form:"client_id"`
		ClientSecret string `json:"client_secret" form:"client_secret"`
		Username     string `json:"username" form:"username"`
		Password     string `json:"password" form:"password"`
		RedirectUri  string `json:"redirect_uri" form:"redirect_uri"`
		GrantType    string `json:"grant_type" form:"grant_type"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	var err error
	var access, refresh string
	switch data.GrantType {
	case "refresh_token":
		// Refresh Token
		access, refresh, err = services.RefreshToken(data.RefreshToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	case "authorization_code":
		// Authorization Code Mode
		access, refresh, err = services.ExchangeOauthToken(data.ClientID, data.ClientSecret, data.RedirectUri, data.Code)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	case "password":
		// Password Mode
		user, err := services.LookupAccount(data.Username)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("account was not found: %v", err.Error()))
		}
		ticket, err := services.NewTicket(user, c.IP(), c.Get(fiber.HeaderUserAgent))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable setup ticket: %v", err.Error()))
		}
		ticket, err = services.ActiveTicketWithPassword(ticket, data.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid password: %v", err.Error()))
		} else if ticket.GrantToken == nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to get grant token to get token"))
		}
		access, refresh, err = services.ExchangeOauthToken(data.ClientID, data.ClientSecret, data.RedirectUri, *ticket.GrantToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	case "grant_token":
		// Internal Usage
		access, refresh, err = services.ExchangeToken(data.Code)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported exchange token type")
	}

	services.SetJwtCookieSet(c, access, refresh)

	return c.JSON(fiber.Map{
		"id_token":      access,
		"access_token":  access,
		"refresh_token": refresh,
		"token_type":    "Bearer",
		"expires_in":    (30 * time.Minute).Seconds(),
	})
}
