package services

import (
	"errors"
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"gorm.io/gorm"
	"math/rand"
)

func CheckDailyCanSign(user models.Account) error {
	var record models.SignRecord
	if err := database.C.Where("account_id = ? AND created_at::date = ?", user.ID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("unable check daliy sign record: %v", err)
	}
	return fmt.Errorf("daliy sign record exists")
}

func DailySign(user models.Account) (models.SignRecord, error) {
	tier := rand.Intn(5)
	record := models.SignRecord{
		ResultTier:       tier,
		ResultExperience: rand.Intn(tier*100) + 100,
		AccountID:        user.ID,
	}

	if err := CheckDailyCanSign(user); err != nil {
		return record, fmt.Errorf("today already signed")
	}

	if err := database.C.Save(&record).Error; err != nil {
		return record, fmt.Errorf("unable do daliy sign: %v", err)
	}

	return record, nil
}
