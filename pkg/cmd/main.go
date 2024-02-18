package main

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/external"
	"code.smartsheep.studio/hydrogen/identity/pkg/server"
	"code.smartsheep.studio/hydrogen/identity/pkg/services"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"

	identity "code.smartsheep.studio/hydrogen/identity/pkg"
	"code.smartsheep.studio/hydrogen/identity/pkg/database"
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

	// Connect to database
	if err := database.NewSource(); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when connect to database.")
	} else if err := database.RunMigration(database.C); err != nil {
		log.Fatal().Err(err).Msg("An error occurred when running database auto migration.")
	}

	// External
	// All the things are optional so when error occurred the server won't crash
	if err := external.SetupFirebase(viper.GetString("external.firebase.credentials")); err != nil {
		log.Error().Err(err).Msg("An error occurred when starting firebase communicating...")
	} else {
		log.Info().Msg("Successfully setup firebase communication.")
	}

	// Server
	server.NewServer()
	go server.Listen()

	// Configure timed tasks
	quartz := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(&log.Logger)))
	quartz.AddFunc("@every 15s", func() {
		log.Info().Msg("Running auto sign off...")
		if tx := services.PerformAutoSignoff(); tx.Error != nil {
			log.Error().Err(tx.Error).Msg("An error occurred when running auto sign off...")
		} else {
			log.Info().Int64("affected", tx.RowsAffected).Msg("Auto sign off accomplished.")
		}
	})
	quartz.Run()

	// Messages
	log.Info().Msgf("Identity v%s is started...", identity.AppVersion)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("Identity v%s is quitting...", identity.AppVersion)

	quartz.Stop()
}
