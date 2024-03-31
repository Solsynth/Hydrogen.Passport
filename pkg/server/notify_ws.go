package server

import (
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"git.solsynth.dev/hydrogen/identity/pkg/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/samber/lo"
)

func listenNotifications(c *websocket.Conn) {
	user := c.Locals("principal").(models.Account)

	// Push connection
	services.WsConn[user.ID] = append(services.WsConn[user.ID], c)

	// Event loop
	var err error
	for {
		message := services.WsNotifyQueue[user.ID]

		if message != nil {
			if err = c.WriteMessage(1, message); err != nil {
				break
			}
		}
	}

	// Pop connection
	services.WsConn[user.ID] = lo.Filter(services.WsConn[user.ID], func(item *websocket.Conn, idx int) bool {
		return item != c
	})
}
