package gcfg

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/go-gase/gcfg/internal/dotenv"
	"github.com/go-gase/gcfg/internal/env"
	"github.com/go-gase/gcfg/internal/providers"
)

const (
	defaultDotEnvFilePath = ".env"

	dotenvProviderName = "DotEnv"
)

// DotEnvProvider reads configuration from .env file.
type DotEnvProvider struct {
	*providers.FSProvider
	*EnvProvider

	filePath string
}

var _ Provider = (*DotEnvProvider)(nil)

// DotEnvOption is a function that configures a DotEnvProvider.
type DotEnvOption func(*DotEnvProvider)

// WithDotEnvFilePath sets the .env file path.
func WithDotEnvFilePath(filePath string) DotEnvOption {
	return func(p *DotEnvProvider) {
		p.filePath = filePath
	}
}

// WithDotEnvSeparator sets the separator for nested map values.
// Given a sep=__ variables like DATABASE__URL become database.url in the resulting map.
func WithDotEnvSeparator(sep string) DotEnvOption {
	return func(p *DotEnvProvider) {
		p.separator = sep
	}
}

// WithDotEnvNormalizeVarNames sets a flag to normalize variable names.
// If set to true, all variable names are converted from snake_case to lowercase identifier (snake case without underscores).
// This is useful to access environment variable names like "DATABASE_URL" with the key "DatabaseUrl".
//
// Default: true.
func WithDotEnvNormalizeVarNames(normalized bool) DotEnvOption {
	return func(p *DotEnvProvider) {
		p.normalizeVarNames = normalized
	}
}

// WithDotEnvFileFS sets the fs of which to read the .env file from.
//
// Default: sysfs.SysFS.
func WithDotEnvFileFS(fs fs.FS) DotEnvOption {
	return func(p *DotEnvProvider) {
		p.SetFS(fs)
	}
}

// NewDotEnvProvider creates .env provider with options.
func NewDotEnvProvider(opts ...DotEnvOption) *DotEnvProvider {
	p := &DotEnvProvider{
		FSProvider:  providers.NewFSProvider(nil),
		EnvProvider: NewEnvProvider(),
		filePath:    defaultDotEnvFilePath,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Load implements the Provider interface.
func (p *DotEnvProvider) Load() (map[string]any, error) {
	if p.filePath == "" {
		return nil, errors.New(".env file path is not set")
	}

	file, err := p.ReadFile(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file %s: %w", p.filePath, err)
	}

	vars, err := dotenv.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse .env file %s: %w", p.filePath, err)
	}

	return env.ParseVariables(vars, p.prefix, p.separator, p.normalizeVarNames), nil
}

func (p *DotEnvProvider) Name() string {
	return dotenvProviderName
}
