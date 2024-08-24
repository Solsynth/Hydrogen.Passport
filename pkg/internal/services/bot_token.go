package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"time"
)

func NewApiKey(user models.Account, key models.ApiKey, ip, ua string, claims []string) (models.ApiKey, error) {
	var expiredAt *time.Time
	if key.Lifecycle != nil {
		expiredAt = lo.ToPtr(time.Now().Add(time.Duration(*key.Lifecycle) * time.Second))
	}

	key.Ticket = models.AuthTicket{
		IpAddress:           ip,
		UserAgent:           ua,
		RequireMFA:          false,
		RequireAuthenticate: false,
		Claims:              claims,
		Audiences:           []string{InternalTokenAudience},
		GrantToken:          lo.ToPtr(uuid.NewString()),
		AccessToken:         lo.ToPtr(uuid.NewString()),
		RefreshToken:        lo.ToPtr(uuid.NewString()),
		AvailableAt:         lo.ToPtr(time.Now()),
		ExpiredAt:           expiredAt,
		Account:             user,
		AccountID:           user.ID,
	}

	if err := database.C.Save(&key).Error; err != nil {
		return key, err
	}
	return key, nil
}

func RollApiKey(key models.ApiKey) (models.ApiKey, error) {
	var ticket models.AuthTicket
	if err := database.C.Where("ticket_id = ?", key.TicketID).First(&ticket).Error; err != nil {
		return key, err
	}

	ticket, err := RotateTicket(ticket)
	if err != nil {
		return key, err
	} else {
		key.Ticket = ticket
	}

	return key, nil
}