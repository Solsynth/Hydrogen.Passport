package services

import (
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type PayloadClaims struct {
	jwt.RegisteredClaims

	// Internal Stuff
	SessionID string `json:"sed"`

	// ID Token Stuff
	Name  string `json:"name,omitempty"`
	Nick  string `json:"preferred_username,omitempty"`
	Email string `json:"email,omitempty"`

	// Additional Stuff
	AuthorizedParties string `json:"azp,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
	Type              string `json:"typ"`
}

const (
	JwtAccessType  = "access"
	JwtRefreshType = "refresh"
)

func EncodeJwt(id string, typ, sub, sed string, nonce *string, aud []string, exp time.Time, idTokenUser ...models.Account) (string, error) {
	var azp string
	for _, item := range aud {
		if item != InternalTokenAudience {
			azp = item
			break
		}
	}

	claims := PayloadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Audience:  aud,
			Issuer:    viper.GetString("security.issuer"),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        id,
		},
		AuthorizedParties: azp,
		SessionID:         sed,
		Type:              typ,
	}

	if len(idTokenUser) > 0 {
		user := idTokenUser[0]
		claims.Name = user.Name
		claims.Nick = user.Nick
		claims.Email = user.GetPrimaryEmail().Content
	}

	if nonce != nil {
		claims.Nonce = *nonce
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

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
