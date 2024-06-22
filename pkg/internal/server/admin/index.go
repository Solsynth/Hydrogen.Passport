package admin

import (
	"github.com/gofiber/fiber/v2"
)

func MapAdminEndpoints(A *fiber.App) {
	admin := A.Group("/api/admin")
	{
		admin.Post("/badges", grantBadge)
		admin.Delete("/badges/:badgeId", revokeBadge)
	}
}
