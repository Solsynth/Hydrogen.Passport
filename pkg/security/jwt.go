package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type PayloadClaims struct {
	jwt.RegisteredClaims

	Type  string `json:"typ"`
	Value any    `json:"val"`
}

const (
	JwtAccessType  = "access"
	JwtRefreshType = "refresh"
)

func EncodeJwt(id string, val any, typ, sub string, aud []string, exp time.Time) (string, error) {
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
		typ,
		val,
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
