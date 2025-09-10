package gcfg_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-gase/gcfg"
)

func TestEnvProvider_DefaultOptions(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "test_value"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Unsetenv("TEST_KEY")
	}()

	p := gcfg.NewEnvProvider()

	values, err := p.Load()
	if err != nil {
		t.Fatal(err)
	}

	if v := values["testkey"]; v != "test_value" {
		t.Errorf("expected 'test_value', got %v", v)
	}
}

func TestEnvProvider_WithEnvPrefix(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "unaccessible_value"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("MYAPP_TEST_KEY", "test_value"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Unsetenv("TEST_KEY")
		_ = os.Unsetenv("MYAPP_TEST_KEY")
	}()

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvPrefix("MYAPP_"), // load only prefixed variables
	)

	values, err := p.Load()
	if err != nil {
		t.Fatal(err)
	}

	if v := values["testkey"]; v != "test_value" {
		t.Errorf("expected 'test_value', got %v", v)
	}
}

func TestEnvProvider_WithEnvSeparator(t *testing.T) {
	if err := os.Setenv("TEST__KEY", "test_value"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Unsetenv("TEST__KEY")
	}()

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvSeparator("__"),
	)

	values, err := p.Load()
	if err != nil {
		t.Fatal(err)
	}

	tvalues, ok := values["test"].(map[string]any)
	if !ok {
		t.Errorf("values.test to be a map, got %v", reflect.TypeOf(values["test"]))
	}

	if tvalues["key"] != "test_value" {
		t.Errorf("expected 'test_value', got %v", values["TEST_KEY"])
	}
}

func TestEnvProvider_WithEnvNormalizeVarNames(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "test_value"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Unsetenv("TEST_KEY")
	}()

	p := gcfg.NewEnvProvider(
		gcfg.WithEnvNormalizeVarNames(false), // keep original variable names
	)

	values, err := p.Load()
	if err != nil {
		t.Fatal(err)
	}

	if v := values["test_key"]; v != "test_value" {
		t.Errorf("expected 'test_value', got %v", v)
	}
}
