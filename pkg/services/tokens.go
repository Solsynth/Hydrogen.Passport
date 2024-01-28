package services

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func NewMagicToken(mode models.MagicTokenType, assignTo *models.Account, expiredAt *time.Time) (models.MagicToken, error) {
	var uid uint
	if assignTo != nil {
		uid = assignTo.ID
	}

	token := models.MagicToken{
		Code:      strings.Replace(uuid.NewString(), "-", "", -1),
		Type:      mode,
		AssignTo:  &uid,
		ExpiredAt: expiredAt,
	}

	if err := database.C.Save(&token).Error; err != nil {
		return token, err
	} else {
		return token, nil
	}
}

func NotifyMagicToken(token models.MagicToken) error {
	if token.AssignTo == nil {
		return fmt.Errorf("could notify a non-assign magic token")
	}

	var user models.Account
	if err := database.C.Where(&models.MagicToken{
		AssignTo: token.AssignTo,
	}).Preload("Contacts").First(&user).Error; err != nil {
		return err
	}

	var subject string
	var content string
	switch token.Type {
	case models.ConfirmMagicToken:
		link := fmt.Sprintf("%s/users/me/confirm?tk=%s", viper.GetString("domain"), token.Code)
		subject = fmt.Sprintf("[%s] Confirm your registration", viper.GetString("name"))
		content = fmt.Sprintf("We got a create account request with this email recently.\n"+
			"So we need you to click the link below to confirm your registeration.\n"+
			"Confirmnation Link: %s\n"+
			"If you didn't do that, you can ignore this email.\n\n"+
			"%s\n"+
			"Best wishes",
			link, viper.GetString("maintainer"))
	default:
		return fmt.Errorf("unsupported magic token type to notify")
	}

	return SendMail(user.GetPrimaryEmail().Content, subject, content)
}
