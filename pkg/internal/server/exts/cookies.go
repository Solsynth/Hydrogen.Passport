package exts

import (
	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"time"
)

func SetAuthCookies(c *fiber.Ctx, atk, rtk string) {
	c.Cookie(&fiber.Cookie{
		Name:     hyper.CookieAtk,
		Value:    atk,
		Domain:   viper.GetString("security.cookie_domain"),
		SameSite: viper.GetString("security.cookie_samesite"),
		Expires:  time.Now().Add(60 * time.Minute),
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     hyper.CookieRtk,
		Value:    rtk,
		Domain:   viper.GetString("security.cookie_domain"),
		SameSite: viper.GetString("security.cookie_samesite"),
		Expires:  time.Now().Add(24 * 30 * time.Hour),
		Path:     "/",
	})
}
