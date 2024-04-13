package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
)

func AddEvent(user models.Account, event, target, ip, ua string) models.ActionEvent {
	evt := models.ActionEvent{
		Type:      event,
		Target:    target,
		IpAddress: ip,
		UserAgent: ua,
		AccountID: user.ID,
	}

	database.C.Save(&evt)

	return evt
}
