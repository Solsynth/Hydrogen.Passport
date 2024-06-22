package api

import (
	"context"
	"fmt"
	pcpb "git.solsynth.dev/hydrogen/paperclip/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func setAvatar(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}

	var data struct {
		AttachmentID uint `json:"attachment" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	pc, err := gap.DiscoverPaperclip()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "attachments services was not available")
	}
	if _, err := pcpb.NewAttachmentsClient(pc).CheckAttachmentExists(context.Background(), &pcpb.AttachmentLookupRequest{
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
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		AttachmentID uint `json:"attachment" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	pc, err := gap.DiscoverPaperclip()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "attachments services was not available")
	}
	if _, err := pcpb.NewAttachmentsClient(pc).CheckAttachmentExists(context.Background(), &pcpb.AttachmentLookupRequest{
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
