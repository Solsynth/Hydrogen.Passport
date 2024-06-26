package api

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func MapAPIs(app *fiber.App) {
	app.Get("/.well-known/openid-configuration", getOidcConfiguration)

	api := app.Group("/api").Name("API")
	{
		notify := api.Group("/notifications").Name("Notifications API")
		{
			notify.Get("/", getNotifications)
			notify.Post("/subscribe", addNotifySubscriber)
			notify.Put("/batch/read", markNotificationReadBatch)
			notify.Put("/:notificationId/read", markNotificationRead)
		}

		me := api.Group("/users/me").Name("Myself Operations")
		{

			me.Get("/avatar", getAvatar)
			me.Get("/banner", getBanner)
			me.Put("/avatar", setAvatar)
			me.Put("/banner", setBanner)

			me.Get("/", getUserinfo)
			me.Put("/", editUserinfo)
			me.Get("/events", getEvents)
			me.Get("/tickets", getTickets)
			me.Delete("/tickets/:ticketId", killTicket)

			me.Post("/confirm", doRegisterConfirm)

			me.Post("/status", setStatus)

			friends := me.Group("/friends").Name("Friends")
			{
				friends.Get("/", listFriendship)
				friends.Get("/:relatedId", getFriendship)
				friends.Post("/", makeFriendship)
				friends.Post("/:relatedId", makeFriendship)
				friends.Put("/:relatedId", editFriendship)
				friends.Delete("/:relatedId", deleteFriendship)
			}
		}

		directory := api.Group("/users/:alias").Name("User Directory")
		{
			directory.Get("/", getOtherUserinfo)
			directory.Get("/status", getStatus)
		}

		api.Post("/users", doRegister)

		auth := api.Group("/auth").Name("Auth")
		{
			auth.Post("/", doAuthenticate)
			auth.Post("/mfa", doMultiFactorAuthenticate)
			auth.Post("/token", getToken)

			auth.Get("/tickets/:ticketId", getTicket)

			auth.Get("/factors", getAvailableFactors)
			auth.Post("/factors/:factorId", requestFactorToken)

			auth.Get("/o/authorize", tryAuthorizeThirdClient)
			auth.Post("/o/authorize", authorizeThirdClient)
		}

		realms := api.Group("/realms").Name("Realms API")
		{
			realms.Get("/", listCommunityRealm)
			realms.Get("/me", listOwnedRealm)
			realms.Get("/me/available", listAvailableRealm)
			realms.Get("/:realm", getRealm)
			realms.Get("/:realm/members", listRealmMembers)
			realms.Get("/:realm/members/me", getMyRealmMember)
			realms.Post("/", createRealm)
			realms.Put("/:realmId", editRealm)
			realms.Delete("/:realmId", deleteRealm)
			realms.Post("/:realm/members", addRealmMember)
			realms.Delete("/:realm/members", removeRealmMember)
			realms.Delete("/:realm/members/me", leaveRealm)
		}

		developers := api.Group("/dev").Name("Developers API")
		{
			developers.Post("/notify", notifyUser)
		}

		api.Use(func(c *fiber.Ctx) error {
			if err := exts.EnsureAuthenticated(c); err != nil {
				return err
			}
			return c.Next()
		}).Get("/ws", websocket.New(listenWebsocket))

		api.All("/*", func(c *fiber.Ctx) error {
			return fiber.ErrNotFound
		})
	}
}
