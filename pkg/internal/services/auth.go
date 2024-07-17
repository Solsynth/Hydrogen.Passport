package services

import (
	"fmt"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var authContextCache sync.Map

func Authenticate(atk, rtk string, rty int) (ctx models.AuthContext, perms map[string]any, newAtk, newRtk string, err error) {
	var claims PayloadClaims
	claims, err = DecodeJwt(atk)
	if err != nil {
		if len(rtk) > 0 && rty < 1 {
			// Auto refresh and retry
			newAtk, newRtk, err = RefreshToken(rtk)
			if err == nil {
				return Authenticate(newAtk, newRtk, rty+1)
			}
		}
		err = fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid auth key: %v", err))
		return
	}

	newAtk = atk
	newRtk = rtk

	if ctx, err = GetAuthContext(claims.ID); err == nil {
		var heldPerms map[string]any
		rawHeldPerms, _ := jsoniter.Marshal(ctx.Account.PermNodes)
		_ = jsoniter.Unmarshal(rawHeldPerms, &heldPerms)

		perms = FilterPermNodes(heldPerms, ctx.Ticket.Claims)
		return
	}

	err = fiber.NewError(fiber.StatusUnauthorized, err.Error())
	return
}

func GetAuthContext(jti string) (models.AuthContext, error) {
	var err error
	var ctx models.AuthContext

	if val, ok := authContextCache.Load(jti); ok {
		ctx = val.(models.AuthContext)
		ctx.LastUsedAt = time.Now()
		authContextCache.Store(jti, ctx)
	} else {
		ctx, err = CacheAuthContext(jti)
		log.Debug().Str("jti", jti).Msg("Created a new auth context cache")
	}

	return ctx, err
}

func CacheAuthContext(jti string) (models.AuthContext, error) {
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
	authContextCache.Store(jti, ctx)

	return ctx, nil
}

func RecycleAuthContext() {
	affected := 0

	authContextCache.Range(func(key, value any) bool {
		val := value.(models.AuthContext)
		if val.LastUsedAt.Add(60*time.Second).Unix() < time.Now().Unix() {
			affected++
			authContextCache.Delete(key)
		}
		return true
	})

	log.Debug().Int("affected", affected).Msg("Recycled auth context...")
}

func InvalidAuthCacheWithUser(userId uint) {
	authContextCache.Range(func(key, value any) bool {
		val := value.(models.AuthContext)
		if val.Account.ID == userId {
			authContextCache.Delete(key)
		}
		return true
	})
}
