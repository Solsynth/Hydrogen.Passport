package models

import (
	"time"

	"gorm.io/datatypes"
)

type Notification struct {
	BaseModel

	Type        string                                `json:"type"`
	Subject     string                                `json:"subject"`
	Content     string                                `json:"content"`
	Metadata    datatypes.JSONMap                     `json:"metadata"`
	Links       datatypes.JSONSlice[NotificationLink] `json:"links"`
	IsImportant bool                                  `json:"is_important"`
	IsRealtime  bool                                  `json:"is_realtime" gorm:"-"`
	ReadAt      *time.Time                            `json:"read_at"`
	SenderID    *uint                                 `json:"sender_id"`
	RecipientID uint                                  `json:"recipient_id"`
}

// NotificationLink Used to embed into notify and render actions
type NotificationLink struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

const (
	NotifySubscriberFirebase = "firebase"
	NotifySubscriberAPNs     = "apple"
)

type NotificationSubscriber struct {
	BaseModel

	UserAgent   string `json:"user_agent"`
	Provider    string `json:"provider"`
	DeviceID    string `json:"device_id" gorm:"uniqueIndex"`
	DeviceToken string `json:"device_token"`
	AccountID   uint   `json:"account_id"`
}
