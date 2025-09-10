package gcfg

import (
	"encoding/json"
	"fmt"

	"github.com/go-gase/gcfg/internal/util"
)

const (
	jsonProviderName = "JSON"
)

// JSONProvider reads configuration from a JSON file.
type JSONProvider struct {
	filePath string
}

var _ Provider = (*JSONProvider)(nil)

// JSONOption is a function that configures a JSONProvider.
type JSONOption func(*JSONProvider)

// WithJSONFilePath sets the JSON file path.
func WithJSONFilePath(filePath string) JSONOption {
	return func(p *JSONProvider) {
		p.filePath = filePath
	}
}

// NewJSONProvider creates a new file provider.
func NewJSONProvider(opts ...JSONOption) *JSONProvider {
	p := &JSONProvider{}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Load implements the Provider interface.
func (p *JSONProvider) Load() (map[string]any, error) {
	file, err := util.SafeReadFile(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON config file %s: %w", p.filePath, err)
	}

	var data map[string]any
	if err = json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON from %s: %w", p.filePath, err)
	}

	return data, nil
}

func (p *JSONProvider) Name() string {
	return jsonProviderName
}
