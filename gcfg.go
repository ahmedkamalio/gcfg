// Package gcfg provides a flexible configuration management system
// that supports reading from multiple providers and binding to user-defined types.
package gcfg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gase/gcfg/internal/maps"
)

var (
	// ErrProviderLoadFailed indicates failure to load configuration from a provider.
	ErrProviderLoadFailed = errors.New("failed to load from provider")

	// ErrNilValues is returned when a nil value is provided where non-nil input is required.
	ErrNilValues = errors.New("values cannot be nil")
)

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

// SetDefault sets a default value for the specified key in the configuration.
// It creates nested maps if they do not exist.
func (c *Config) SetDefault(key string, value any) {
	if key == "" {
		return
	}

	pathParts, finalKey := c.keyToPathParts(key)

	finalMap := maps.FindNestedMap(c.values, pathParts, true)
	if finalMap != nil {
		finalMap[finalKey] = value
	}
}

// SetDefaults sets default configuration values from a struct or map. Returns an error if the input is invalid or nil.
func (c *Config) SetDefaults(values any) error {
	if values == nil {
		return ErrNilValues
	}

	if val, ok := values.(map[string]any); ok {
		maps.Merge(c.values, val)

		return nil
	}

	if val, ok := values.(*map[string]any); ok {
		maps.Merge(c.values, *val)

		return nil
	}

	if err := maps.Unbind(values, c.values); err != nil {
		return err
	}

	maps.LowercaseKeys(c.values)

	return nil
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

// Get retrieves a configuration value by key. Supports hierarchical paths like "database.host".
func (c *Config) Get(key string) any {
	if key == "" {
		return nil
	}

	pathParts, finalKey := c.keyToPathParts(key)

	finalMap := maps.FindNestedMap(c.values, pathParts, false)
	if finalMap != nil {
		return finalMap[finalKey]
	}

	return nil
}

// Values returns the configuration values.
func (c *Config) Values() map[string]any {
	return c.values
}

func (c *Config) keyToPathParts(key string) (pathParts []string, finalKey string) {
	parts := strings.Split(strings.ToLower(key), ".")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts[:len(parts)-1], parts[len(parts)-1]
}
