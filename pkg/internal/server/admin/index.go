package admin

import (
	"github.com/gofiber/fiber/v2"
)

func MapAdminAPIs(app *fiber.App) {
	admin := app.Group("/api/admin")
	{
		admin.Post("/badges", grantBadge)
		admin.Delete("/badges/:badgeId", revokeBadge)

		admin.Post("/notify/all", notifyAllUser)

		admin.Put("/users/:user/permissions", editUserPermission)
		admin.Post("/users/:user/confirm", forceConfirmAccount)
	}
}
