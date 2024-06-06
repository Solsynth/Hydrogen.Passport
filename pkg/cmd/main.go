package main

import (
	"os"
	"os/signal"
	"syscall"

	"git.solsynth.dev/hydrogen/passport/pkg"
	"git.solsynth.dev/hydrogen/passport/pkg/i18n"

	"git.solsynth.dev/hydrogen/passport/pkg/grpc"
	"git.solsynth.dev/hydrogen/passport/pkg/server"
	"git.solsynth.dev/hydrogen/passport/pkg/services"
	"github.com/robfig/cron/v3"

	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	// Configure settings
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")

	// Load settings
	if err := viper.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("An error occurred when loading settings.")
	}

	i18n.InitInternationalization()

	// Connect to database
	if err := database.NewGorm(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connect to database.")
	} else if err := database.RunMigration(database.C); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when running database auto migration.")
	}

	// Connect other services
	if err := services.SetupFirebase(); err != nil {
		log.Error().Err(err).Msg("An error occurred when connecting firebase...")
	}
	if err := grpc.ConnectPaperclip(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connecting to paperclip...")
	}

	// Server
	server.NewServer()
	go server.Listen()

	// Grpc Server
	go func() {
		if err := grpc.StartGrpc(); err != nil {
			log.Fatal().Err(err).Msg("An message occurred when starting grpc server.")
		}
	}()

	// Configure timed tasks
	quartz := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)))
	quartz.AddFunc("@every 60m", services.DoAutoSignoff)
	quartz.AddFunc("@every 60m", services.DoAutoDatabaseCleanup)
	quartz.AddFunc("@every 60s", services.RecycleAuthContext)
	quartz.AddFunc("@every 5m", services.KexCleanup)
	quartz.Start()

	// Messages
	log.Info().Msgf("Passport v%s is started...", pkg.AppVersion)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("Passport v%s is quitting...", pkg.AppVersion)

	quartz.Stop()
}
