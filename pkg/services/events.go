package services

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
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
