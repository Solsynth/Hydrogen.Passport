package server

import (
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/publisher"
	"code.smartsheep.studio/hydrogen/bus/pkg/wire"
)

var Commands = map[string]publisher.CommandManifest{
	"passport.accounts.new": {
		Name:         "Create a new account",
		Description:  "Create a new account on passport platform.",
		Requirements: wire.CommandRequirements{},
		Handle:       doRegister,
	},
	"passport.auth.challenges.new": {
		Name:         "Create a new challenge",
		Description:  "Create a new challenge to get access session.",
		Requirements: wire.CommandRequirements{},
		Handle:       startChallenge,
	},
	"passport.auth.challenges.do": {
		Name:         "Challenge a challenge",
		Description:  "Getting closer to get access session.",
		Requirements: wire.CommandRequirements{},
		Handle:       doChallenge,
	},
	"passport.auth.factor.token": {
		Name:         "Get a factor token",
		Description:  "Get the factor token to finish the challenge.",
		Requirements: wire.CommandRequirements{},
		Handle:       getFactorToken,
	},
	"passport.auth.tokens.exchange": {
		Name:         "Exchange a pair of token",
		Description:  "Use the grant token to exchange the first token pair.",
		Requirements: wire.CommandRequirements{},
		Handle:       exchangeToken,
	},
	"passport.auth.tokens.refresh": {
		Name:         "Refresh a pair token",
		Description:  "Use the refresh token to refresh the token pair.",
		Requirements: wire.CommandRequirements{},
		Handle:       refreshToken,
	},
}
