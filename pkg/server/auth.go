package server

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/security"
	"code.smartsheep.studio/hydrogen/identity/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

var auth = keyauth.New(keyauth.Config{
	KeyLookup:  "header:Authorization",
	AuthScheme: "Bearer",
	Validator: func(c *fiber.Ctx, token string) (bool, error) {
		claims, err := security.DecodeJwt(token)
		if err != nil {
			return false, err
		}

		session, err := services.LookupSessionWithToken(claims.ID)
		if err != nil {
			return false, err
		} else if err := session.IsAvailable(); err != nil {
			return false, err
		}

		user, err := services.GetAccount(session.AccountID)
		if err != nil {
			return false, err
		}

		c.Locals("principal", user)

		return true, nil
	},
	ContextKey: "token",
})
