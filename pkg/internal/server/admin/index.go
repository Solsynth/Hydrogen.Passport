package admin

import (
	"github.com/gofiber/fiber/v2"
)

func MapAdminEndpoints(A *fiber.App, authMiddleware fiber.Handler) {
	admin := A.Group("/api/admin").Use(authMiddleware)
	{
		admin.Post("/badges", grantBadge)
		admin.Delete("/badges/:badgeId", revokeBadge)
	}
}
