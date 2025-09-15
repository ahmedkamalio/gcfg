package gcfg_test

import (
	"testing"
	"testing/fstest"

	"github.com/go-gase/gcfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDotEnvProvider_DefaultOptions(t *testing.T) {
	t.Parallel()

	p := gcfg.NewDotEnvProvider()
	_, err := p.Load()
	assert.Error(t, err)
}

func TestDotEnvProvider_WithDotEnvFile_FileNotFound(t *testing.T) {
	t.Parallel()

	p := gcfg.NewDotEnvProvider(
		gcfg.WithDotEnvFilePath(".env.non-existing"),
	)
	_, err := p.Load()
	assert.Error(t, err)
}

func TestDotEnvProvider_WithDotEnvFile(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		".env": &fstest.MapFile{
			Data: []byte(`
				TEST_KEY=test_value
			`),
		},
	}

	p := gcfg.NewDotEnvProvider(
		gcfg.WithDotEnvFilePath(".env"),
		gcfg.WithDotEnvFileFS(&fsys),
	)

	values, err := p.Load()
	require.NoError(t, err)
	assert.Equal(t, "test_value", values["testkey"])
}

func TestDotEnvProvider_WithEnvSeparator(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		".env": &fstest.MapFile{
			Data: []byte(`
				TEST__KEY=test_value
			`),
		},
	}

	p := gcfg.NewDotEnvProvider(
		gcfg.WithDotEnvFilePath(".env"),
		gcfg.WithDotEnvFileFS(&fsys),
		gcfg.WithDotEnvSeparator("__"),
	)

	values, err := p.Load()
	require.NoError(t, err)
	assert.IsType(t, map[string]any{}, values["test"])
	assert.Equal(t, "test_value", values["test"].(map[string]any)["key"])
}

func TestDotEnvProvider_WithEnvNormalizeVarNames(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		".env": &fstest.MapFile{
			Data: []byte(`
				TEST_KEY=test_value
			`),
		},
	}

	p := gcfg.NewDotEnvProvider(
		gcfg.WithDotEnvFilePath(".env"),
		gcfg.WithDotEnvFileFS(&fsys),
		gcfg.WithDotEnvNormalizeVarNames(false), // keep original variable names
	)

	values, err := p.Load()
	require.NoError(t, err)
	assert.Equal(t, "test_value", values["test_key"])
}

func TestDotEnvProvider_Syntax(t *testing.T) {
	t.Parallel()

	fsys := fstest.MapFS{
		".env": &fstest.MapFile{
			Data: []byte(`
				# This is a comment
				TEST_KEY=test_value
				TEST_KEY2=test_value2 # This is an inline comment
			`),
		},
	}

	p := gcfg.NewDotEnvProvider(
		gcfg.WithDotEnvFilePath(".env"),
		gcfg.WithDotEnvFileFS(&fsys),
	)

	values, err := p.Load()
	require.NoError(t, err)
	assert.Equal(t, "test_value", values["testkey"])
	assert.Equal(t, "test_value2", values["testkey2"])
}
