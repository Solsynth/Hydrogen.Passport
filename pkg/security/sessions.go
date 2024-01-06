package security

import (
	"fmt"
	"strconv"
	"time"

	"code.smartsheep.studio/hydrogen/passport/pkg/database"
	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func GrantSession(challenge models.AuthChallenge, claims []string, expired *time.Time, available *time.Time) (models.AuthSession, error) {
	var session models.AuthSession
	if err := challenge.IsAvailable(); err != nil {
		return session, err
	}
	if challenge.Progress < challenge.Requirements {
		return session, fmt.Errorf("challenge haven't passed")
	}

	challenge.State = models.FinishChallengeState

	session = models.AuthSession{
		Claims:       claims,
		Challenge:    challenge,
		GrantToken:   uuid.NewString(),
		AccessToken:  uuid.NewString(),
		RefreshToken: uuid.NewString(),
		ExpiredAt:    expired,
		AvailableAt:  available,
		AccountID:    challenge.AccountID,
	}

	if err := database.C.Save(&challenge).Error; err != nil {
		return session, err
	} else if err := database.C.Save(&session).Error; err != nil {
		return session, err
	}

	return session, nil
}

func GetToken(session models.AuthSession, aud ...string) (string, string, error) {
	var refresh, access string
	if err := session.IsAvailable(); err != nil {
		return refresh, access, err
	}

	var err error

	sub := strconv.Itoa(int(session.ID))
	access, err = EncodeJwt(session.AccessToken, nil, JwtAccessType, sub, aud, time.Now().Add(30*time.Minute))
	if err != nil {
		return refresh, access, err
	}
	refresh, err = EncodeJwt(session.RefreshToken, nil, JwtRefreshType, sub, aud, time.Now().Add(30*24*time.Hour))
	if err != nil {
		return refresh, access, err
	}

	session.LastGrantAt = lo.ToPtr(time.Now())
	database.C.Save(&session)

	return access, refresh, nil
}

func ExchangeToken(token string, aud ...string) (string, string, error) {
	var session models.AuthSession
	if err := database.C.Where(models.AuthSession{GrantToken: token}).First(&session).Error; err != nil {
		return "404", "403", err
	} else if session.LastGrantAt != nil {
		return "404", "403", fmt.Errorf("session was granted the first token, use refresh token instead")
	}

	return GetToken(session, aud...)
}

func RefreshToken(token string, aud ...string) (string, string, error) {
	parseInt := func(str string) int {
		val, _ := strconv.Atoi(str)
		return val
	}

	var session models.AuthSession
	if claims, err := DecodeJwt(token); err != nil {
		return "404", "403", err
	} else if claims.Type != JwtRefreshType {
		return "404", "403", fmt.Errorf("invalid token type, expected refresh token")
	} else if err := database.C.Where(models.AuthSession{
		BaseModel: models.BaseModel{ID: uint(parseInt(claims.Subject))},
	}).First(&session).Error; err != nil {
		return "404", "403", err
	}

	return GetToken(session, aud...)
}
