package server

import (
	"net/http"
	"strings"
	"time"

	"code.smartsheep.studio/hydrogen/identity/pkg/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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

		api.Get("/notifications", authMiddleware, getNotifications)
		api.Put("/notifications/:notificationId/read", authMiddleware, markNotificationRead)
		api.Post("/notifications/subscribe", authMiddleware, addNotifySubscriber)

		me := api.Group("/users/me").Name("Myself Operations")
		{

			me.Put("/avatar", authMiddleware, setAvatar)
			me.Put("/banner", authMiddleware, setBanner)

			me.Get("/", authMiddleware, getUserinfo)
			me.Put("/", authMiddleware, editUserinfo)
			me.Get("/events", authMiddleware, getEvents)
			me.Get("/challenges", authMiddleware, getChallenges)
			me.Get("/sessions", authMiddleware, getSessions)
			me.Delete("/sessions/:sessionId", authMiddleware, killSession)

			me.Post("/confirm", doRegisterConfirm)
		}

		api.Post("/users", doRegister)

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
		Root:         http.FS(views.FS),
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
