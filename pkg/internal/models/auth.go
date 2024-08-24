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

	Type   int8    `json:"type"`
	Secret string  `json:"-"`
	Config JSONMap `json:"config"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}

type AuthTicket struct {
	BaseModel

	Location            string                      `json:"location"`
	IpAddress           string                      `json:"ip_address"`
	UserAgent           string                      `json:"user_agent"`
	RequireMFA          bool                        `json:"require_mfa"`
	RequireAuthenticate bool                        `json:"require_authenticate"`
	Claims              datatypes.JSONSlice[string] `json:"claims"`
	Audiences           datatypes.JSONSlice[string] `json:"audiences"`
	GrantToken          *string                     `json:"grant_token"`
	AccessToken         *string                     `json:"access_token"`
	RefreshToken        *string                     `json:"refresh_token"`
	ExpiredAt           *time.Time                  `json:"expired_at"`
	AvailableAt         *time.Time                  `json:"available_at"`
	LastGrantAt         *time.Time                  `json:"last_grant_at"`
	Nonce               *string                     `json:"nonce"`
	ClientID            *uint                       `json:"client_id"`

	Account   Account `json:"account"`
	AccountID uint    `json:"account_id"`
}

func (v AuthTicket) IsAvailable() error {
	if v.RequireMFA || v.RequireAuthenticate {
		return fmt.Errorf("ticket isn't authenticated yet")
	}
	if v.AvailableAt != nil && time.Now().Unix() < v.AvailableAt.Unix() {
		return fmt.Errorf("ticket isn't available yet")
	}
	if v.ExpiredAt != nil && time.Now().Unix() > v.ExpiredAt.Unix() {
		return fmt.Errorf("ticket expired")
	}

	return nil
}

type AuthContext struct {
	Ticket     AuthTicket `json:"ticket"`
	Account    Account    `json:"account"`
	LastUsedAt time.Time  `json:"last_used_at"`
}
