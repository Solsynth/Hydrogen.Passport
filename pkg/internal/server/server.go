package server

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/api"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/exts"
	"net/http"
	"strings"

	"git.solsynth.dev/hydrogen/passport/pkg/internal"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/i18n"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/admin"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/server/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var A *fiber.App

func NewServer() {
	templates := html.NewFileSystem(http.FS(pkg.FS), ".gohtml")

	A = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		EnableIPValidation:    true,
		ServerHeader:          "Hydrogen.Passport",
		AppName:               "Hydrogen.Passport",
		ProxyHeader:           fiber.HeaderXForwardedFor,
		JSONEncoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
		JSONDecoder:           jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
		EnablePrintRoutes:     viper.GetBool("debug.print_routes"),
		Views:                 templates,
		ViewsLayout:           "views/index",
	})

	A.Use(idempotency.New())
	A.Use(cors.New(cors.Config{
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

	A.Use(logger.New(logger.Config{
		Format: "${status} | ${latency} | ${method} ${path}\n",
		Output: log.Logger,
	}))

	A.Use(exts.AuthMiddleware)
	A.Use(i18n.I18nMiddleware)

	A.Use(favicon.New(favicon.Config{
		FileSystem: http.FS(pkg.FS),
		File:       "views/favicon.png",
		URL:        "/favicon.png",
	}))

	api.MapAPIs(A)
	admin.MapAdminEndpoints(A)
	ui.MapUserInterface(A)
}

func Listen() {
	if err := A.Listen(viper.GetString("bind")); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when starting server...")
	}
}
