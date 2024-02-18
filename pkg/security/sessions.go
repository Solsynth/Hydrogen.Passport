package security

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"time"

	"code.smartsheep.studio/hydrogen/identity/pkg/database"
	"code.smartsheep.studio/hydrogen/identity/pkg/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func GrantSession(challenge models.AuthChallenge, claims, audiences []string, expired, available *time.Time) (models.AuthSession, error) {
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
		Audiences:    audiences,
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

func GrantOauthSession(user models.Account, client models.ThirdClient, claims, audiences []string, expired, available *time.Time, ip, ua string) (models.AuthSession, error) {
	session := models.AuthSession{
		Claims:    claims,
		Audiences: audiences,
		Challenge: models.AuthChallenge{
			IpAddress: ip,
			UserAgent: ua,
			RiskLevel: CalcRisk(user, ip, ua),
			State:     models.FinishChallengeState,
			AccountID: user.ID,
		},
		GrantToken:   uuid.NewString(),
		AccessToken:  uuid.NewString(),
		RefreshToken: uuid.NewString(),
		ExpiredAt:    expired,
		AvailableAt:  available,
		ClientID:     &client.ID,
		AccountID:    user.ID,
	}

	if err := database.C.Save(&session).Error; err != nil {
		return session, err
	}

	return session, nil
}

func RegenSession(session models.AuthSession) (models.AuthSession, error) {
	session.GrantToken = uuid.NewString()
	session.AccessToken = uuid.NewString()
	session.RefreshToken = uuid.NewString()
	err := database.C.Save(&session).Error
	return session, err
}

func GetToken(session models.AuthSession) (string, string, error) {
	var refresh, access string
	if err := session.IsAvailable(); err != nil {
		return refresh, access, err
	}

	accessDuration := time.Duration(viper.GetInt64("security.access_token_duration")) * time.Second
	refreshDuration := time.Duration(viper.GetInt64("security.refresh_token_duration")) * time.Second

	var err error
	sub := strconv.Itoa(int(session.AccountID))
	sed := strconv.Itoa(int(session.ID))
	access, err = EncodeJwt(session.AccessToken, JwtAccessType, sub, sed, session.Audiences, time.Now().Add(accessDuration))
	if err != nil {
		return refresh, access, err
	}
	refresh, err = EncodeJwt(session.RefreshToken, JwtRefreshType, sub, sed, session.Audiences, time.Now().Add(refreshDuration))
	if err != nil {
		return refresh, access, err
	}

	session.LastGrantAt = lo.ToPtr(time.Now())
	database.C.Save(&session)

	return access, refresh, nil
}

func ExchangeToken(token string) (string, string, error) {
	var session models.AuthSession
	if err := database.C.Where(models.AuthSession{GrantToken: token}).First(&session).Error; err != nil {
		return "404", "403", err
	} else if session.LastGrantAt != nil {
		return "404", "403", fmt.Errorf("session was granted the first token, use refresh token instead")
	} else if len(session.Audiences) > 1 {
		return "404", "403", fmt.Errorf("should use authorization code grant type")
	}

	return GetToken(session)
}

func ExchangeOauthToken(clientId, clientSecret, redirectUri, token string) (string, string, error) {
	var client models.ThirdClient
	if err := database.C.Where(models.ThirdClient{Alias: clientId}).First(&client).Error; err != nil {
		return "404", "403", err
	} else if client.Secret != clientSecret {
		return "404", "403", fmt.Errorf("invalid client secret")
	} else if !client.IsDraft && !lo.Contains(client.Callbacks, redirectUri) {
		return "404", "403", fmt.Errorf("invalid redirect uri")
	}

	var session models.AuthSession
	if err := database.C.Where(models.AuthSession{GrantToken: token}).First(&session).Error; err != nil {
		return "404", "403", err
	} else if session.LastGrantAt != nil {
		return "404", "403", fmt.Errorf("session was granted the first token, use refresh token instead")
	}

	return GetToken(session)
}

func RefreshToken(token string) (string, string, error) {
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
		BaseModel: models.BaseModel{ID: uint(parseInt(claims.SessionID))},
	}).First(&session).Error; err != nil {
		return "404", "403", err
	}

	if session, err := RegenSession(session); err != nil {
		return "404", "403", err
	} else {
		return GetToken(session)
	}
}
