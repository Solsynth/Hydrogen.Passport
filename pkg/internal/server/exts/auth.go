package exts

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/hyper"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware(c *fiber.Ctx) error {
	var atk string
	if cookie := c.Cookies(hyper.CookieAtk); len(cookie) > 0 {
		atk = cookie
	}
	if header := c.Get(fiber.HeaderAuthorization); len(header) > 0 {
		tk := strings.Replace(header, "Bearer", "", 1)
		atk = strings.TrimSpace(tk)
	}
	if tk := c.Query("tk"); len(tk) > 0 {
		atk = strings.TrimSpace(tk)
	}

	c.Locals("p_token", atk)

	rtk := c.Cookies(hyper.CookieRtk)
	if ctx, perms, newAtk, newRtk, err := services.Authenticate(atk, rtk, 0); err == nil {
		if newAtk != atk {
			SetAuthCookies(c, newAtk, newRtk)
		}
		c.Locals("permissions", perms)
		c.Locals("user", ctx.Account)
	}

	return c.Next()
}

func EnsureAuthenticated(c *fiber.Ctx) error {
	if _, ok := c.Locals("user").(models.Account); !ok {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	return nil
}

func EnsureGrantedPerm(c *fiber.Ctx, key string, val any) error {
	if err := EnsureAuthenticated(c); err != nil {
		return err
	}
	perms := c.Locals("permissions").(map[string]any)
	if !services.HasPermNode(perms, key, val) {
		return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("missing permission: %s", key))
	}
	return nil
}
