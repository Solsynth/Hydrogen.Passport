package admin

import (
	"github.com/gofiber/fiber/v2"
)

func MapAdminEndpoints(app *fiber.App) {
	admin := app.Group("/api/admin")
	{
		admin.Post("/badges", grantBadge)
		admin.Delete("/badges/:badgeId", revokeBadge)
	}
}
