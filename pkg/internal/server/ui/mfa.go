package ui

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/samber/lo"
	"github.com/sujit-baniya/flash"
)

func mfaRequestPage(c *fiber.Ctx) error {
	ticketId := c.QueryInt("ticket", 0)

	ticket, err := services.GetTicket(uint(ticketId))
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": "you must provide ticket id to perform multi-factor authenticate",
		}).Redirect("/sign-in")
	}
	user, err := services.GetAccount(ticket.AccountID)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": "ticket related user just weirdly disappear",
		}).Redirect("/sign-in")
	}
	factors, err := services.ListUserFactor(user.ID)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable to get your factors: %v", err.Error()),
		}).Redirect("/sign-in")
	}

	factors = lo.Filter(factors, func(item models.AuthFactor, index int) bool {
		return item.Type != models.PasswordAuthFactor
	})

	localizer := c.Locals("localizer").(*i18n.Localizer)

	next, _ := localizer.LocalizeMessage(&i18n.Message{ID: "next"})
	title, _ := localizer.LocalizeMessage(&i18n.Message{ID: "mfaTitle"})
	caption, _ := localizer.LocalizeMessage(&i18n.Message{ID: "mfaCaption"})

	return c.Render("views/mfa", fiber.Map{
		"info":         flash.Get(c)["message"],
		"redirect_uri": flash.Get(c)["redirect_uri"],
		"ticket_id":    ticket.ID,
		"factors": lo.Map(factors, func(item models.AuthFactor, index int) fiber.Map {
			return fiber.Map{
				"name": services.GetFactorName(item.Type, localizer),
				"id":   item.ID,
			}
		}),
		"i18n": fiber.Map{
			"next":    next,
			"title":   title,
			"caption": caption,
		},
	}, "views/layouts/auth")
}

func mfaRequestAction(c *fiber.Ctx) error {
	var data struct {
		TicketID uint `form:"ticket_id" validate:"required"`
		FactorID uint `form:"factor_id" validate:"required"`
	}

	redirectBackUri := "/sign-in"
	err := utils.BindAndValidate(c, &data)

	if data.TicketID > 0 {
		redirectBackUri = fmt.Sprintf("/mfa?ticket=%d", data.TicketID)
	}

	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": err.Error(),
		}).Redirect(redirectBackUri)
	}

	factor, err := services.GetFactor(data.FactorID)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("factor was not found: %v", err.Error()),
		}).Redirect(redirectBackUri)
	}

	_, err = services.GetFactorCode(factor)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable to get factor code: %v", err.Error()),
		}).Redirect(redirectBackUri)
	}

	return flash.WithData(c, fiber.Map{
		"redirect_uri": utils.GetRedirectUri(c),
	}).Redirect(fmt.Sprintf("/mfa/apply?ticket=%d&factor=%d", data.TicketID, factor.ID))
}

func mfaApplyPage(c *fiber.Ctx) error {
	ticketId := c.QueryInt("ticket", 0)
	factorId := c.QueryInt("factor", 0)

	ticket, err := services.GetTicket(uint(ticketId))
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable to find your ticket: %v", err.Error()),
		}).Redirect("/sign-in")
	}
	factor, err := services.GetFactor(uint(factorId))
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable to find your factors: %v", err.Error()),
		}).Redirect("/sign-in")
	}

	localizer := c.Locals("localizer").(*i18n.Localizer)

	next, _ := localizer.LocalizeMessage(&i18n.Message{ID: "next"})
	password, _ := localizer.LocalizeMessage(&i18n.Message{ID: "password"})
	title, _ := localizer.LocalizeMessage(&i18n.Message{ID: "mfaTitle"})
	caption, _ := localizer.LocalizeMessage(&i18n.Message{ID: "mfaCaption"})

	return c.Render("views/mfa-apply", fiber.Map{
		"info":      flash.Get(c)["message"],
		"label":     services.GetFactorName(factor.Type, localizer),
		"ticket_id": ticket.ID,
		"factor_id": factor.ID,
		"i18n": fiber.Map{
			"next":     next,
			"password": password,
			"title":    title,
			"caption":  caption,
		},
	}, "views/layouts/auth")
}

func mfaApplyAction(c *fiber.Ctx) error {
	var data struct {
		TicketID uint   `form:"ticket_id" validate:"required"`
		FactorID uint   `form:"factor_id" validate:"required"`
		Code     string `form:"code" validate:"required"`
	}

	redirectBackUri := "/sign-in"
	err := utils.BindAndValidate(c, &data)

	if data.TicketID > 0 {
		redirectBackUri = fmt.Sprintf("/mfa/apply?ticket=%d&factor=%d", data.TicketID, data.FactorID)
	}

	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": err.Error(),
		}).Redirect(redirectBackUri)
	}

	ticket, err := services.GetTicket(data.TicketID)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("unable to find your ticket: %v", err.Error()),
		}).Redirect("/sign-in")
	}
	factor, err := services.GetFactor(data.FactorID)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("factor was not found: %v", err.Error()),
		}).Redirect(redirectBackUri)
	}

	ticket, err = services.ActiveTicketWithMFA(ticket, factor, data.Code)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("invalid multi-factor authenticate code: %v", err.Error()),
		}).Redirect(redirectBackUri)
	} else if ticket.IsAvailable() != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": "ticket weirdly still unavailable after multi-factor authenticate",
		}).Redirect("/sign-in")
	}

	access, refresh, err := services.ExchangeToken(*ticket.GrantToken)
	if err != nil {
		return flash.WithInfo(c, fiber.Map{
			"message": fmt.Sprintf("failed to exchange token: %v", err.Error()),
		}).Redirect("/sign-in")
	} else {
		services.SetJwtCookieSet(c, access, refresh)
	}

	return c.Redirect(lo.FromPtr(utils.GetRedirectUri(c, "/users/me")))
}
