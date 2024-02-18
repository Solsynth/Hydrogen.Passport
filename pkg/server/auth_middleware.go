package server

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/security"
	"code.smartsheep.studio/hydrogen/identity/pkg/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func authMiddleware(c *fiber.Ctx) error {
	var token string
	if cookie := c.Cookies(security.CookieAccessKey); len(cookie) > 0 {
		token = cookie
	}
	if header := c.Get(fiber.HeaderAuthorization); len(header) > 0 {
		tk := strings.Replace(header, "Bearer", "", 1)
		token = strings.TrimSpace(tk)
	}

	c.Locals("token", token)

	if err := authFunc(c); err != nil {
		fmt.Println(err)
		return err
	}

	return c.Next()
}

func authFunc(c *fiber.Ctx, overrides ...string) error {
	var token string
	if len(overrides) > 0 {
		token = overrides[0]
	} else {
		if tk, ok := c.Locals("token").(string); !ok {
			return fiber.NewError(fiber.StatusUnauthorized)
		} else {
			token = tk
		}
	}

	claims, err := security.DecodeJwt(token)
	if err != nil {
		rtk := c.Cookies(security.CookieRefreshKey)
		if len(rtk) > 0 && len(overrides) < 1 {
			// Auto refresh and retry
			access, refresh, err := security.RefreshToken(rtk)
			if err == nil {
				security.SetJwtCookieSet(c, access, refresh)
				return authFunc(c, access)
			}
		}
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid auth key: %v", err))
	}

	session, err := services.LookupSessionWithToken(claims.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid auth session: %v", err))
	} else if err := session.IsAvailable(); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("unavailable auth session: %v", err))
	}

	user, err := services.GetAccount(session.AccountID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid account: %v", err))
	}

	c.Locals("principal", user)

	return nil
}
