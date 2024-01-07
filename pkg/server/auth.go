package server

import (
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/adaptor"
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/publisher"
	"code.smartsheep.studio/hydrogen/bus/pkg/wire"
	"code.smartsheep.studio/hydrogen/passport/pkg/security"
	"code.smartsheep.studio/hydrogen/passport/pkg/services"
)

func doAuth(c *publisher.RequestCtx) error {
	token := adaptor.ParseAnyToStruct[string](c.Parameters)

	claims, err := security.DecodeJwt(token)
	if err != nil {
		return c.SendError(wire.Unauthorized, err)
	}

	session, err := services.LookupSessionWithToken(claims.ID)
	if err != nil {
		return c.SendError(wire.Unauthorized, err)
	} else if err := session.IsAvailable(); err != nil {
		return c.SendError(wire.Unauthorized, err)
	}

	user, err := services.GetAccount(session.AccountID)
	if err != nil {
		return c.SendError(wire.Unauthorized, err)
	}

	return c.SendResponse(user.Permissions.Data())
}
