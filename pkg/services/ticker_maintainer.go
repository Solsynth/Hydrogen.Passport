package services

import (
	"git.solsynth.dev/hydrogen/passport/pkg/database"
	"git.solsynth.dev/hydrogen/passport/pkg/models"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.etcd.io/bbolt"
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

func DoAutoAuthCleanup() {
	log.Debug().Msg("Now cleaning up cached auth context...")

	count := 0
	err := database.B.Batch(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(authContextBucket))
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()

		var ctx models.AuthContext
		for key, val := cursor.First(); key != nil; key, val = cursor.Next() {
			if err := jsoniter.Unmarshal(val, &ctx); err != nil {
				bucket.Delete(key)
				count++
			} else if time.Now().Unix() >= ctx.ExpiredAt.Unix() {
				bucket.Delete(key)
				count++
			}
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("An error occurred when running auth context cleanup...")
	} else {
		log.Debug().Int("affected", count).Msg("Clean up auth context accomplished.")
	}
}
