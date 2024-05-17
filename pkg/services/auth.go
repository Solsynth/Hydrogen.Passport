package services

import (
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var authContextCache = make(map[string]models.AuthContext)

func Authenticate(access, refresh string, depth int) (user models.Account, perms map[string]any, newAccess, newRefresh string, err error) {
	var claims PayloadClaims
	claims, err = DecodeJwt(access)
	if err != nil {
		if len(refresh) > 0 && depth < 1 {
			// Auto refresh and retry
			newAccess, newRefresh, err = RefreshToken(refresh)
			if err == nil {
				return Authenticate(newAccess, newRefresh, depth+1)
			}
		}
		err = fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid auth key: %v", err))
		return
	}

	newAccess = access
	newRefresh = refresh

	var ctx models.AuthContext

	ctx, lookupErr := GetAuthContext(claims.ID)
	if lookupErr == nil {
		log.Debug().Str("jti", claims.ID).Msg("Hit auth context cache once!")
		perms = FilterPermNodes(ctx.Account.PermNodes, ctx.Ticket.Claims)
		user = ctx.Account
		return
	}

	ctx, err = GrantAuthContext(claims.ID)
	if err == nil {

		perms = FilterPermNodes(ctx.Account.PermNodes, ctx.Ticket.Claims)
		user = ctx.Account
		return
	}

	err = fiber.NewError(fiber.StatusUnauthorized, err.Error())
	return
}

func GetAuthContext(jti string) (models.AuthContext, error) {
	var err error
	var ctx models.AuthContext

	if val, ok := authContextCache[jti]; ok {
		ctx = val
		ctx.LastUsedAt = time.Now()
		authContextCache[jti] = ctx
		log.Debug().Str("jti", jti).Msg("Used an auth context cache")
	} else {
		ctx, err = GrantAuthContext(jti)
		log.Debug().Str("jti", jti).Msg("Created a new auth context cache")
	}

	return ctx, err
}

func GrantAuthContext(jti string) (models.AuthContext, error) {
	var ctx models.AuthContext

	// Query data from primary database
	ticket, err := GetTicketWithToken(jti)
	if err != nil {
		return ctx, fmt.Errorf("invalid auth ticket: %v", err)
	} else if err := ticket.IsAvailable(); err != nil {
		return ctx, fmt.Errorf("unavailable auth ticket: %v", err)
	}

	user, err := GetAccount(ticket.AccountID)
	if err != nil {
		return ctx, fmt.Errorf("invalid account: %v", err)
	}

	ctx = models.AuthContext{
		Ticket:     ticket,
		Account:    user,
		LastUsedAt: time.Now(),
	}

	// Put the data into memory for cache
	authContextCache[jti] = ctx

	return ctx, nil
}

func RecycleAuthContext() {
	if len(authContextCache) == 0 {
		return
	}

	affected := 0
	for key, val := range authContextCache {
		if val.LastUsedAt.Add(60*time.Second).Unix() < time.Now().Unix() {
			affected++
			delete(authContextCache, key)
		}
	}
	log.Debug().Int("affected", affected).Msg("Recycled auth context...")
}

func InvalidAuthCacheWithUser(userId uint) {
	for key, val := range authContextCache {
		if val.Account.ID == userId {
			delete(authContextCache, key)
		}
	}
}
