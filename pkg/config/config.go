package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var requiredKeys []string

func RegisterRequiredKey(keys ...string) {
	for _, key := range keys {
		requiredKeys = append(requiredKeys, key)
	}
}

func ValidateRequiredKeys(logger zerolog.Logger, v *viper.Viper) {
	for _, key := range requiredKeys {
		if !v.IsSet(key) {
			logger.Panic().Str("config_key", key).Msg("missing required config key")
		}
	}
}

func NewConfig(logger zerolog.Logger, path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigName(path)

	err := conf.ReadInConfig()

	if err != nil {
		zerolog.DefaultContextLogger.Panic().Err(err).Msg("failed to load app config")
	}

	ValidateRequiredKeys(logger, conf)

	return conf
}
