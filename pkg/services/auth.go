package services

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/models"
	"code.smartsheep.studio/hydrogen/identity/pkg/security"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func Authenticate(access, refresh string, depth int) (models.Account, string, string, error) {
	var user models.Account
	claims, err := security.DecodeJwt(access)
	if err != nil {
		if len(refresh) > 0 && depth < 1 {
			// Auto refresh and retry
			access, refresh, err := security.RefreshToken(refresh)
			if err == nil {
				return Authenticate(access, refresh, depth+1)
			}
		}
		return user, access, refresh, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid auth key: %v", err))
	}

	session, err := LookupSessionWithToken(claims.ID)
	if err != nil {
		return user, access, refresh, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid auth session: %v", err))
	} else if err := session.IsAvailable(); err != nil {
		return user, access, refresh, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("unavailable auth session: %v", err))
	}

	user, err = GetAccount(session.AccountID)
	if err != nil {
		return user, access, refresh, fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("invalid account: %v", err))
	}

	return user, access, refresh, nil
}
