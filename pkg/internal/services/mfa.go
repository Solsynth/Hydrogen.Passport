package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetFactorName(w models.AuthFactorType, localizer *i18n.Localizer) string {
	unknown, _ := localizer.LocalizeMessage(&i18n.Message{ID: "unknown"})
	mfaEmail, _ := localizer.LocalizeMessage(&i18n.Message{ID: "mfaFactorEmail"})

	switch w {
	case models.EmailPasswordFactor:
		return mfaEmail
	default:
		return unknown
	}
}
