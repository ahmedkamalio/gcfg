// Package gcfg provides a flexible configuration management system
// that supports reading from multiple providers and binding to user-defined types.
package gcfg

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ahmedkamalio/gcfg/internal/maps"
	"github.com/ahmedkamalio/gcfg/internal/reflection"
)

var (
	// ErrProviderLoadFailed indicates failure to load configuration from a provider.
	ErrProviderLoadFailed = errors.New("failed to load from provider")

	// ErrExtensionPreLoadHookFailed indicates a failure while executing the pre-load hook of an extension.
	ErrExtensionPreLoadHookFailed = errors.New("failed to execute extension pre-load hook")

	// ErrExtensionPostLoadHookFailed indicates a failure when executing the post-load hook of an extension.
	ErrExtensionPostLoadHookFailed = errors.New("failed to execute extension post-load hook")

	// ErrNilValues is returned when a nil value is provided where non-nil input is required.
	ErrNilValues = errors.New("values cannot be nil")
)

// Config represents the configuration loaded from various providers.
type Config struct {
	providers []Provider

	extensions []Extension

	values map[string]any
	mu     sync.RWMutex
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

// WithExtensions appends one or more extensions to the configuration and returns the updated Config instance.
func (c *Config) WithExtensions(extensions ...Extension) *Config {
	c.extensions = append(c.extensions, extensions...)

	return c
}

// SetDefault sets a default value for the specified key in the configuration.
// It creates nested maps if they do not exist.
func (c *Config) SetDefault(key string, value any) {
	if key == "" {
		return
	}

	pathParts, finalKey := keyToPathParts(key)

	c.mu.Lock()
	defer c.mu.Unlock()

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

	c.mu.Lock()
	defer c.mu.Unlock()

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

// Load loads configuration from all registered providers and applies pre/post-load hooks
// defined by extensions.
//
// Returns an error if any provider or extension hook fails during the loading process.
func (c *Config) Load() error {
	return c.LoadWithContext(context.Background())
}

// LoadWithContext loads configuration with the provided context, executing pre-load and post-load
// hooks for extensions.
func (c *Config) LoadWithContext(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ext := range c.extensions {
		if err := ext.PreLoad(ctx, c); err != nil {
			return fmt.Errorf("%w %s: %w", ErrExtensionPreLoadHookFailed, ext.Name(), err)
		}
	}

	for _, p := range c.providers {
		values, err := p.Load()
		if err != nil {
			return fmt.Errorf("%w %s: %w", ErrProviderLoadFailed, p.Name(), err)
		}
		// Merge values, later providers override
		maps.Merge(c.values, values)
	}

	for _, ext := range c.extensions {
		if err := ext.PostLoad(ctx, c); err != nil {
			return fmt.Errorf("%w %s: %w", ErrExtensionPostLoadHookFailed, ext.Name(), err)
		}
	}

	return nil
}

// Bind binds the configuration to the provided struct.
func (c *Config) Bind(dest any) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return maps.Bind(c.values, dest)
}

// Get retrieves a configuration value by key. Supports hierarchical paths like "database.host".
func (c *Config) Get(key string) any {
	if key == "" {
		return nil
	}

	pathParts, finalKey := keyToPathParts(key)

	c.mu.RLock()
	defer c.mu.RUnlock()

	finalMap := maps.FindNestedMap(c.values, pathParts, false)
	if finalMap != nil {
		return reflection.Clone(finalMap[finalKey])
	}

	return nil
}

// Values returns the configuration values.
func (c *Config) Values() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return reflection.Clone(c.values)
}

func keyToPathParts(key string) (pathParts []string, finalKey string) {
	parts := strings.Split(strings.ToLower(key), ".")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts[:len(parts)-1], parts[len(parts)-1]
}
