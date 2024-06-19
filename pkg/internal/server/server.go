package server

import (
	"net/http"
	"strings"

	"github.com/gofiber/contrib/websocket"

	"git.solsynth.dev/hydrogen/passport/pkg/internal"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/i18n"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/admin"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
		notify := api.Group("/notifications").Name("Notifications API")
		{
			notify.Get("/", authMiddleware, getNotifications)
			notify.Post("/subscribe", authMiddleware, addNotifySubscriber)
			notify.Put("/batch/read", authMiddleware, markNotificationReadBatch)
			notify.Put("/:notificationId/read", authMiddleware, markNotificationRead)
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
			me.Delete("/tickets/:ticketId", authMiddleware, killSession)

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

		realms := api.Group("/realms").Name("Realms API")
		{
			realms.Get("/", listCommunityRealm)
			realms.Get("/me", authMiddleware, listOwnedRealm)
			realms.Get("/me/available", authMiddleware, listAvailableRealm)
			realms.Get("/:realm", getRealm)
			realms.Get("/:realm/members", listRealmMembers)
			realms.Get("/:realm/members/me", authMiddleware, getMyRealmMember)
			realms.Post("/", authMiddleware, createRealm)
			realms.Put("/:realmId", authMiddleware, editRealm)
			realms.Delete("/:realmId", authMiddleware, deleteRealm)
			realms.Post("/:realm/members", authMiddleware, addRealmMember)
			realms.Delete("/:realm/members", authMiddleware, removeRealmMember)
			realms.Delete("/:realm/members/me", authMiddleware, leaveRealm)
		}

		developers := api.Group("/dev").Name("Developers API")
		{
			developers.Post("/notify", notifyUser)
		}

		api.Get("/ws", authMiddleware, websocket.New(listenWebsocket))
	}

	A.Use(favicon.New(favicon.Config{
		FileSystem: http.FS(pkg.FS),
		File:       "views/favicon.png",
		URL:        "/favicon.png",
	}))

	admin.MapAdminEndpoints(A, authMiddleware)
	ui.MapUserInterface(A, authFunc)
}

func Listen() {
	if err := A.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
