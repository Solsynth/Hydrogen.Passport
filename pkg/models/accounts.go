package models

import (
	"time"

	"gorm.io/datatypes"
)

type AccountState = int8

const (
	PendingAccountState = AccountState(iota)
	ActiveAccountState
)

type Account struct {
	BaseModel

	Name        string                       `json:"name" gorm:"uniqueIndex"`
	Nick        string                       `json:"nick"`
	State       AccountState                 `json:"state"`
	Profile     AccountProfile               `json:"profile"`
	Session     []AuthSession                `json:"sessions"`
	Challenges  []AuthChallenge              `json:"challenges"`
	Factors     []AuthFactor                 `json:"factors"`
	Contacts    []AccountContact             `json:"contacts"`
	ConfirmedAt *time.Time                   `json:"confirmed_at"`
	Permissions datatypes.JSONType[[]string] `json:"permissions"`
}

type AccountContactType = int8

const (
	EmailAccountContact = AccountContactType(iota)
)

type AccountContact struct {
	BaseModel

	Type       int8       `json:"type"`
	Content    string     `json:"content" gorm:"uniqueIndex"`
	IsPublic   bool       `json:"is_public"`
	IsPrimary  bool       `json:"is_primary"`
	VerifiedAt *time.Time `json:"verified_at"`
	AccountID  uint       `json:"account_id"`
}
