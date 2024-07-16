package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/admin"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/api"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type HTTPApp struct {
	app *fiber.App
}

func NewServer() *HTTPApp {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "Hydrogen.Passport",
		AppName:               "Hydrogen.Passport",
		ProxyHeader:           fiber.HeaderXForwardedFor,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
	})

	app.Use(idempotency.New())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodOptions,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	app.Use(logger.New(logger.Config{
		Format: "${status} | ${latency} | ${method} ${path}\n",
		Output: log.Logger,
	}))

	app.Use(exts.AuthMiddleware)

	admin.MapAdminAPIs(app)
	api.MapAPIs(app)

	app.Use(filesystem.New(filesystem.Config{
		Root:         http.Dir(viper.GetString("frontend_app")),
		Index:        "index.html",
		NotFoundFile: "index.html",
		MaxAge:       3600,
	}))

	app.Use(favicon.New(favicon.Config{
		File: filepath.Join(viper.GetString("frontend_app"), "favicon.png"),
		URL:  "/favicon.png",
	}))

	return &HTTPApp{app}
}

func (v *HTTPApp) Listen() {
	if err := v.app.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
