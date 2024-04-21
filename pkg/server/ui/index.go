package ui

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/gofiber/fiber/v2"
)

func MapUserInterface(A *fiber.App, authFunc func(c *fiber.Ctx, overrides ...string) error) {
	authCheckWare := func(c *fiber.Ctx) error {
		var token string
		if cookie := c.Cookies(services.CookieAccessKey); len(cookie) > 0 {
			token = cookie
		}

		c.Locals("token", token)

		if err := authFunc(c); err != nil {
			uri := c.Request().URI().FullURI()
			return c.Redirect(fmt.Sprintf("/sign-in?redirect_uri=%s", string(uri)))
		} else {
			return c.Next()
		}
	}

	pages := A.Group("/").Name("Pages")

	pages.Get("/sign-up", signupPage)
	pages.Get("/sign-in", signinPage)
	pages.Get("/mfa", mfaRequestPage)
	pages.Get("/mfa/apply", mfaApplyPage)

	pages.Post("/sign-up", signupAction)
	pages.Post("/sign-in", signinAction)
	pages.Post("/mfa", mfaRequestAction)
	pages.Post("/mfa/apply", mfaApplyAction)

	pages.Get("/users/me", authCheckWare, selfUserinfoPage)
}
