package models

import (
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"path/filepath"
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

	Name         string                       `json:"name" gorm:"uniqueIndex"`
	Nick         string                       `json:"nick"`
	Avatar       string                       `json:"avatar"`
	State        AccountState                 `json:"state"`
	Profile      AccountProfile               `json:"profile"`
	Sessions     []AuthSession                `json:"sessions"`
	Challenges   []AuthChallenge              `json:"challenges"`
	Factors      []AuthFactor                 `json:"factors"`
	Contacts     []AccountContact             `json:"contacts"`
	Events       []ActionEvent                `json:"events"`
	MagicTokens  []MagicToken                 `json:"-" gorm:"foreignKey:AssignTo"`
	ThirdClients []ThirdClient                `json:"clients"`
	ConfirmedAt  *time.Time                   `json:"confirmed_at"`
	Permissions  datatypes.JSONType[[]string] `json:"permissions"`
}

func (v Account) GetPrimaryEmail() AccountContact {
	val, _ := lo.Find(v.Contacts, func(item AccountContact) bool {
		return item.Type == EmailAccountContact && item.IsPrimary
	})
	return val
}

func (v Account) GetAvatarPath() string {
	basepath := viper.GetString("content")
	return filepath.Join(basepath, v.Avatar)
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
