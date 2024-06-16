package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"git.solsynth.dev/hydrogen/passport/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func listRealmMembers(c *fiber.Ctx) error {
	alias := c.Params("realm")

	if realm, err := services.GetRealmWithAlias(alias); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if members, err := services.ListRealmMember(realm.ID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(members)
	}
}

func getMyRealmMember(c *fiber.Ctx) error {
	alias := c.Params("realm")
	user := c.Locals("principal").(models.Account)

	if realm, err := services.GetRealmWithAlias(alias); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if member, err := services.GetRealmMember(user.ID, realm.ID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(member)
	}
}

func addRealmMember(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	alias := c.Params("realm")

	var data struct {
		Target string `json:"target" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	realm, err := services.GetRealmWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var account models.Account
	if err := database.C.Where(&models.Account{
		Name: data.Target,
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.AddRealmMember(user, account, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func removeRealmMember(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	alias := c.Params("realm")

	var data struct {
		Target string `json:"target" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	realm, err := services.GetRealmWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var account models.Account
	if err := database.C.Where(&models.Account{
		Name: data.Target,
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.RemoveRealmMember(user, account, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func leaveRealm(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	alias := c.Params("realm")

	realm, err := services.GetRealmWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if user.ID == realm.AccountID {
		return fiber.NewError(fiber.StatusBadRequest, "you cannot leave your own realm")
	}

	var account models.Account
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: user.ID},
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.RemoveRealmMember(user, account, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}
