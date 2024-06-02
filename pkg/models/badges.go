package models

type Badge struct {
	BaseModel

	Type      string `json:"type"`
	Label     string `json:"label"`
	Color     string `json:"color"`
	AccountID uint   `json:"account_id"`
}
