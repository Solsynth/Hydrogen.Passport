package ui

import "github.com/gofiber/fiber/v2"

func MapUserInterface(A *fiber.App) {
	pages := A.Group("/").Name("Pages")
	pages.Get("/sign-in", signinPage)
}
