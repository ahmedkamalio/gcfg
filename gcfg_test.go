package gcfg_test

import (
	"os"
	"testing"

	"github.com/go-gase/gcfg"
)

func TestConfig_Bind(t *testing.T) {
	if err := os.Setenv("MY_KEY", "my_value"); err != nil {
		t.Error(err)
	}

	cfg := gcfg.New(gcfg.NewEnvProvider())

	if err := cfg.Load(); err != nil {
		t.Error(err)
	}

	obj := struct {
		MyKey string
	}{}

	if err := cfg.Bind(&obj); err != nil {
		t.Error(err)
	}

	if obj.MyKey != "my_value" {
		t.Error("expected my_value, got", obj.MyKey)
	}
}
