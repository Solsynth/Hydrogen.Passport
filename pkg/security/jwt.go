package security

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var CookieAccessKey = "identity_auth_key"
var CookieRefreshKey = "identity_refresh_key"

type PayloadClaims struct {
	jwt.RegisteredClaims

	SessionID string `json:"sed"`
	Type      string `json:"typ"`
}

const (
	JwtAccessType  = "access"
	JwtRefreshType = "refresh"
)

func EncodeJwt(id string, typ, sub, sed string, aud []string, exp time.Time) (string, error) {
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, PayloadClaims{
		jwt.RegisteredClaims{
			Subject:   sub,
			Audience:  aud,
			Issuer:    fmt.Sprintf("https://%s", viper.GetString("domain")),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        id,
		},
		sed,
		typ,
	})

	return tk.SignedString([]byte(viper.GetString("secret")))
}

func DecodeJwt(str string) (PayloadClaims, error) {
	var claims PayloadClaims
	tk, err := jwt.ParseWithClaims(str, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("secret")), nil
	})
	if err != nil {
		return claims, err
	}

	if data, ok := tk.Claims.(*PayloadClaims); ok {
		return *data, nil
	} else {
		return claims, fmt.Errorf("unexpected token payload: not payload claims type")
	}
}

func SetJwtCookieSet(c *fiber.Ctx, access, refresh string) {
	c.Cookie(&fiber.Cookie{
		Name:     CookieAccessKey,
		Value:    access,
		Domain:   viper.GetString("security.cookie_domain"),
		SameSite: viper.GetString("security.cookie_samesite"),
		Expires:  time.Now().Add(60 * time.Minute),
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     CookieRefreshKey,
		Value:    refresh,
		Domain:   viper.GetString("security.cookie_domain"),
		SameSite: viper.GetString("security.cookie_samesite"),
		Expires:  time.Now().Add(24 * 30 * time.Hour),
		Path:     "/",
	})
}
