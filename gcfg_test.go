package gcfg_test

import (
	"errors"
	"testing"

	"github.com/go-gase/gcfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockProvider is a mock implementation of the Provider interface for testing.
type mockProvider struct {
	name string
	data map[string]any
	err  error
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Load() (map[string]any, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.data, nil
}

func TestConfig_WithProviders(t *testing.T) {
	t.Parallel()

	mockP1 := &mockProvider{name: "mock1", data: map[string]any{"mock1key": "m1value"}}
	mockP2 := &mockProvider{name: "mock2", data: map[string]any{"mock2key": "m2value"}}
	cfg := gcfg.New(mockP1, mockP2)

	require.NoError(t, cfg.Load())

	// Check that data from mocks is loaded
	assert.Equal(t, "m1value", cfg.Get("mock1key"))
	assert.Equal(t, "m2value", cfg.Get("mock2key"))
}

func TestConfig_WithEnvProvider_AlreadyPresent(t *testing.T) {
	t.Parallel()

	mockP1 := &mockProvider{name: "env", data: map[string]any{"customEnv": "custom"}}
	mockP2 := &mockProvider{name: "mock2", data: map[string]any{"mockKey": "mockValue"}}
	cfg := gcfg.New(mockP1, mockP2)

	err := cfg.Load()
	require.NoError(t, err)

	assert.Equal(t, "custom", cfg.Get("customEnv"))
	assert.Equal(t, "mockValue", cfg.Get("mockKey"))
}

func TestConfig_NoProviders(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	require.NoError(t, cfg.Load())
	// NOTE: actual env variables may be present in Values()
	assert.NotNil(t, cfg.Values())
}

func TestConfig_LoadProviderError(t *testing.T) {
	t.Parallel()

	//nolint:err113
	mockP1 := &mockProvider{name: "mock1", err: errors.New("load failed")}
	cfg := gcfg.New(mockP1)

	err := cfg.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load from provider mock1: load failed")
	assert.ErrorIs(t, err, gcfg.ErrProviderLoadFailed)
}

func TestConfig_Bind(t *testing.T) {
	t.Parallel()

	mockP1 := &mockProvider{name: "mock", data: map[string]any{
		"myKey": "value",
	}}

	cfg := gcfg.New(mockP1)

	require.NoError(t, cfg.Load())

	obj := struct {
		MyKey string
	}{}

	require.NoError(t, cfg.Bind(&obj))

	assert.Equal(t, "value", obj.MyKey)
}

func TestConfig_BindError(t *testing.T) {
	t.Parallel()

	cfg := &gcfg.Config{}
	obj := "not a struct"
	err := cfg.Bind(obj)
	assert.Error(t, err)
}

func TestConfig_Get(t *testing.T) {
	t.Parallel()

	mockP1 := &mockProvider{name: "mock", data: map[string]any{
		"simple": "value",
		"nested": map[string]any{
			"key": "nestedvalue",
		},
	}}
	cfg := gcfg.New(mockP1)

	err := cfg.Load()
	require.NoError(t, err)

	assert.Equal(t, "value", cfg.Get("simple"))
	assert.Equal(t, "nestedvalue", cfg.Get("nested.key"))
	assert.Nil(t, cfg.Get("nonexistent"))
	assert.Nil(t, cfg.Get("nested.nonexistent"))
	assert.Nil(t, cfg.Get("")) // empty key
}

func TestConfig_Values(t *testing.T) {
	t.Parallel()

	mockP1 := &mockProvider{name: "mock", data: map[string]any{"key": "value"}}
	cfg := gcfg.New(mockP1)

	err := cfg.Load()
	require.NoError(t, err)

	values := cfg.Values()
	assert.NotNil(t, values)
	// Since env vars are loaded, check specific key
	assert.Equal(t, "value", cfg.Get("key"))
}

func TestConfig_SetDefault_Basic(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	cfg.SetDefault("key", "value")

	assert.Equal(t, "value", cfg.Get("key"))
}

func TestConfig_SetDefault_Nested(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	cfg.SetDefault("database.host", "localhost")

	assert.Equal(t, "localhost", cfg.Get("database.host"))

	// Verify nested structure
	values := cfg.Values()
	db, ok := values["database"]
	assert.True(t, ok)
	dbMap, ok := db.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "localhost", dbMap["host"])
}

func TestConfig_SetDefault_MultipleLevels(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	cfg.SetDefault("a.b.c.d", "deepvalue")

	assert.Equal(t, "deepvalue", cfg.Get("a.b.c.d"))

	values := cfg.Values()
	a, ok := values["a"].(map[string]any)
	require.True(t, ok)
	b, ok := a["b"].(map[string]any)
	require.True(t, ok)
	c, ok := b["c"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "deepvalue", c["d"])
}

func TestConfig_SetDefault_EmptyKey(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	cfg.SetDefault("", "value")

	values := cfg.Values()
	// Empty key should be ignored, no value should be set
	assert.NotNil(t, values)
	// But since it's ignored, no specific key
}

func TestConfig_SetDefault_Override(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	cfg.SetDefault("key", "oldvalue")
	assert.Equal(t, "oldvalue", cfg.Get("key"))

	cfg.SetDefault("key", "newvalue")
	assert.Equal(t, "newvalue", cfg.Get("key"))
}

func TestConfig_SetDefault_WithSpaces(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	cfg.SetDefault("key.with spaces", "value")

	assert.Equal(t, "value", cfg.Get("key.with spaces"))
}

func TestConfig_SetDefaults_ValidStruct(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	s := struct {
		Key   string
		Value int
	}{
		Key:   "testkey",
		Value: 42,
	}

	err := cfg.SetDefaults(&s)
	require.NoError(t, err)

	assert.Equal(t, "testkey", cfg.Get("key"))
	assert.Equal(t, 42, cfg.Get("value"))
}

func TestConfig_SetDefaults_StructWithoutJSONTags(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	s := struct {
		Key   string
		Value int
	}{
		Key:   "test",
		Value: 123,
	}

	err := cfg.SetDefaults(&s)
	require.NoError(t, err)

	assert.Equal(t, "test", cfg.Get("key"))
	assert.Equal(t, 123, cfg.Get("value"))
}

func TestConfig_SetDefaults_Map(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	m := map[string]any{
		"key1": "value1",
		"key2": map[string]any{
			"nested": "nestedvalue",
		},
	}

	err := cfg.SetDefaults(&m)
	require.NoError(t, err)

	assert.Equal(t, "value1", cfg.Get("key1"))
	assert.Equal(t, "nestedvalue", cfg.Get("key2.nested"))
}

func TestConfig_SetDefaults_NilValues(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	err := cfg.SetDefaults(nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, gcfg.ErrNilValues)
}

func TestConfig_SetDefaults_InvalidType(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	err := cfg.SetDefaults(42) // not a pointer
	require.Error(t, err)
}

func TestConfig_SetDefaults_NilPointer(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	var s *struct{}

	err := cfg.SetDefaults(s) // nil pointer
	require.Error(t, err)
}

func TestConfig_SetDefaults_PointerToNonStructNonMap(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	x := 42

	err := cfg.SetDefaults(&x)
	require.Error(t, err)
}

func TestConfig_SetDefaults_EmptyStruct(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	s := struct{}{}

	err := cfg.SetDefaults(&s)
	require.NoError(t, err)

	values := cfg.Values()
	assert.NotNil(t, values)
	// Empty struct should not add any keys
}

func TestConfig_SetDefaults_Merge(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	s1 := struct {
		Key1 string `json:"key1"`
	}{Key1: "value1"}

	s2 := struct {
		Key2 string `json:"key2"`
	}{Key2: "value2"}

	err := cfg.SetDefaults(&s1)
	require.NoError(t, err)
	err = cfg.SetDefaults(&s2)
	require.NoError(t, err)

	assert.Equal(t, "value1", cfg.Get("key1"))
	assert.Equal(t, "value2", cfg.Get("key2"))
}

func TestConfig_SetDefaults_NestedStruct(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	s := struct {
		Database struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		} `json:"database"`
	}{}

	s.Database.Host = "localhost"
	s.Database.Port = 5432

	err := cfg.SetDefaults(&s)
	require.NoError(t, err)

	assert.Equal(t, "localhost", cfg.Get("database.host"))
	assert.Equal(t, 5432, cfg.Get("database.port"))
}

func TestConfig_GetAfterSetDefaults(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	m := map[string]any{
		"app": map[string]any{
			"name":    "myapp",
			"version": "1.0.0",
		},
	}

	err := cfg.SetDefaults(&m)
	require.NoError(t, err)

	assert.Equal(t, "myapp", cfg.Get("app.name"))
	assert.Equal(t, "1.0.0", cfg.Get("app.version"))
}

func TestConfig_SetDefaultAndSetDefaultsInteraction(t *testing.T) {
	t.Parallel()

	cfg := gcfg.New()

	cfg.SetDefault("key", "setdefault")
	assert.Equal(t, "setdefault", cfg.Get("key"))

	s := struct {
		Key string `json:"key"`
	}{Key: "setdefaults"}

	err := cfg.SetDefaults(&s)
	require.NoError(t, err)

	// SetDefaults should merge/override
	assert.Equal(t, "setdefaults", cfg.Get("key"))
}
