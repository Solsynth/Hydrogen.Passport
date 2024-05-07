package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/gofiber/contrib/websocket"
)

func listenNotifications(c *websocket.Conn) {
	user := c.Locals("principal").(models.Account)

	// Push connection
	services.ClientRegister(user, c)

	// Event loop
	var err error
	for {
		if _, _, err = c.ReadMessage(); err != nil {
			break
		}
	}

	// Pop connection
	services.ClientUnregister(user, c)
}
