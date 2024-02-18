package server

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/view"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

var A *fiber.App

func NewServer() {
	A = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "Hydrogen.Identity",
		AppName:               "Hydrogen.Identity",
		ProxyHeader:           fiber.HeaderXForwardedFor,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		EnablePrintRoutes:     viper.GetBool("debug"),
	})

	A.Use(idempotency.New())
	A.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodOptions,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	A.Use(logger.New(logger.Config{
		Format: "${status} | ${latency} | ${method} ${path}\n",
		Output: log.Logger,
	}))

	A.Get("/.well-known", getMetadata)
	A.Get("/.well-known/openid-configuration", getOidcConfiguration)

	api := A.Group("/api").Name("API")
	{
		api.Get("/avatar/:avatarId", getAvatar)
		api.Put("/avatar", authMiddleware, setAvatar)

		api.Get("/notifications", authMiddleware, getNotifications)
		api.Put("/notifications/:notificationId/read", authMiddleware, markNotificationRead)
		api.Post("/notifications/subscribe", authMiddleware, addNotifySubscriber)

		api.Get("/users/me", authMiddleware, getUserinfo)
		api.Put("/users/me", authMiddleware, editUserinfo)
		api.Get("/users/me/events", authMiddleware, getEvents)
		api.Get("/users/me/challenges", authMiddleware, getChallenges)
		api.Get("/users/me/sessions", authMiddleware, getSessions)
		api.Delete("/users/me/sessions/:sessionId", authMiddleware, killSession)

		api.Post("/users", doRegister)
		api.Post("/users/me/confirm", doRegisterConfirm)

		api.Put("/auth", startChallenge)
		api.Post("/auth", doChallenge)
		api.Post("/auth/token", exchangeToken)
		api.Post("/auth/factors/:factorId", requestFactorToken)

		api.Get("/auth/o/connect", authMiddleware, preConnect)
		api.Post("/auth/o/connect", authMiddleware, doConnect)

		developers := api.Group("/dev").Name("Developers API")
		{
			developers.Post("/notify", notifyUser)
		}
	}

	A.Use("/", cache.New(cache.Config{
		Expiration:   24 * time.Hour,
		CacheControl: true,
	}), filesystem.New(filesystem.Config{
		Root:         http.FS(view.FS),
		PathPrefix:   "dist",
		Index:        "index.html",
		NotFoundFile: "dist/index.html",
	}))
}

func Listen() {
	if err := A.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
