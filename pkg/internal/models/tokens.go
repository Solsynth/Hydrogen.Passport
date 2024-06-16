package models

import "time"

type MagicTokenType = int8

const (
	ConfirmMagicToken = MagicTokenType(iota)
	RegistrationMagicToken
)

type MagicToken struct {
	BaseModel

	Code      string     `json:"code"`
	Type      int8       `json:"type"`
	AssignTo  *uint      `json:"assign_to"`
	ExpiredAt *time.Time `json:"expired_at"`
}
