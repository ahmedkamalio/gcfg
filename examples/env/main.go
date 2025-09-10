// Example usage of the gcfg package
package main

import (
	"fmt"
	"os"

	"github.com/go-gase/gcfg"
)

type AppConfig struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
	}
	Server struct {
		Host string
		Port int
	}
	Logging struct {
		Level string
	}
}

func initEnvVars() {
	_ = os.Setenv("MYAPP_DATABASE__HOST", "localhost")
	_ = os.Setenv("MYAPP_DATABASE__PORT", "5432")
	_ = os.Setenv("MYAPP_DATABASE__USER", "admin")
	_ = os.Setenv("MYAPP_DATABASE__PASSWORD", "admin")
	_ = os.Setenv("MYAPP_SERVER__HOST", "0.0.0.0")
	_ = os.Setenv("MYAPP_SERVER__PORT", "8080")
	_ = os.Setenv("MYAPP_LOGGING__LEVEL", "debug")
}

func main() {
	initEnvVars()

	// initialize config instance
	config := gcfg.New(
		gcfg.NewEnvProvider(
			gcfg.WithEnvPrefix("MYAPP_"),
		),
	)

	// Load configuration
	if err := config.Load(); err != nil {
		panic(err)
	}

	// Bind to user-defined type
	var appCfg AppConfig
	if err := config.Bind(&appCfg); err != nil {
		panic(err)
	}

	// Use the config
	fmt.Printf("Server: %s:%d\n", appCfg.Server.Host, appCfg.Server.Port)
	fmt.Printf(
		"DB: postgresql://%s:%s@%s:%d\n",
		appCfg.Database.User,
		appCfg.Database.Password,
		appCfg.Database.Host,
		appCfg.Database.Port,
	)
	fmt.Printf("Log Level: %s\n", appCfg.Logging.Level)
}
