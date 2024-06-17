package server

import (
	"context"
	"fmt"

	pcpb "git.solsynth.dev/hydrogen/paperclip/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/grpc"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func setAvatar(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		AttachmentID uint `json:"attachment" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	if _, err := grpc.Attachments.CheckAttachmentExists(context.Background(), &pcpb.AttachmentLookupRequest{
		Id:    lo.ToPtr(uint64(data.AttachmentID)),
		Usage: lo.ToPtr("p.avatar"),
	}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("avatar was not found in repository: %v", err))
	}

	user.Avatar = &data.AttachmentID

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.InvalidAuthCacheWithUser(user.ID)
	}

	return c.SendStatus(fiber.StatusOK)
}

func setBanner(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		AttachmentID uint `json:"attachment" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &data); err != nil {
		return err
	}

	if _, err := grpc.Attachments.CheckAttachmentExists(context.Background(), &pcpb.AttachmentLookupRequest{
		Id:    lo.ToPtr(uint64(data.AttachmentID)),
		Usage: lo.ToPtr("p.banner"),
	}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("banner was not found in repository: %v", err))
	}

	user.Banner = &data.AttachmentID

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		services.InvalidAuthCacheWithUser(user.ID)
	}

	return c.SendStatus(fiber.StatusOK)
}
