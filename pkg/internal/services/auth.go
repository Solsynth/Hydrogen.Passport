package services

import (
	"context"
	"fmt"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	jsoniter "github.com/json-iterator/go"

	localCache "git.solsynth.dev/hydrogen/passport/pkg/internal/cache"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Authenticate(atk, rtk string, rty int) (ctx models.AuthContext, perms map[string]any, err error) {
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

func GetAuthContextCacheKey(jti string) string {
	return fmt.Sprintf("auth-context#%s", jti)
}

func GetAuthContext(jti string) (models.AuthContext, error) {
	var err error
	var ctx models.AuthContext

	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)
	contx := context.Background()

	if val, err := marshal.Get(contx, GetAuthContextCacheKey(jti), new(models.AuthContext)); err == nil {
		ctx = *val.(*models.AuthContext)
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
	groups, err := GetUserAccountGroup(user)
	if err != nil {
		return ctx, fmt.Errorf("unable to get account groups: %v", err)
	}

	for _, group := range groups {
		for k, v := range group.PermNodes {
			if _, ok := user.PermNodes[k]; !ok {
				user.PermNodes[k] = v
			}
		}
	}

	ctx = models.AuthContext{
		Ticket:  ticket,
		Account: user,
	}

	// Put the data into cache
	cacheManager := cache.New[any](localCache.S)
	marshal := marshaler.New(cacheManager)
	contx := context.Background()

	marshal.Set(
		contx,
		GetAuthContextCacheKey(jti),
		ctx,
		store.WithExpiration(3*time.Minute),
		store.WithTags([]string{"auth-context", fmt.Sprintf("user#%d", user.ID)}),
	)

	return ctx, nil
}

func InvalidAuthCacheWithUser(userId uint) {
	cacheManager := cache.New[any](localCache.S)
	contx := context.Background()

	cacheManager.Invalidate(
		contx,
		store.WithInvalidateTags([]string{"auth-context", fmt.Sprintf("user#%d", userId)}),
	)
}
