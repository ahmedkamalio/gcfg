package gcfg_test

import (
	"testing"

	"github.com/go-gase/gcfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Bind(t *testing.T) {
	t.Setenv("MY_KEY", "my_value")

	cfg := gcfg.New(gcfg.NewEnvProvider())

	require.NoError(t, cfg.Load())

	obj := struct {
		MyKey string
	}{}

	require.NoError(t, cfg.Bind(&obj))

	assert.Equal(t, "my_value", obj.MyKey)
}
