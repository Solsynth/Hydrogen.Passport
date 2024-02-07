package main

import (
	"code.smartsheep.studio/hydrogen/passport/pkg/external"
	"code.smartsheep.studio/hydrogen/passport/pkg/server"
	"os"
	"os/signal"
	"syscall"

	passport "code.smartsheep.studio/hydrogen/passport/pkg"
	"code.smartsheep.studio/hydrogen/passport/pkg/database"
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
	}

	// Server
	server.NewServer()
	go server.Listen()

	// Messages
	log.Info().Msgf("Passport v%s is started...", passport.AppVersion)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("Passport v%s is quitting...", passport.AppVersion)
}
