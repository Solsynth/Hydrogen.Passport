package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
)

func LookupChallenge(id uint) (models.AuthChallenge, error) {
	var challenge models.AuthChallenge
	err := database.C.Where(models.AuthChallenge{
		BaseModel: models.BaseModel{ID: id},
	}).First(&challenge).Error

	return challenge, err
}

func LookupChallengeWithFingerprint(id uint, ip string, ua string) (models.AuthChallenge, error) {
	var challenge models.AuthChallenge
	err := database.C.Where(models.AuthChallenge{
		BaseModel: models.BaseModel{ID: id},
		IpAddress: ip,
		UserAgent: ua,
	}).First(&challenge).Error

	return challenge, err
}
