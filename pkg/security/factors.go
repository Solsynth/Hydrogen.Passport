package security

import (
	"fmt"

	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"github.com/samber/lo"
)

func GetFactorCode(factor models.AuthFactor) error {
	switch factor.Type {
	case models.EmailPasswordFactor:
		// TODO
		return nil
	default:
		return fmt.Errorf("unsupported factor to get code")
	}
}

func VerifyFactor(factor models.AuthFactor, code string) error {
	switch factor.Type {
	case models.PasswordAuthFactor:
		return lo.Ternary(
			VerifyPassword(code, factor.Secret),
			nil,
			fmt.Errorf("invalid password"),
		)
	}

	return nil
}
