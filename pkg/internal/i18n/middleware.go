package i18n

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func I18nMiddleware(c *fiber.Ctx) error {
	accept := c.Get(fiber.HeaderAcceptLanguage)
	localizer := i18n.NewLocalizer(Bundle, accept)

	c.Locals("localizer", localizer)

	return c.Next()
}
