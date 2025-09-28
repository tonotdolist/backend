package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	const (
		key    = "key"
		badKey = "bkey"
		value  = "value"

		length = 1
	)
	RegisterRequiredKey(key)

	t.Run("SuccessfulConfig", func(t *testing.T) {
		const (
			data = key + ": " + value + "g\n"
		)

		conf, err := NewConfig(strings.NewReader(data), "yaml")

		require.NoError(t, err, "Unexpected error when initializing config.")
		assert.Equal(t, length, len(conf.AllKeys()), "The length of the config keys do not match.")
		assert.Equal(t, value, conf.Get(key), "The value of %s in the config does not match the expected value.", key)
	})

	t.Run("ValidateKeyFail", func(t *testing.T) {
		const (
			data = badKey + ": " + value + "\n"
		)

		_, err := NewConfig(strings.NewReader(data), "yaml")
		require.NoError(t, err, "Unexpected error when initializing config.")
		assert.NotNil(t, err, "Expected error caused by missing required key.")
	})

	t.Run("BadConfigFormat", func(t *testing.T) {
		const (
			data = key + " " + value + "\n"
		)

		_, err := NewConfig(strings.NewReader(data), "yaml")
		require.NoError(t, err, "Unexpected error when initializing config.")
		assert.NotNil(t, err, "Expected error caused by malformatted YAML config.")
	})
}
