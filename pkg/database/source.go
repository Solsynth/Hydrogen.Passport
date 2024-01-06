package database

import (
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var C *gorm.DB

func NewSource() error {
	var err error

	dialector := postgres.Open(viper.GetString("database.dsn"))
	C, err = gorm.Open(dialector, &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix: viper.GetString("database.prefix"),
	}, Logger: logger.New(&log.Logger, logger.Config{
		Colorful:                  true,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  lo.Ternary(viper.GetBool("debug"), logger.Info, logger.Silent),
	})})

	return err
}
