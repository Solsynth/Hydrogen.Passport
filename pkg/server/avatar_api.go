package server

import (
	"context"
	"fmt"
	pcpb "git.solsynth.dev/hydrogen/paperclip/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/grpc"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func setAvatar(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		AttachmentID string `json:"attachment"`
	}

	if _, err := grpc.Attachments.CheckAttachmentExists(context.Background(), &pcpb.AttachmentLookupRequest{
		Uuid:  &data.AttachmentID,
		Usage: lo.ToPtr("p.avatar"),
	}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("avatar was not found in repository: %v", err))
	}

	user.Avatar = data.AttachmentID

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func setBanner(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	var data struct {
		AttachmentID string `json:"attachment"`
	}

	if _, err := grpc.Attachments.CheckAttachmentExists(context.Background(), &pcpb.AttachmentLookupRequest{
		Uuid:  &data.AttachmentID,
		Usage: lo.ToPtr("p.banner"),
	}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("banner was not found in repository: %v", err))
	}

	user.Banner = data.AttachmentID

	if err := database.C.Save(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
