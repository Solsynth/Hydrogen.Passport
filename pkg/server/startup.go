package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg"
	"git.solsynth.dev/hydrogen/passport/pkg/i18n"
	"git.solsynth.dev/hydrogen/passport/pkg/server/ui"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var A *fiber.App

func NewServer() {
	templates := html.NewFileSystem(http.FS(pkg.FS), ".gohtml")

	A = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "Hydrogen.Passport",
		AppName:               "Hydrogen.Passport",
		ProxyHeader:           fiber.HeaderXForwardedFor,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
		Views:                 templates,
		ViewsLayout:           "views/index",
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

	A.Use(i18n.I18nMiddleware)

	A.Get("/.well-known", getMetadata)
	A.Get("/.well-known/openid-configuration", getOidcConfiguration)

	api := A.Group("/api").Name("API")
	{
		api.Get("/avatar/:avatarId", getAvatar)

		notify := api.Group("/notifications").Name("Notifications API")
		{
			notify.Get("/", authMiddleware, getNotifications)
			notify.Put("/:notificationId/read", authMiddleware, markNotificationRead)
			notify.Post("/subscribe", authMiddleware, addNotifySubscriber)

			notify.Get("/listen", authMiddleware, websocket.New(listenNotifications))
		}

		me := api.Group("/users/me").Name("Myself Operations")
		{

			me.Put("/avatar", authMiddleware, setAvatar)
			me.Put("/banner", authMiddleware, setBanner)

			me.Get("/", authMiddleware, getUserinfo)
			me.Get("/page", authMiddleware, getOwnPersonalPage)
			me.Put("/", authMiddleware, editUserinfo)
			me.Put("/page", authMiddleware, editPersonalPage)
			me.Get("/events", authMiddleware, getEvents)
			me.Get("/tickets", authMiddleware, getTickets)
			me.Delete("/sessions/:sessionId", authMiddleware, killSession)

			me.Post("/confirm", doRegisterConfirm)

			friends := me.Group("/friends").Name("Friends")
			{
				friends.Get("/", authMiddleware, listFriendship)
				friends.Get("/:relatedId", authMiddleware, getFriendship)
				friends.Post("/", authMiddleware, makeFriendship)
				friends.Post("/:relatedId", authMiddleware, makeFriendship)
				friends.Put("/:relatedId", authMiddleware, editFriendship)
				friends.Delete("/:relatedId", authMiddleware, deleteFriendship)
			}
		}

		directory := api.Group("/users/:alias").Name("User Directory")
		{
			directory.Get("/", getOtherUserinfo)
			directory.Get("/page", getPersonalPage)
		}

		api.Post("/users", doRegister)

		api.Post("/auth", doAuthenticate)
		api.Post("/auth/token", getToken)
		api.Post("/auth/factors/:factorId", requestFactorToken)

		api.Get("/auth/o/connect", authMiddleware, preConnect)
		api.Post("/auth/o/connect", authMiddleware, doConnect)

		developers := api.Group("/dev").Name("Developers API")
		{
			developers.Post("/notify", notifyUser)
		}
	}

	A.Use(favicon.New(favicon.Config{
		FileSystem: http.FS(pkg.FS),
		File:       "views/favicon.png",
		URL:        "/favicon.png",
	}))

	ui.MapUserInterface(A, authFunc)
}

func Listen() {
	if err := A.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
