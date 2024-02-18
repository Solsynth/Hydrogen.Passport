package security

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/models"
	"fmt"
	"github.com/samber/lo"
)

func VerifyFactor(factor models.AuthFactor, code string) error {
	switch factor.Type {
	case models.PasswordAuthFactor:
		return lo.Ternary(
			VerifyPassword(code, factor.Secret),
			nil,
			fmt.Errorf("invalid password"),
		)
	case models.EmailPasswordFactor:
		return lo.Ternary(
			code == factor.Secret,
			nil,
			fmt.Errorf("invalid verification code"),
		)
	}

	return nil
}
