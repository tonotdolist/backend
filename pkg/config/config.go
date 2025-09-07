package config

import (
	"fmt"
	"github.com/spf13/viper"
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

func NewConfig(path string) (*viper.Viper, error) {
	conf := viper.New()
	conf.SetConfigFile(path)

	err := conf.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("failed to load app config: %w", err)
	}

	err = ValidateRequiredKeys(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return conf, nil
}
