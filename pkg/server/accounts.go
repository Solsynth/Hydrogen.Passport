package server

import (
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/adaptor"
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/publisher"
	"code.smartsheep.studio/hydrogen/bus/pkg/wire"
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"code.smartsheep.studio/hydrogen/passport/pkg/security"
)

func doRegister(c *publisher.RequestCtx) error {
	data := adaptor.ParseAnyToStruct[struct {
		Name     string `json:"name"`
		Nick     string `json:"nick"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}](c.Parameters)

	user := models.Account{
		Name:  data.Name,
		Nick:  data.Nick,
		State: models.PendingAccountState,
		Factors: []models.AuthFactor{
			{
				Type:   models.PasswordAuthFactor,
				Secret: security.HashPassword(data.Password),
			},
		},
		Contacts: []models.AccountContact{
			{
				Type:       models.EmailAccountContact,
				Content:    data.Email,
				VerifiedAt: nil,
			},
		},
	}

	if err := database.C.Create(&user).Error; err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	return c.SendResponse(user)
}
