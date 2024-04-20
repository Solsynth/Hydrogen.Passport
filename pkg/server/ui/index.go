package ui

import "github.com/gofiber/fiber/v2"

func MapUserInterface(A *fiber.App) {
	pages := A.Group("/").Name("Pages")

	pages.Get("/sign-up", signupPage)
	pages.Get("/sign-in", signinPage)

	pages.Post("/sign-up", signupAction)
	pages.Post("/sign-in", signinAction)
}
