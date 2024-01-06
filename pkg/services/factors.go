package services

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
)

func LookupFactor(id uint) (models.AuthFactor, error) {
	var factor models.AuthFactor
	err := database.C.Where(models.AuthFactor{
		BaseModel: models.BaseModel{ID: id},
	}).First(&factor).Error

	return factor, err
}

func LookupFactorsByUser(uid uint) ([]models.AuthFactor, error) {
	var factors []models.AuthFactor
	err := database.C.Where(models.AuthFactor{
		AccountID: uid,
	}).Find(&factors).Error

	return factors, err
}
