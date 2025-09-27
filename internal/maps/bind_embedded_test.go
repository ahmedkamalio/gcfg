package maps_test

import (
	"testing"

	"github.com/ahmedkamalio/gcfg/internal/maps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBind_EmbeddedStruct(t *testing.T) {
	t.Parallel()

	type Embedded struct {
		MyKey string `json:"mykey"`
		Value int    `json:"value"`
	}

	type Parent struct {
		Embedded

		Name string `json:"name"`
	}

	src := map[string]any{
		"mykey": "embedded_value",
		"value": 42,
		"name":  "parent_name",
	}

	var dest Parent

	err := maps.Bind(src, &dest)
	require.NoError(t, err)

	// These should work with embedded structs
	assert.Equal(t, "embedded_value", dest.MyKey)
	assert.Equal(t, 42, dest.Value)
	assert.Equal(t, "parent_name", dest.Name)
}

func TestBind_NestedEmbeddedStruct(t *testing.T) {
	t.Parallel()

	type Level1 struct {
		Field1 string `json:"field1"`
	}

	type Level2 struct {
		Level1

		Field2 string `json:"field2"`
	}

	type Parent struct {
		Level2

		ParentField string `json:"parent_field"`
	}

	src := map[string]any{
		"field1":       "level1_value",
		"field2":       "level2_value",
		"parent_field": "parent_value",
	}

	var dest Parent

	err := maps.Bind(src, &dest)
	require.NoError(t, err)

	assert.Equal(t, "level1_value", dest.Field1)
	assert.Equal(t, "level2_value", dest.Field2)
	assert.Equal(t, "parent_value", dest.ParentField)
}

func TestBind_EmbeddedStructWithPointer(t *testing.T) {
	t.Parallel()

	type Embedded struct {
		MyKey string `json:"mykey"`
	}

	type Parent struct {
		*Embedded

		Name string `json:"name"`
	}

	src := map[string]any{
		"mykey": "embedded_value",
		"name":  "parent_name",
	}

	var dest Parent

	err := maps.Bind(src, &dest)
	require.NoError(t, err)

	assert.Equal(t, "parent_name", dest.Name)
	// Embedded pointer should be initialized and populated
	require.NotNil(t, dest.Embedded)
	assert.Equal(t, "embedded_value", dest.MyKey)
}

func TestUnbind_EmbeddedStruct(t *testing.T) {
	t.Parallel()

	type Embedded struct {
		MyKey string `json:"mykey"`
		Value int    `json:"value"`
	}

	type Parent struct {
		Embedded

		Name string `json:"name"`
	}

	src := Parent{
		Embedded: Embedded{
			MyKey: "embedded_value",
			Value: 42,
		},
		Name: "parent_name",
	}

	dest := make(map[string]any)
	err := maps.Unbind(&src, dest)
	require.NoError(t, err)

	// All fields should be available at the top level
	assert.Equal(t, "embedded_value", dest["mykey"])
	assert.Equal(t, 42, dest["value"])
	assert.Equal(t, "parent_name", dest["name"])
}
