package security

import (
	"fmt"

	"code.smartsheep.studio/hydrogen/passport/pkg/models"
	"github.com/samber/lo"
)

func GetFactorCode(factor models.AuthFactor) (bool, error) {
	switch factor.Type {
	case models.EmailPasswordFactor:
		// TODO
		return true, nil
	default:
		return false, nil
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
