package services

import (
	"fmt"
	"strings"
	"time"

	"git.solsynth.dev/hydrogen/identity/pkg/database"
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const ConfirmRegistrationTemplate = `Dear %s,

Thank you for choosing to register with %s. We are excited to welcome you to our community and appreciate your trust in us.

Your registration details have been successfully received, and you are now a valued member of %s. Here are the confirm link of your registration:

	%s

As a confirmed registered member, you will have access to all our services.
We encourage you to explore our services and take full advantage of the resources available to you.

Once again, thank you for choosing us. We look forward to serving you and hope you have a positive experience with us.

Best regards,
%s`

func ValidateMagicToken(code string, mode models.MagicTokenType) (models.MagicToken, error) {
	var tk models.MagicToken
	if err := database.C.Where(models.MagicToken{Code: code, Type: mode}).First(&tk).Error; err != nil {
		return tk, err
	} else if tk.ExpiredAt != nil && time.Now().Unix() >= tk.ExpiredAt.Unix() {
		return tk, fmt.Errorf("token has been expired")
	}

	return tk, nil
}

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
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: *token.AssignTo},
	}).Preload("Contacts").First(&user).Error; err != nil {
		return err
	}

	var subject string
	var content string
	switch token.Type {
	case models.ConfirmMagicToken:
		link := fmt.Sprintf("https://%s/me/confirm?tk=%s", viper.GetString("domain"), token.Code)
		subject = fmt.Sprintf("[%s] Confirm your registration", viper.GetString("name"))
		content = fmt.Sprintf(
			ConfirmRegistrationTemplate,
			user.Name,
			viper.GetString("name"),
			viper.GetString("maintainer"),
			link,
			viper.GetString("maintainer"),
		)
	default:
		return fmt.Errorf("unsupported magic token type to notify")
	}

	return SendMail(user.GetPrimaryEmail().Content, subject, content)
}
