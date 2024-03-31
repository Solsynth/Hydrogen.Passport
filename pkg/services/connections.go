package services

import (
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"github.com/gofiber/contrib/websocket"
)

var WsConn = make(map[uint][]*websocket.Conn)
var WsNotifyQueue = make(map[uint][]byte)

func CheckOnline(user models.Account) bool {
	return len(WsConn[user.ID]) > 0
}
