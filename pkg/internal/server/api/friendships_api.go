package api

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func listFriendship(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	status := c.QueryInt("status", -1)

	var err error
	var friends []models.AccountFriendship
	if status < 0 {
		if friends, err = services.ListAllFriend(user); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		if friends, err = services.ListFriend(user, models.FriendshipStatus(status)); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(friends)
}

func getFriendship(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
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
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
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
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	relatedId, _ := c.ParamsInt("relatedId", 0)

	var data struct {
		Status uint8 `json:"status"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
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
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
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
