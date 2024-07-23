package api

import (
	"github.com/gofiber/fiber/v2"
)

func MapAPIs(app *fiber.App, baseURL string) {
	app.Get("/.well-known/openid-configuration", getOidcConfiguration)

	api := app.Group(baseURL).Name("API")
	{
		notify := api.Group("/notifications").Name("Notifications API")
		{
			notify.Get("/", getNotifications)
			notify.Post("/subscribe", addNotifySubscriber)
			notify.Put("/read", markNotificationReadBatch)
			notify.Put("/read/:notificationId", markNotificationRead)
		}

		api.Get("/users/lookup", lookupAccount)

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
			me.Post("/password-reset", requestResetPassword)
			me.Patch("/password-reset", confirmResetPassword)

			me.Get("/status", getMyselfStatus)
			me.Post("/status", setStatus)
			me.Put("/status", editStatus)
			me.Delete("/status", clearStatus)

			relations := me.Group("/relations").Name("Relations")
			{
				relations.Get("/", listRelationship)
				relations.Get("/:relatedId", getRelationship)
				relations.Put("/:relatedId", editRelationship)
				relations.Delete("/:relatedId", deleteRelationship)

				relations.Post("/", makeFriendship)
				relations.Post("/:relatedId", makeFriendship)
				relations.Post("/:relatedId/accept", acceptFriend)
				relations.Post("/:relatedId/decline", declineFriend)
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

		api.All("/*", func(c *fiber.Ctx) error {
			return fiber.ErrNotFound
		})
	}
}
