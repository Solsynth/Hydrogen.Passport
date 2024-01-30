package server

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/view"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
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
		EnablePrintRoutes:     viper.GetBool("debug"),
	})

	A.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(view.FS),
		PathPrefix:   "dist",
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	A.Get("/.well-known", getMetadata)
	A.Get("/.well-known/openid-configuration", getOidcConfiguration)

	api := A.Group("/api").Name("API")
	{
		api.Get("/users/me", auth, getPrincipal)
		api.Get("/users/me/events", auth, getEvents)
		api.Delete("/users/me/sessions/:sessionId", auth, killSession)

		api.Post("/users", doRegister)
		api.Post("/users/me/confirm", doRegisterConfirm)

		api.Put("/auth", startChallenge)
		api.Post("/auth", doChallenge)
		api.Post("/auth/token", exchangeToken)
		api.Post("/auth/factors/:factorId", requestFactorToken)

		api.Get("/auth/oauth/connect", auth, preConnect)
		api.Post("/auth/oauth/connect", auth, doConnect)
	}
}

func Listen() {
	if err := A.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
