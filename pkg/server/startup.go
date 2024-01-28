package server

import (
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var A *fiber.App

func NewServer() {
	A = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "Hydrogen.Passport",
		AppName:               "Hydrogen.Passport",
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
	})

	api := A.Group("/api").Name("API")
	{
		api.Get("/users/me", auth, getPrincipal)
		api.Post("/users", doRegister)

		api.Put("/auth", startChallenge)
		api.Post("/auth", doChallenge)
		api.Post("/auth/token", exchangeToken)
		api.Post("/auth/factors/:factorId", requestFactorToken)
	}
}

func Listen() {
	if err := A.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
