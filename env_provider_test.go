package gcfg_test

import (
	"os"
	"testing"

	"github.com/go-gase/gcfg"
	"github.com/stretchr/testify/assert"
)

func TestEnvProvider_DefaultOptions(t *testing.T) {
	assert.NoError(t, os.Setenv("TEST_KEY", "test_value"))
	defer func() {
		_ = os.Unsetenv("TEST_KEY")
	}()

	p := gcfg.NewEnvProvider()

	values, err := p.Load()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", values["testkey"])
}

func TestEnvProvider_WithEnvPrefix(t *testing.T) {
	assert.NoError(t, os.Setenv("TEST_KEY", "unaccessible_value"))
	assert.NoError(t, os.Setenv("MYAPP_TEST_KEY", "test_value"))
	defer func() {
		_ = os.Unsetenv("TEST_KEY")
		_ = os.Unsetenv("MYAPP_TEST_KEY")
	}()

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvPrefix("MYAPP_"), // load only prefixed variables
	)

	values, err := p.Load()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", values["testkey"])
}

func TestEnvProvider_WithEnvSeparator(t *testing.T) {
	assert.NoError(t, os.Setenv("TEST__KEY", "test_value"))
	defer func() {
		_ = os.Unsetenv("TEST__KEY")
	}()

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvSeparator("__"),
	)

	values, err := p.Load()
	assert.NoError(t, err)
	assert.IsType(t, map[string]any{}, values["test"])
	assert.Equal(t, "test_value", values["test"].(map[string]any)["key"])
}

func TestEnvProvider_WithEnvNormalizeVarNames(t *testing.T) {
	assert.NoError(t, os.Setenv("TEST_KEY", "test_value"))
	defer func() {
		_ = os.Unsetenv("TEST_KEY")
	}()

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvNormalizeVarNames(false), // keep original variable names
	)

	values, err := p.Load()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", values["test_key"])
}
