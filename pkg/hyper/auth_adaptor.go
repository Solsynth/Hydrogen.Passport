package hyper

import (
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

const CookieAtk = "__hydrogen_atk"
const CookieRtk = "__hydrogen_rtk"

func (v *HyperConn) AuthMiddleware(c *fiber.Ctx) error {
	// Detect token
	var atk string
	if cookie := c.Cookies(CookieAtk); len(cookie) > 0 {
		atk = cookie
	}
	if header := c.Get(fiber.HeaderAuthorization); len(header) > 0 {
		tk := strings.Replace(header, "Bearer", "", 1)
		atk = strings.TrimSpace(tk)
	}

	c.Locals("p_token", atk)

	if user, newAtk, newRtk, err := v.DoAuthenticate(atk, c.Cookies(CookieRtk)); err == nil {
		if newAtk != atk {
			c.Cookie(&fiber.Cookie{
				Name:     CookieAtk,
				Value:    newAtk,
				SameSite: "Lax",
				Expires:  time.Now().Add(60 * time.Minute),
				Path:     "/",
			})
			c.Cookie(&fiber.Cookie{
				Name:     CookieRtk,
				Value:    newRtk,
				SameSite: "Lax",
				Expires:  time.Now().Add(24 * 30 * time.Hour),
				Path:     "/",
			})
		}
		c.Locals("p_user", user)
		return nil
	}

	return c.Next()
}

func (v *HyperConn) EnsureAuthenticated(c *fiber.Ctx) error {
	if _, ok := c.Locals("p_user").(*proto.Userinfo); !ok {
		return fiber.NewError(fiber.StatusUnauthorized)
	}

	return nil
}

func (v *HyperConn) EnsureGrantedPerm(c *fiber.Ctx, key string, val any) error {
	if err := v.EnsureAuthenticated(c); err != nil {
		return err
	}
	encodedVal, _ := jsoniter.Marshal(val)
	if err := v.DoCheckPerm(c.Locals("p_token").(string), key, encodedVal); err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}
	return nil
}
