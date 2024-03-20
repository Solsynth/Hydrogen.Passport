package security

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"git.solsynth.dev/hydrogen/identity/pkg/database"
	"git.solsynth.dev/hydrogen/identity/pkg/models"
	"github.com/samber/lo"
	"gorm.io/datatypes"
)

func CalcRisk(user models.Account, ip, ua string) int {
	risk := 3
	var secureFactor int64
	if err := database.C.Where(models.AuthChallenge{
		AccountID: user.ID,
		IpAddress: ip,
	}).Model(models.AuthChallenge{}).Count(&secureFactor).Error; err == nil {
		if secureFactor >= 3 {
			risk -= 3
		} else if secureFactor >= 1 {
			risk -= 2
		}
	}

	return risk
}

func NewChallenge(user models.Account, factors []models.AuthFactor, ip, ua string) (models.AuthChallenge, error) {
	var challenge models.AuthChallenge
	// Pickup any challenge if possible
	if err := database.C.Where(models.AuthChallenge{
		AccountID: user.ID,
	}).Where("state = ?", models.ActiveChallengeState).First(&challenge).Error; err == nil {
		return challenge, nil
	}

	// Calculate the risk level
	risk := CalcRisk(user, ip, ua)

	// Clamp risk in the exists requirements factor count
	requirements := lo.Clamp(risk, 1, len(factors))

	challenge = models.AuthChallenge{
		IpAddress:        ip,
		UserAgent:        ua,
		RiskLevel:        risk,
		Requirements:     requirements,
		BlacklistFactors: datatypes.NewJSONType([]uint{}),
		State:            models.ActiveChallengeState,
		ExpiredAt:        time.Now().Add(2 * time.Hour),
		AccountID:        user.ID,
	}

	err := database.C.Save(&challenge).Error

	return challenge, err
}

func DoChallenge(challenge models.AuthChallenge, factor models.AuthFactor, code string) error {
	if err := challenge.IsAvailable(); err != nil {
		challenge.State = models.ExpiredChallengeState
		database.C.Save(&challenge)
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

	// Revoke some factor passwords
	if factor.Type == models.EmailPasswordFactor {
		factor.Secret = strings.ReplaceAll(uuid.NewString(), "-", "")
		database.C.Save(&factor)
	}

	return nil
}
