// Example usage of the gcfg package
package main

import (
	"fmt"
	"testing/fstest"

	"github.com/ahmedkamalio/gcfg"
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

func main() {
	fsys := fstest.MapFS{
		"config.json": &fstest.MapFile{
			Data: []byte(`{
			  "database": {
				"host": "localhost",
				"port": 5432,
				"user": "admin",
				"password": "admin"
			  },
			  "server": {
				"host": "0.0.0.0",
				"port": 8080
			  },
			  "logging": {
				"level": "debug"
			  }
			}`),
		},
	}

	// initialize config instance
	config := gcfg.New(
		gcfg.NewJSONProvider(
			gcfg.WithJSONFilePath("config.json"),
			gcfg.WithJSONFileFS(&fsys),
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
