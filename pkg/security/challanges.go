package security

import (
	"fmt"
	"math"
	"time"

	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"github.com/samber/lo"
	"gorm.io/datatypes"
)

func NewChallenge(account models.Account, factors []models.AuthFactor, ip, ua string) (models.AuthChallenge, error) {
	risk := 3
	var challenge models.AuthChallenge
	// Pickup any challenge if possible
	if err := database.C.Where(models.AuthChallenge{
		AccountID: account.ID,
	}).Where("state = ?", models.ActiveChallengeState).First(&challenge).Error; err == nil {
		return challenge, nil
	}

	// Reduce the risk level
	var secureFactor int64
	if err := database.C.Where(models.AuthChallenge{
		AccountID: account.ID,
		IpAddress: ip,
	}).Model(models.AuthChallenge{}).Count(&secureFactor).Error; err != nil {
		return challenge, err
	}
	if secureFactor >= 3 {
		risk -= 2
	} else if secureFactor >= 1 {
		risk -= 1
	}

	// Thinking of the requirements factors
	requirements := int(math.Max(float64(len(factors)), math.Min(float64(risk), 1)))

	challenge = models.AuthChallenge{
		IpAddress:        ip,
		UserAgent:        ua,
		RiskLevel:        risk,
		Requirements:     requirements,
		BlacklistFactors: datatypes.NewJSONType([]uint{}),
		State:            models.ActiveChallengeState,
		ExpiredAt:        time.Now().Add(2 * time.Hour),
		AccountID:        account.ID,
	}

	err := database.C.Save(&challenge).Error

	return challenge, err
}

func DoChallenge(challenge models.AuthChallenge, factor models.AuthFactor, code string) error {
	if err := challenge.IsAvailable(); err != nil {
		return err
	}
	if challenge.Progress >= challenge.Requirements {
		return fmt.Errorf("challenge already passed")
	}

	blacklist := challenge.BlacklistFactors.Data()
	if lo.Contains(blacklist, factor.ID) {
		return fmt.Errorf("factor in blacklist, please change another factor to challenge")
	}
	if err := VerifyFactor(factor, code); err != nil {
		return err
	}

	challenge.Progress++
	challenge.BlacklistFactors = datatypes.NewJSONType(append(blacklist, factor.ID))

	if err := database.C.Save(&challenge).Error; err != nil {
		return err
	}

	return nil
}
