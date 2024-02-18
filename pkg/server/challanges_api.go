package server

import (
	"github.com/gofiber/fiber/v2"
	"time"

	"code.smartsheep.studio/hydrogen/identity/pkg/security"
	"code.smartsheep.studio/hydrogen/identity/pkg/services"
	"github.com/samber/lo"
)

func startChallenge(c *fiber.Ctx) error {
	var data struct {
		ID string `json:"id" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	user, err := services.LookupAccount(data.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	factors, err := services.LookupFactorsByUser(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	challenge, err := security.NewChallenge(user, factors, c.IP(), c.Get(fiber.HeaderUserAgent))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	services.AddEvent(user, "challenges.start", data.ID, c.IP(), c.Get(fiber.HeaderUserAgent))
	return c.JSON(fiber.Map{
		"display_name": user.Nick,
		"challenge":    challenge,
		"factors":      factors,
	})
}

func doChallenge(c *fiber.Ctx) error {
	var data struct {
		ChallengeID uint   `json:"challenge_id" validate:"required"`
		FactorID    uint   `json:"factor_id" validate:"required"`
		Secret      string `json:"secret" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	challenge, err := services.LookupChallengeWithFingerprint(data.ChallengeID, c.IP(), c.Get(fiber.HeaderUserAgent))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	factor, err := services.LookupFactor(data.FactorID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := security.DoChallenge(challenge, factor, data.Secret); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	challenge, err = services.LookupChallenge(data.ChallengeID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if challenge.Progress >= challenge.Requirements {
		session, err := security.GrantSession(challenge, []string{"*"}, []string{"identity"}, nil, lo.ToPtr(time.Now()))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.JSON(fiber.Map{
			"is_finished": true,
			"challenge":   challenge,
			"session":     session,
		})
	}

	return c.JSON(fiber.Map{
		"is_finished": false,
		"challenge":   challenge,
		"session":     nil,
	})
}

func exchangeToken(c *fiber.Ctx) error {
	var data struct {
		Code         string `json:"code" form:"code"`
		RefreshToken string `json:"refresh_token" form:"refresh_token"`
		ClientID     string `json:"client_id" form:"client_id"`
		ClientSecret string `json:"client_secret" form:"client_secret"`
		RedirectUri  string `json:"redirect_uri" form:"redirect_uri"`
		GrantType    string `json:"grant_type" form:"grant_type"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var err error
	var access, refresh string
	switch data.GrantType {
	case "authorization_code":
		// Authorization Code Mode
		access, refresh, err = security.ExchangeOauthToken(data.ClientID, data.ClientSecret, data.RedirectUri, data.Code)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	case "grant_token":
		// Internal Usage
		access, refresh, err = security.ExchangeToken(data.Code)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	case "refresh_token":
		// Refresh Token
		access, refresh, err = security.RefreshToken(data.RefreshToken)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported exchange token type")
	}

	return c.JSON(fiber.Map{
		"id_token":      access,
		"access_token":  access,
		"refresh_token": refresh,
		"token_type":    "Bearer",
		"expires_in":    (30 * time.Minute).Seconds(),
	})
}
