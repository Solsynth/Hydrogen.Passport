package ui

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func DoAuthRedirect(c *fiber.Ctx) error {
	uri := c.Request().URI().FullURI()
	return c.Redirect(fmt.Sprintf("/sign-in?redirect_uri=%s", string(uri)))
}

func MapUserInterface(A *fiber.App) {
	pages := A.Group("/").Name("Pages")

	pages.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/users/me")
	})

	pages.Get("/sign-up", signupPage)
	pages.Get("/sign-in", signinPage)
	pages.Get("/mfa", mfaRequestPage)
	pages.Get("/mfa/apply", mfaApplyPage)
	pages.Get("/authorize", authorizePage)

	pages.Post("/sign-up", signupAction)
	pages.Post("/sign-in", signinAction)
	pages.Post("/mfa", mfaRequestAction)
	pages.Post("/mfa/apply", mfaApplyAction)
	pages.Post("/authorize", authorizeAction)

	pages.Get("/users/me", selfUserinfoPage)
}
