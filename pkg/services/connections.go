package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/gofiber/contrib/websocket"
)

type WsPushRequest struct {
	Payload     []byte
	RecipientID uint
}

var WsConn = make(map[uint][]*websocket.Conn)

func CheckOnline(user models.Account) bool {
	return len(WsConn[user.ID]) > 0
}
