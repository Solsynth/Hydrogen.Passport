package services

import (
	"context"
	"fmt"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/gap"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/samber/lo"
)

var statusCache = make(map[uint]models.Status)

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
		Where("clear_at > ?", time.Now()).
		First(&status).Error; err != nil {
		return status, err
	} else {
		statusCache[uid] = status
	}
	return status, nil
}

func GetUserOnline(uid uint) bool {
	pc := proto.NewStreamControllerClient(gap.H.GetDealerGrpcConn())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := pc.CountStreamConnection(ctx, &proto.CountConnectionRequest{
		UserId: uint64(uid),
	})
	if err != nil {
		return false
	}
	return resp.Count > 0
}

func GetStatusDisturbable(uid uint) error {
	status, err := GetStatus(uid)
	isOnline := GetUserOnline(uid)
	if isOnline && err != nil {
		return nil
	} else if err == nil && status.IsNoDisturb {
		return fmt.Errorf("do not disturb")
	} else {
		return nil
	}
}

func GetStatusOnline(uid uint) error {
	status, err := GetStatus(uid)
	isOnline := GetUserOnline(uid)
	if isOnline && err != nil {
		return nil
	} else if err == nil && status.IsInvisible {
		return fmt.Errorf("invisible")
	} else if !isOnline {
		return fmt.Errorf("offline")
	} else {
		return nil
	}
}

func NewStatus(user models.Account, status models.Status) (models.Status, error) {
	if err := database.C.Save(&status).Error; err != nil {
		return status, err
	} else {
		statusCache[user.ID] = status
	}
	return status, nil
}

func EditStatus(user models.Account, status models.Status) (models.Status, error) {
	if err := database.C.Save(&status).Error; err != nil {
		return status, err
	} else {
		statusCache[user.ID] = status
	}
	return status, nil
}

func ClearStatus(user models.Account) error {
	if err := database.C.
		Where("account_id = ?", user.ID).
		Where("clear_at > ?", time.Now()).
		Updates(models.Status{ClearAt: lo.ToPtr(time.Now())}).Error; err != nil {
		return err
	} else {
		delete(statusCache, user.ID)
	}

	return nil
}
