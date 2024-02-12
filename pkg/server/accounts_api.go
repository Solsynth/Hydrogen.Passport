package server

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"code.smartsheep.studio/hydrogen/passport/pkg/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func getUserinfo(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		Preload("Profile").
		Preload("Contacts").
		Preload("Notifications", "read_at IS NULL").
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var resp fiber.Map
	raw, _ := jsoniter.Marshal(data)
	jsoniter.Unmarshal(raw, &resp)

	resp["sub"] = strconv.Itoa(int(data.ID))
	resp["family_name"] = data.Profile.FirstName
	resp["given_name"] = data.Profile.LastName
	resp["name"] = data.Name
	resp["email"] = data.GetPrimaryEmail().Content
	resp["preferred_username"] = data.Nick

	if len(data.Avatar) > 0 {
		resp["picture"] = fmt.Sprintf("https://%s/api/avatar/%s", viper.GetString("domain"), data.Avatar)
	}

	return c.JSON(resp)
}

func getEvents(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var count int64
	var events []models.ActionEvent
	if err := database.C.
		Where(&models.ActionEvent{AccountID: user.ID}).
		Model(&models.ActionEvent{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := database.C.
		Order("created_at desc").
		Where(&models.ActionEvent{AccountID: user.ID}).
		Limit(take).
		Offset(offset).
		Find(&events).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  events,
	})
}

func editUserinfo(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Nick       string    `json:"nick" validate:"required,min=4,max=24"`
		FirstName  string    `json:"first_name"`
		MiddleName string    `json:"middle_name"`
		LastName   string    `json:"last_name"`
		Birthday   time.Time `json:"birthday"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var account models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		Preload("Profile").
		First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	account.Nick = data.Nick
	account.Profile.FirstName = data.FirstName
	account.Profile.MiddleName = data.MiddleName
	account.Profile.LastName = data.LastName
	account.Profile.Birthday = &data.Birthday

	if err := database.C.Save(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if err := database.C.Save(&account.Profile).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func killSession(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("sessionId", 0)

	if err := database.C.Delete(&models.AuthSession{}, &models.AuthSession{
		BaseModel: models.BaseModel{ID: uint(id)},
		AccountID: user.ID,
	}).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func doRegister(c *fiber.Ctx) error {
	var data struct {
		Name       string `json:"name" validate:"required,lowercase,alphanum,min=4,max=16"`
		Nick       string `json:"nick" validate:"required,min=4,max=24"`
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,min=4,max=32"`
		MagicToken string `json:"magic_token"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if viper.GetBool("use_registration_magic_token") && len(data.MagicToken) <= 0 {
		return fmt.Errorf("missing magic token in request")
	}

	if tk, err := services.ValidateMagicToken(data.MagicToken, models.RegistrationMagicToken); err != nil {
		return err
	} else {
		database.C.Delete(&tk)
	}

	if user, err := services.CreateAccount(
		data.Name,
		data.Nick,
		data.Email,
		data.Password,
	); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(user)
	}
}

func doRegisterConfirm(c *fiber.Ctx) error {
	var data struct {
		Code string `json:"code" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	if err := services.ConfirmAccount(data.Code); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
