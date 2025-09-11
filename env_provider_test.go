package gcfg_test

import (
	"testing"

	"github.com/go-gase/gcfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvProvider_DefaultOptions(t *testing.T) {
	t.Setenv("TEST_KEY", "test_value")

	p := gcfg.NewEnvProvider()

	values, err := p.Load()
	require.NoError(t, err)
	assert.Equal(t, "test_value", values["testkey"])
}

func TestEnvProvider_WithEnvPrefix(t *testing.T) {
	t.Setenv("TEST_KEY", "unaccessible_value")
	t.Setenv("MYAPP_TEST_KEY", "test_value")

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvPrefix("MYAPP_"), // load only prefixed variables
	)

	values, err := p.Load()
	require.NoError(t, err)
	assert.Equal(t, "test_value", values["testkey"])
}

func TestEnvProvider_WithEnvSeparator(t *testing.T) {
	t.Setenv("TEST__KEY", "test_value")

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvSeparator("__"),
	)

	values, err := p.Load()
	require.NoError(t, err)
	assert.IsType(t, map[string]any{}, values["test"])
	//nolint:forcetypeassert
	assert.Equal(t, "test_value", values["test"].(map[string]any)["key"])
}

func TestEnvProvider_WithEnvNormalizeVarNames(t *testing.T) {
	t.Setenv("TEST_KEY", "test_value")

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvNormalizeVarNames(false), // keep original variable names
	)

	values, err := p.Load()
	require.NoError(t, err)
	assert.Equal(t, "test_value", values["test_key"])
}
