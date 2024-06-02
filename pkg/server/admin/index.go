package admin

import (
	"git.solsynth.dev/hydrogen/passport/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func MapAdminEndpoints(A *fiber.App, authFunc utils.AuthFunc) {
	admin := A.Group("/api/admin").Use(authFunc)
	{
		admin.Post("/badges", grantBadge)
		admin.Delete("/badges/:badgeId", revokeBadge)
	}
}
