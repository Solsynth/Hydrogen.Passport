package models

import (
	"gorm.io/datatypes"
	"time"
)

type AccountProfile struct {
	BaseModel

	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Experience uint64     `json:"experience"`
	Birthday   *time.Time `json:"birthday"`
	AccountID  uint       `json:"account_id"`
}

type AccountPage struct {
	BaseModel

	Content   string                                `json:"content"`
	Script    string                                `json:"script"`
	Style     string                                `json:"style"`
	Links     datatypes.JSONSlice[AccountPageLinks] `json:"links"`
	AccountID uint                                  `json:"account_id"`
}

type AccountPageLinks struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}
