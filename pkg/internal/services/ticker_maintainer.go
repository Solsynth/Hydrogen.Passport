package services

import (
	"time"

	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"git.solsynth.dev/hydrogen/passport/pkg/internal/models"
	"github.com/rs/zerolog/log"
)

func DoAutoSignoff() {
	duration := 7 * 24 * time.Hour
	deadline := time.Now().Add(-duration)

	log.Debug().Time("before", deadline).Msg("Now signing off tickets...")

	if tx := database.C.
		Where("last_grant_at < ?", deadline).
		Delete(&models.AuthTicket{}); tx.Error != nil {
		log.Error().Err(tx.Error).Msg("An error occurred when running auto sign off...")
	} else {
		log.Debug().Int64("affected", tx.RowsAffected).Msg("Auto sign off accomplished.")
	}
}
