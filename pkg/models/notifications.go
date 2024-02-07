package models

import "time"

type Notification struct {
	BaseModel

	Subject     string     `json:"subject"`
	Content     string     `json:"content"`
	IsImportant bool       `json:"is_important"`
	ReadAt      *time.Time `json:"read_at"`
	SenderID    *uint      `json:"sender_id"`
	RecipientID uint       `json:"recipient_id"`
}

const (
	NotifySubscriberFirebase = "firebase"
)

type NotificationSubscriber struct {
	BaseModel

	UserAgent string `json:"user_agent"`
	Provider  string `json:"provider"`
	DeviceID  string `json:"device_id"`
	AccountID uint   `json:"account_id"`
}