# gcfg

A flexible configuration management system for Go that supports reading from multiple providers and binding to
user-defined types.

## Features

- **Multiple Providers**: Support for environment variables, JSON files, dotenv files, and more
- **Flexible API**: Easy-to-use interface for loading and binding configuration
- **Extensible**: Create custom providers by implementing the `Provider` interface
- **Hierarchical**: Support for nested configuration using dot notation (e.g., `database.host`)
- **Merge Strategy**: Later providers override earlier ones for flexible configuration layering

## Installation

```bash
go get github.com/ahmedkamalio/gcfg
```

## Quick Start

### Basic Usage

```go
package main

import (
	"fmt"
	"github.com/ahmedkamalio/gcfg"
)

type AppConfig struct {
	Database struct {
		Host string
		Port int
	}
	Server struct {
		Host string
		Port int
	}
}

func main() {
	// Create config (environment provider is registered by default).
	config := gcfg.New()

	// Load configuration
	if err := config.Load(); err != nil {
		panic(err)
	}

	// Bind to your config struct
	var appCfg AppConfig
	if err := config.Bind(&appCfg); err != nil {
		panic(err)
	}

	fmt.Printf("Server: %s:%d\n", appCfg.Server.Host, appCfg.Server.Port)
}
```

### Using JSON Configuration

```go
package main

import (
	"github.com/ahmedkamalio/gcfg"
)

func main() {
	// Initialize config with JSON provider
	config := gcfg.New(
		gcfg.NewJSONProvider(
			gcfg.WithJSONFilePath("config.json"),
		),
	)

	// Load and bind as shown above
	// ...
}

```

Example `config.json`:

```json
{
  "database": {
    "host": "localhost",
    "port": 5432
  },
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  }
}
```

### Using environment variables

Set environment variables:

```bash
export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
```

Then in your Go code:

```go
config := gcfg.New(gcfg.NewEnvProvider())
```

### Using .env files

```go
package main

import (
	"github.com/ahmedkamalio/gcfg"
)

func main() {
	// Initialize config with dotenv provider
	config := gcfg.New(
		gcfg.NewDotEnvProvider(), // defaults to ".env"
	)

	// Load and bind as shown above
	// ...
}
```

Example `.env` file:

```env
DATABASE_HOST=localhost
DATABASE_PORT=5432
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

### Combining Multiple Providers

You can combine multiple providers, with later providers overriding earlier ones:

```go
config := gcfg.New(
// Default values from JSON
gcfg.NewJSONProvider(gcfg.WithJSONFilePath("config.json")),
// Override with environment variables
gcfg.NewEnvProvider(),
// Override with .env file
gcfg.NewDotEnvProvider(),
)
```

## API Reference

### Config

#### `New(providers ...Provider) *Config`

Creates a new configuration instance with the given providers. If no `EnvProvider` is provided, one will be added
automatically.

#### `SetDefault(key string, value any)`

Sets a default value for the specified key in the configuration. Supports hierarchical paths like "database.host"

#### `SetDefaults(values any) error`

Sets default configuration values from a struct or map. Returns an error if the input is invalid or nil.

#### `Load() error`

Loads configuration from all providers, merging values. Later providers override earlier ones.

#### `Bind(dest any) error`

Binds the loaded configuration to a Go struct using reflection.

#### `Get(key string) any`

Retrieves a configuration value by key (supports hierarchical paths like "database.host").

#### `Values() map[string]any`

Returns all configuration values as a map.

### Providers

#### `Provider` interface

```go
type Provider interface {
    Name() string
    Load() (map[string]any, error)
}
```

#### Built-in Providers

- `NewEnvProvider()` - Loads from environment variables
- `NewJSONProvider(options ...JSONProviderOption)` - Loads from JSON files
- `NewDotEnvProvider()` - Loads from dotenv files

#### Custom Providers

```go
type CustomProvider struct{}

func (c *CustomProvider) Name() string {
    return "custom"
}

func (c *CustomProvider) Load() (map[string]any, error) {
    // Load configuration from your custom source
    return map[string]any{
        "my.setting": "value",
    }, nil
}

// Use it
config := gcfg.New(&CustomProvider{})
```

## Examples

See the [`examples/`](./examples/) directory for complete examples:

- [`basic/`](./examples/basic/) - Basic usage
- [`json/`](./examples/json/) - JSON configuration
- [`dotenv/`](./examples/dotenv/) - Dotenv files
- [`json_fs/`](./examples/json_fs/) - JSON with filesystem provider

## Error Handling

All operations that can fail return errors that should be handled appropriately:

```go
config := gcfg.New(gcfg.NewJSONProvider(gcfg.WithJSONFilePath("config.json")))
if err := config.Load(); err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

var appCfg AppConfig
if err := config.Bind(&appCfg); err != nil {
    log.Fatalf("Failed to bind config: %v", err)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
