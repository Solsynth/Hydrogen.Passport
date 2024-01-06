package server

import (
	"time"

	"code.smartsheep.studio/hydrogen/bus/pkg/kit/adaptor"
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/publisher"
	"code.smartsheep.studio/hydrogen/bus/pkg/wire"
	"code.smartsheep.studio/hydrogen/passport/pkg/security"
	"code.smartsheep.studio/hydrogen/passport/pkg/services"
	"github.com/samber/lo"
)

func startChallenge(c *publisher.RequestCtx) error {
	meta := adaptor.ParseAnyToStruct[wire.ClientMetadata](c.Metadata)
	data := adaptor.ParseAnyToStruct[struct {
		ID string `json:"id"`
	}](c.Parameters)

	user, err := services.LookupAccount(data.ID)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	}
	factors, err := services.LookupFactorsByUser(user.ID)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	challenge, err := security.NewChallenge(user, factors, meta.ClientIp, meta.UserAgent)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	return c.SendResponse(map[string]any{
		"display_name": user.Nick,
		"challenge":    challenge,
		"factors":      factors,
	})
}

func doChallenge(c *publisher.RequestCtx) error {
	meta := adaptor.ParseAnyToStruct[wire.ClientMetadata](c.Metadata)
	data := adaptor.ParseAnyToStruct[struct {
		ChallengeID uint   `json:"challenge_id"`
		FactorID    uint   `json:"factor_id"`
		Secret      string `json:"secret"`
	}](c.Parameters)

	challenge, err := services.LookupChallengeWithFingerprint(data.ChallengeID, meta.ClientIp, meta.UserAgent)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	factor, err := services.LookupFactor(data.FactorID)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	if err := security.DoChallenge(challenge, factor, data.Secret); err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	challenge, err = services.LookupChallenge(data.ChallengeID)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	} else if challenge.Progress >= challenge.Requirements {
		session, err := security.GrantSession(challenge, []string{"*"}, nil, lo.ToPtr(time.Now()))
		if err != nil {
			return c.SendError(wire.InvalidActions, err)
		}

		return c.SendResponse(map[string]any{
			"is_finished": true,
			"challenge":   challenge,
			"session":     session,
		})
	}

	return c.SendResponse(map[string]any{
		"is_finished": false,
		"challenge":   challenge,
		"session":     nil,
	})
}

func exchangeToken(c *publisher.RequestCtx) error {
	data := adaptor.ParseAnyToStruct[struct {
		GrantToken string `json:"token"`
	}](c.Parameters)

	access, refresh, err := security.ExchangeToken(data.GrantToken)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	return c.SendResponse(map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func refreshToken(c *publisher.RequestCtx) error {
	data := adaptor.ParseAnyToStruct[struct {
		RefreshToken string `json:"token"`
	}](c.Parameters)

	access, refresh, err := security.RefreshToken(data.RefreshToken)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	return c.SendResponse(map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
