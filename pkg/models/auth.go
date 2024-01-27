package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
)

type AuthFactorType = int8

const (
	PasswordAuthFactor = AuthFactorType(iota)
	EmailPasswordFactor
)

type AuthFactor struct {
	BaseModel

	Type      int8    `json:"type"`
	Secret    string  `json:"secret"`
	Config    JSONMap `json:"config"`
	AccountID uint    `json:"account_id"`
}

type AuthSession struct {
	BaseModel

	Claims       datatypes.JSONSlice[string] `json:"claims"`
	Challenge    AuthChallenge               `json:"challenge" gorm:"foreignKey:SessionID"`
	GrantToken   string                      `json:"grant_token"`
	AccessToken  string                      `json:"access_token"`
	RefreshToken string                      `json:"refresh_token"`
	ExpiredAt    *time.Time                  `json:"expired_at"`
	AvailableAt  *time.Time                  `json:"available_at"`
	LastGrantAt  *time.Time                  `json:"last_grant_at"`
	AccountID    uint                        `json:"account_id"`
}

func (v AuthSession) IsAvailable() error {
	if v.AvailableAt != nil && time.Now().Unix() < v.AvailableAt.Unix() {
		return fmt.Errorf("session isn't available yet")
	}
	if v.ExpiredAt != nil && time.Now().Unix() > v.ExpiredAt.Unix() {
		return fmt.Errorf("session expired")
	}

	return nil
}

type AuthChallengeState = int8

const (
	ActiveChallengeState = AuthChallengeState(iota)
	ExpiredChallengeState
	FinishChallengeState
)

type AuthChallenge struct {
	BaseModel

	IpAddress        string                     `json:"ip_address"`
	UserAgent        string                     `json:"user_agent"`
	RiskLevel        int                        `json:"risk_level"`
	Progress         int                        `json:"progress"`
	Requirements     int                        `json:"requirements"`
	BlacklistFactors datatypes.JSONType[[]uint] `json:"blacklist_factors"`
	State            int8                       `json:"state"`
	ExpiredAt        time.Time                  `json:"expired_at"`
	SessionID        *uint                      `json:"session_id"`
	AccountID        uint                       `json:"account_id"`
}

func (v AuthChallenge) IsAvailable() error {
	if time.Now().Unix() > v.ExpiredAt.Unix() {
		return fmt.Errorf("challenge expired")
	}

	return nil
}
