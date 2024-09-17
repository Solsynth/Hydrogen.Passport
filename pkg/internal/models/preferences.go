package models

import "gorm.io/datatypes"

type PreferenceNotification struct {
	BaseModel

	Config    datatypes.JSONMap `json:"config"`
	AccountID uint              `json:"account_id"`
	Account   Account           `json:"account"`
}
