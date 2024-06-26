package services

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"time"
)

var statusCache = make(map[uint]models.Status)

func NewStatus(user models.Account, status models.Status) (models.Status, error) {
	if err := database.C.Save(&status).Error; err != nil {
		return status, err
	} else {
		statusCache[user.ID] = status
	}
	return status, nil
}

func GetStatus(uid uint) (models.Status, error) {
	if status, ok := statusCache[uid]; ok {
		if status.ClearAt != nil && status.ClearAt.Before(time.Now()) {
			delete(statusCache, uid)
		} else {
			return status, nil
		}
	}
	var status models.Status
	if err := database.C.
		Where("account_id = ?", uid).
		Where("clear_at < ?", time.Now()).
		First(&status).Error; err != nil {
		return status, err
	} else {
		statusCache[uid] = status
	}
	return status, nil
}

func GetStatusDisturbable(uid uint) error {
	status, err := GetStatus(uid)
	isOnline := wsConn[uid] == nil || len(wsConn[uid]) < 0
	if isOnline && err != nil {
		return nil
	} else if err == nil && status.IsNoDisturb {
		return fmt.Errorf("do not disturb")
	} else {
		return fmt.Errorf("offline")
	}
}

func GetStatusOnline(uid uint) error {
	status, err := GetStatus(uid)
	isOnline := wsConn[uid] == nil || len(wsConn[uid]) < 0
	if isOnline && err != nil {
		return nil
	} else if err == nil && status.IsInvisible {
		return fmt.Errorf("invisible")
	} else {
		return fmt.Errorf("offline")
	}
}
