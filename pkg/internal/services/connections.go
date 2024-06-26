package services

import (
	"sync"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/gofiber/contrib/websocket"
)

var (
	wsMutex sync.Mutex
	wsConn  = make(map[uint]map[*websocket.Conn]bool)
)

func ClientRegister(user models.Account, conn *websocket.Conn) {
	wsMutex.Lock()
	if wsConn[user.ID] == nil {
		wsConn[user.ID] = make(map[*websocket.Conn]bool)
	}
	wsConn[user.ID][conn] = true
	wsMutex.Unlock()
}

func ClientUnregister(user models.Account, conn *websocket.Conn) {
	wsMutex.Lock()
	if wsConn[user.ID] == nil {
		wsConn[user.ID] = make(map[*websocket.Conn]bool)
	}
	delete(wsConn[user.ID], conn)
	wsMutex.Unlock()

	if status, err := GetStatus(user.ID); err != nil || !status.IsInvisible {
		if len(wsConn[user.ID]) == 0 {
			_ = SetAccountLastSeen(user.ID)
		}
	}
}
