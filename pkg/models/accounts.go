package models

import "time"

type AccountState = int8

const (
	PendingAccountState = AccountState(iota)
	ActiveAccountState
)

type Account struct {
	BaseModel

	Name       string           `json:"name" gorm:"uniqueIndex"`
	Nick       string           `json:"nick"`
	State      AccountState     `json:"state"`
	Session    []AuthSession    `json:"sessions"`
	Challenges []AuthChallenge  `json:"challenges"`
	Factors    []AuthFactor     `json:"factors"`
	Contacts   []AccountContact `json:"contacts"`
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
