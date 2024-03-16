package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func getMetadata(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":              viper.GetString("name"),
		"domain":            viper.GetString("domain"),
		"open_registration": !viper.GetBool("use_registration_magic_token"),
	})
}

func getOidcConfiguration(c *fiber.Ctx) error {
	domain := viper.GetString("domain")
	basepath := fmt.Sprintf("https://%s", domain)

	return c.JSON(fiber.Map{
		"issuer":                                           basepath,
		"authorization_endpoint":                           fmt.Sprintf("%s/auth/o/connect", basepath),
		"token_endpoint":                                   fmt.Sprintf("%s/api/auth/token", basepath),
		"userinfo_endpoint":                                fmt.Sprintf("%s/api/users/me", basepath),
		"response_types_supported":                         []string{"code", "token"},
		"grant_types_supported":                            []string{"authorization_code", "implicit", "refresh_token"},
		"subject_types_supported":                          []string{"public"},
		"token_endpoint_auth_methods_supported":            []string{"client_secret_post"},
		"id_token_signing_alg_values_supported":            []string{"HS512"},
		"token_endpoint_auth_signing_alg_values_supported": []string{"HS512"},
	})
}
