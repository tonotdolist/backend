package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
)

var requiredKeys []string

func RegisterRequiredKey(keys ...string) {
	for _, key := range keys {
		requiredKeys = append(requiredKeys, key)
	}
}

func ValidateRequiredKeys(v *viper.Viper) error {
	for _, key := range requiredKeys {
		if !v.IsSet(key) {
			return fmt.Errorf("missing required config key: %s", key)
		}
	}

	return nil
}

func NewConfig(reader io.Reader, configType string) (*viper.Viper, error) {
	conf := viper.New()

	conf.SetConfigType(configType)
	err := conf.ReadConfig(reader)

	if err != nil {
		return nil, fmt.Errorf("failed to load app config: %w", err)
	}

	err = ValidateRequiredKeys(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return conf, nil
}
