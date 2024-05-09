package server

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/gofiber/contrib/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
)

func listenWebsocket(c *websocket.Conn) {
	user := c.Locals("principal").(models.Account)

	// Push connection
	services.ClientRegister(user, c)

	// Event loop
	var task models.UnifiedCommand

	var messageType int
	var payload []byte
	var packet []byte
	var err error

	for {
		if messageType, packet, err = c.ReadMessage(); err != nil {
			break
		} else if err := jsoniter.Unmarshal(packet, &task); err != nil {
			_ = c.WriteMessage(messageType, models.UnifiedCommand{
				Action:  "error",
				Message: "unable to unmarshal your command, requires json request",
			}.Marshal())
			continue
		} else {
			payload, _ = jsoniter.Marshal(task.Payload)
		}

		var message *models.UnifiedCommand
		switch task.Action {
		case "kex.request":
			var req struct {
				RequestID string `json:"request_id"`
				KeypairID string `json:"keypair_id"`
				OwnerID   uint   `json:"owner_id"`
				Deadline  int64  `json:"deadline"`
			}
			_ = jsoniter.Unmarshal(payload, &req)
			if len(req.RequestID) <= 0 || len(req.KeypairID) <= 0 || req.OwnerID <= 0 {
				message = lo.ToPtr(models.UnifiedCommandFromError(fmt.Errorf("invalid request")))
			}
			services.KexRequest(c, req.RequestID, req.KeypairID, user.ID, req.Deadline)
		case "kex.provide":
			var req struct {
				RequestID string `json:"request_id"`
				KeypairID string `json:"keypair_id"`
				PublicKey []byte `json:"public_key"`
			}
			_ = jsoniter.Unmarshal(payload, &req)
			if len(req.RequestID) <= 0 || len(req.KeypairID) <= 0 {
				message = lo.ToPtr(models.UnifiedCommandFromError(fmt.Errorf("invalid request")))
			}
			services.KexProvide(user.ID, req.RequestID, req.KeypairID, packet)
		default:
			message = lo.ToPtr(models.UnifiedCommandFromError(fmt.Errorf("unknown action")))
		}

		if message != nil {
			if err = c.WriteMessage(messageType, message.Marshal()); err != nil {
				break
			}
		}
	}

	// Pop connection
	services.ClientUnregister(user, c)
}
