package models

import (
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
