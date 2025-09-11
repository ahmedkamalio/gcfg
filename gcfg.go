// Package gcfg provides a flexible configuration management system
// that supports reading from multiple providers and binding to user-defined types.
package gcfg

import (
	"errors"
	"fmt"

	"github.com/go-gase/gcfg/internal/maps"
)

// ErrProviderLoadFailed indicates failure to load configuration from a provider.
var ErrProviderLoadFailed = errors.New("failed to load from provider")

// Config represents the configuration loaded from various providers.
type Config struct {
	values    map[string]any
	providers []Provider
}

// New creates a new config instance with given providers.
func New(providers ...Provider) *Config {
	pvd := append([]Provider{}, providers...)

	hasEnvProvider := false

	for _, p := range providers {
		if p.Name() == envProviderName {
			hasEnvProvider = true

			break
		}
	}

	if !hasEnvProvider {
		pvd = append([]Provider{NewEnvProvider()}, pvd...)
	}

	return &Config{
		values:    make(map[string]any),
		providers: pvd,
	}
}

// Provider defines the interface for configuration providers.
// Implement this interface to create custom providers like env, json, yml, etc.
type Provider interface {
	Name() string
	// Load reads configuration from the source and returns it as a map.
	// Keys should be hierarchical paths (e.g., "database.host").
	Load() (map[string]any, error)
}

// Load merges configuration from all providers.
// Later providers override earlier ones.
func (c *Config) Load() error {
	for _, p := range c.providers {
		values, err := p.Load()
		if err != nil {
			return fmt.Errorf("%w %s: %w", ErrProviderLoadFailed, p.Name(), err)
		}
		// Merge values, later providers override
		maps.Merge(c.values, values)
	}

	return nil
}

// Bind binds the configuration to the provided struct.
func (c *Config) Bind(dest any) error {
	return maps.Bind(c.values, dest)
}

// Get retrieves a configuration value by key.
func (c *Config) Get(key string) any {
	return c.values[key]
}

// Values returns the configuration values.
func (c *Config) Values() map[string]any {
	return c.values
}
