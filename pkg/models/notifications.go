package models

import "time"

type Notification struct {
	BaseModel

	Subject     string     `json:"subject"`
	Content     string     `json:"content"`
	IsImportant bool       `json:"is_important"`
	ReadAt      *time.Time `json:"read_at"`
	RecipientID uint       `json:"recipient_id"`
}
