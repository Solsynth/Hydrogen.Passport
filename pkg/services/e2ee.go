package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"time"
)

type kexRequest struct {
	OwnerID  uint
	Conn     *websocket.Conn
	Deadline time.Time
}

var kexRequests = make(map[string]map[string]kexRequest)

func KexRequest(conn *websocket.Conn, requestId, keypairId string, ownerId uint, deadline int64) {
	if kexRequests[keypairId] == nil {
		kexRequests[keypairId] = make(map[string]kexRequest)
	}

	ddl := time.Now().Add(time.Second * time.Duration(deadline))
	request := kexRequest{
		OwnerID:  ownerId,
		Conn:     conn,
		Deadline: ddl,
	}

	flag := false
	for c := range wsConn[ownerId] {
		if c == conn {
			continue
		}

		if c.WriteMessage(1, models.UnifiedCommand{
			Action: "kex.request",
			Payload: fiber.Map{
				"request_id": requestId,
				"keypair_id": keypairId,
				"owner_id":   ownerId,
				"deadline":   deadline,
			},
		}.Marshal()) == nil {
			flag = true
		}
	}

	if flag {
		kexRequests[keypairId][requestId] = request
	}
}

func KexProvide(userId uint, requestId, keypairId string, pkt []byte) {
	if kexRequests[keypairId] == nil {
		return
	}

	val, ok := kexRequests[keypairId][requestId]
	if !ok {
		return
	} else if val.OwnerID != userId {
		return
	} else {
		_ = val.Conn.WriteMessage(1, pkt)
	}
}

func KexCleanup() {
	if len(kexRequests) <= 0 {
		return
	}

	for kp, data := range kexRequests {
		for idx, req := range data {
			if req.Deadline.Unix() <= time.Now().Unix() {
				delete(kexRequests[kp], idx)
			}
		}
	}
}
