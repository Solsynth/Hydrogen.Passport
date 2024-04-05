package server

import (
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"git.solsynth.dev/hydrogen/identity/pkg/services"
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
	relatedId, _ := c.ParamsInt("relatedId", 0)

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
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

	if err := BindAndValidate(c, &data); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	related, err := services.GetAccount(uint(relatedId))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	friendship, err := services.GetFriendWithTwoSides(user.ID, related.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if friendship.Status == models.FriendshipPending || data.Status == uint8(models.FriendshipPending) {
		if friendship.RelatedID != user.ID {
			return fiber.NewError(fiber.StatusNotFound, "only related person can accept or revoke accept friendship")
		}
	}

	if friendship, err := services.EditFriend(friendship); err != nil {
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
