package api

import (
	"context"
	"fmt"
	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"git.solsynth.dev/hydrogen/paperclip/pkg/proto"
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

	var data struct {
		AttachmentID uint `json:"attachment" validate:"required"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeFileProvider)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "attachments services was not available")
	}
	if _, err := proto.NewAttachmentsClient(pc).CheckAttachmentExists(context.Background(), &proto.AttachmentLookupRequest{
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

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeFileProvider)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "attachments services was not available")
	}
	if _, err := proto.NewAttachmentsClient(pc).CheckAttachmentExists(context.Background(), &proto.AttachmentLookupRequest{
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

func getAvatar(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if content := user.GetAvatar(); content == nil {
		return c.SendStatus(fiber.StatusNotFound)
	} else {
		return c.Redirect(*content, fiber.StatusFound)
	}
}

func getBanner(c *fiber.Ctx) error {
	if err := exts.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	if content := user.GetBanner(); content == nil {
		return c.SendStatus(fiber.StatusNotFound)
	} else {
		return c.Redirect(*content, fiber.StatusFound)
	}
}
