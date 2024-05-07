package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/gofiber/contrib/websocket"
)

var wsConn = make(map[uint]map[*websocket.Conn]bool)

func ClientRegister(user models.Account, conn *websocket.Conn) {
	if wsConn[user.ID] == nil {
		wsConn[user.ID] = make(map[*websocket.Conn]bool)
	}
	wsConn[user.ID][conn] = true
}

func ClientUnregister(user models.Account, conn *websocket.Conn) {
	if wsConn[user.ID] == nil {
		wsConn[user.ID] = make(map[*websocket.Conn]bool)
	}
	delete(wsConn[user.ID], conn)
}
