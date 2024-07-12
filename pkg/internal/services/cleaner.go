package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/internal/database"
	"github.com/rs/zerolog/log"
)

func DoAutoDatabaseCleanup() {
	log.Debug().Msg("Now cleaning up entire database...")

	var count int64
	for _, model := range database.AutoMaintainRange {
		tx := database.C.Unscoped().Delete(model, "deleted_at IS NOT NULL")
		if tx.Error != nil {
			log.Error().Err(tx.Error).Msg("An error occurred when running cleaning up entire database...")
		}
		count += tx.RowsAffected
	}

	log.Debug().Int64("affected", count).Msg("Clean up entire database accomplished.")
}
