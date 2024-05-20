package models

import (
	"time"

	"github.com/samber/lo"
	"gorm.io/datatypes"
)

type Account struct {
	BaseModel

	Name        string            `json:"name" gorm:"uniqueIndex"`
	Nick        string            `json:"nick"`
	Description string            `json:"description"`
	Avatar      *uint             `json:"avatar"`
	Banner      *uint             `json:"banner"`
	ConfirmedAt *time.Time        `json:"confirmed_at"`
	PermNodes   datatypes.JSONMap `json:"perm_nodes"`

	Profile         AccountProfile   `json:"profile"`
	PersonalPage    AccountPage      `json:"personal_page"`
	Contacts        []AccountContact `json:"contacts"`
	RealmIdentities []RealmMember    `json:"realm_identities"`

	Tickets []AuthTicket `json:"tickets"`
	Factors []AuthFactor `json:"factors"`

	Events      []ActionEvent `json:"events"`
	MagicTokens []MagicToken  `json:"-" gorm:"foreignKey:AssignTo"`

	ThirdClients []ThirdClient `json:"clients"`

	Notifications     []Notification           `json:"notifications" gorm:"foreignKey:RecipientID"`
	NotifySubscribers []NotificationSubscriber `json:"notify_subscribers"`

	Friendships        []AccountFriendship `json:"friendships" gorm:"foreignKey:AccountID"`
	RelatedFriendships []AccountFriendship `json:"related_friendships" gorm:"foreignKey:RelatedID"`
}

func (v Account) GetPrimaryEmail() AccountContact {
	val, _ := lo.Find(v.Contacts, func(item AccountContact) bool {
		return item.Type == EmailAccountContact && item.IsPrimary
	})
	return val
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
