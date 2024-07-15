package models

import (
	"gorm.io/datatypes"
)

type Notification struct {
	BaseModel

	Topic    string            `json:"topic"`
	Title    string            `json:"title"`
	Subtitle *string           `json:"subtitle"`
	Body     string            `json:"body"`
	Metadata datatypes.JSONMap `json:"metadata"`
	UserID   uint              `json:"user_id"`
	SenderID *uint             `json:"sender_id"`

	IsRealtime  bool `json:"is_realtime" gorm:"-"`
	IsForcePush bool `json:"is_force_push" gorm:"-"`
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
