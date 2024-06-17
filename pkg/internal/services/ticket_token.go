package services

import (
	"fmt"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func GetToken(ticket models.AuthTicket) (string, string, error) {
	var refresh, access string
	if err := ticket.IsAvailable(); err != nil {
		return refresh, access, err
	}
	if ticket.AccessToken == nil || ticket.RefreshToken == nil {
		return refresh, access, fmt.Errorf("unable to encode token, access or refresh token id missing")
	}

	accessDuration := time.Duration(viper.GetInt64("security.access_token_duration")) * time.Second
	refreshDuration := time.Duration(viper.GetInt64("security.refresh_token_duration")) * time.Second

	var err error
	sub := strconv.Itoa(int(ticket.AccountID))
	sed := strconv.Itoa(int(ticket.ID))
	access, err = EncodeJwt(*ticket.AccessToken, JwtAccessType, sub, sed, ticket.Audiences, time.Now().Add(accessDuration))
	if err != nil {
		return refresh, access, err
	}
	refresh, err = EncodeJwt(*ticket.RefreshToken, JwtRefreshType, sub, sed, ticket.Audiences, time.Now().Add(refreshDuration))
	if err != nil {
		return refresh, access, err
	}

	ticket.LastGrantAt = lo.ToPtr(time.Now())
	database.C.Save(&ticket)

	return access, refresh, nil
}

func ExchangeToken(token string) (string, string, error) {
	var ticket models.AuthTicket
	if err := database.C.Where(models.AuthTicket{GrantToken: &token}).First(&ticket).Error; err != nil {
		return "", "", err
	} else if ticket.LastGrantAt != nil {
		return "", "", fmt.Errorf("ticket was granted the first token, use refresh token instead")
	} else if len(ticket.Audiences) > 1 {
		return "", "", fmt.Errorf("should use authorization code grant type")
	}

	return GetToken(ticket)
}

func ExchangeOauthToken(clientId, clientSecret, redirectUri, token string) (string, string, error) {
	var client models.ThirdClient
	if err := database.C.Where(models.ThirdClient{Alias: clientId}).First(&client).Error; err != nil {
		return "", "", err
	} else if client.Secret != clientSecret {
		return "", "", fmt.Errorf("invalid client secret")
	} else if !client.IsDraft && !lo.Contains(client.Callbacks, redirectUri) {
		return "", "", fmt.Errorf("invalid redirect uri")
	}

	var ticket models.AuthTicket
	if err := database.C.Where(models.AuthTicket{GrantToken: &token}).First(&ticket).Error; err != nil {
		return "", "", err
	} else if ticket.LastGrantAt != nil {
		return "", "", fmt.Errorf("ticket was granted the first token, use refresh token instead")
	}

	return GetToken(ticket)
}

func RefreshToken(token string) (string, string, error) {
	parseInt := func(str string) int {
		val, _ := strconv.Atoi(str)
		return val
	}

	var ticket models.AuthTicket
	if claims, err := DecodeJwt(token); err != nil {
		return "404", "403", err
	} else if claims.Type != JwtRefreshType {
		return "404", "403", fmt.Errorf("invalid token type, expected refresh token")
	} else if err := database.C.Where(models.AuthTicket{
		BaseModel: models.BaseModel{ID: uint(parseInt(claims.SessionID))},
	}).First(&ticket).Error; err != nil {
		return "404", "403", err
	}

	if ticket, err := RegenSession(ticket); err != nil {
		return "404", "403", err
	} else {
		return GetToken(ticket)
	}
}
