package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

func DoAutoSignoff() {
	duration := time.Duration(viper.GetInt64("security.auto_signoff_duration")) * time.Second
	divider := time.Now().Add(-duration)

	log.Debug().Time("before", divider).Msg("Now signing off tickets...")

	if tx := database.C.
		Where("last_grant_at < ?", divider).
		Delete(&models.AuthTicket{}); tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred when running auto sign off...")
	} else {
		log.Debug().Int64("affected", tx.RowsAffected).Msg("Auto sign off accomplished.")
	}
}
