package server

import (
	"strings"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func authMiddleware(c *fiber.Ctx) error {
	var token string
	if cookie := c.Cookies(services.CookieAccessKey); len(cookie) > 0 {
		token = cookie
	}
	if header := c.Get(fiber.HeaderAuthorization); len(header) > 0 {
		tk := strings.Replace(header, "Bearer", "", 1)
		token = strings.TrimSpace(tk)
	}
	if query := c.Query("tk"); len(query) > 0 {
		token = strings.TrimSpace(query)
	}

	c.Locals("token", token)

	if err := authFunc(c); err != nil {
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

	rtk := c.Cookies(services.CookieRefreshKey)
	if ctx, perms, atk, rtk, err := services.Authenticate(token, rtk, 0); err == nil {
		if atk != token {
			services.SetJwtCookieSet(c, atk, rtk)
		}
		c.Locals("permissions", perms)
		c.Locals("principal", ctx.Account)
		return nil
	} else {
		return err
	}
}
