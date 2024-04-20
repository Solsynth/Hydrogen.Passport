package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"git.solsynth.dev/hydrogen/passport/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func listFriendship(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	status := c.QueryInt("status", int(models.FriendshipActive))

	if friends, err := services.ListFriend(user, models.FriendshipStatus(status)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(friends)
	}
}

func getFriendship(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if friend, err := services.GetFriendWithTwoSides(user.ID, related.ID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(friend)
	}
}

func makeFriendship(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	relatedName := c.Query("related")
	relatedId, _ := c.ParamsInt("relatedId", 0)

	var err error
	var related models.Account
	if relatedId > 0 {
		related, err = services.GetAccount(uint(relatedId))
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	} else if len(relatedName) > 0 {
		related, err = services.LookupAccount(relatedName)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
	} else {
		return fiber.NewError(fiber.StatusBadRequest, "must one of username or user id")
	}

	friend, err := services.NewFriend(user, related, models.FriendshipPending)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(friend)
	}
}

func editFriendship(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	var data struct {
		Status uint8 `json:"status"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	friendship, err := services.GetFriendWithTwoSides(user.ID, related.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	originalStatus := friendship.Status
	friendship.Status = models.FriendshipStatus(data.Status)

	if friendship, err := services.EditFriendWithCheck(friendship, user, originalStatus); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(friendship)
	}
}

func deleteFriendship(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	friendship, err := services.GetFriendWithTwoSides(user.ID, related.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeleteFriend(friendship); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(friendship)
	}
}
