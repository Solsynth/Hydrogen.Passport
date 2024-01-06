package server

import (
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/adaptor"
	"code.smartsheep.studio/hydrogen/bus/pkg/kit/publisher"
	"code.smartsheep.studio/hydrogen/bus/pkg/wire"
	"code.smartsheep.studio/hydrogen/passport/pkg/security"
	"code.smartsheep.studio/hydrogen/passport/pkg/services"
)

func getFactorToken(c *publisher.RequestCtx) error {
	data := adaptor.ParseAnyToStruct[struct {
		ID uint `json:"id"`
	}](c.Parameters)

	factor, err := services.LookupFactor(data.ID)
	if err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	if err := security.GetFactorCode(factor); err != nil {
		return c.SendError(wire.InvalidActions, err)
	}

	return c.SendResponse(nil)
}
