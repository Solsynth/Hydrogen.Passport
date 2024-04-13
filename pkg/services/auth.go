package services

import (
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/security"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"
)

const authContextBucket = "AuthContext"

func Authenticate(access, refresh string, depth int) (user models.Account, newAccess, newRefresh string, err error) {
	var claims security.PayloadClaims
	claims, err = security.DecodeJwt(access)
	if err != nil {
		if len(refresh) > 0 && depth < 1 {
			// Auto refresh and retry
			newAccess, newRefresh, err = security.RefreshToken(refresh)
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
		user = ctx.Account
		return
	}

	ctx, err = GrantAuthContext(claims.ID)
	if err == nil {
		log.Debug().Str("jti", claims.ID).Err(lookupErr).Msg("Missed auth context cache once!")
		user = ctx.Account
		return
	}

	err = fiber.NewError(fiber.StatusUnauthorized, err.Error())
	return
}

func GetAuthContext(jti string) (models.AuthContext, error) {
	var err error
	var ctx models.AuthContext

	err = database.B.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(authContextBucket))
		if bucket == nil {
			return fmt.Errorf("unable to find auth context bucket")
		}

		raw := bucket.Get([]byte(jti))
		if raw == nil {
			return fmt.Errorf("unable to find auth context")
		} else if err := jsoniter.Unmarshal(raw, &ctx); err != nil {
			return fmt.Errorf("unable to unmarshal auth context: %v", err)
		}

		return nil
	})

	if err == nil && time.Now().Unix() >= ctx.ExpiredAt.Unix() {
		RevokeAuthContext(jti)

		return ctx, fmt.Errorf("auth context has been expired")
	}

	return ctx, err
}

func GrantAuthContext(jti string) (models.AuthContext, error) {
	var ctx models.AuthContext

	// Query data from primary database
	session, err := LookupSessionWithToken(jti)
	if err != nil {
		return ctx, fmt.Errorf("invalid auth session: %v", err)
	} else if err := session.IsAvailable(); err != nil {
		return ctx, fmt.Errorf("unavailable auth session: %v", err)
	}

	user, err := GetAccount(session.AccountID)
	if err != nil {
		return ctx, fmt.Errorf("invalid account: %v", err)
	}

	// Every context should expires in some while
	// Once user update their account info, this will have delay to update
	ctx = models.AuthContext{
		Session:   session,
		Account:   user,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	// Save data into KV cache
	return ctx, database.B.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(authContextBucket))
		if err != nil {
			return err
		}

		raw, err := jsoniter.Marshal(ctx)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(jti), raw)
	})
}

func RevokeAuthContext(jti string) error {
	return database.B.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(authContextBucket))
		if err != nil {
			return err
		}

		return bucket.Delete([]byte(jti))
	})
}
