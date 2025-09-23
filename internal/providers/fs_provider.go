// Package providers implements a base FS-based configuration provider.
package providers

import (
	"io/fs"

	"github.com/ahmedkamalio/gcfg/internal/sysfs"
)

// FSProvider provides file system operations by wrapping an fs.FS implementation.
// It is used as a base provider for other file-based configuration providers.
type FSProvider struct {
	fs fs.FS
}

// NewFSProvider creates a new FSProvider with the given fs.FS implementation.
// If fs is nil, it defaults to using sysfs.NewSysFS().
func NewFSProvider(fs fs.FS) *FSProvider {
	fsOrDefault := fs
	if fs == nil {
		fsOrDefault = sysfs.NewSysFS()
	}

	return &FSProvider{
		fs: fsOrDefault,
	}
}

// SetFS sets the underlying fs.FS implementation.
func (p *FSProvider) SetFS(fs fs.FS) {
	p.fs = fs
}

// OpenFile opens the named file using the underlying fs.FS implementation.
func (p *FSProvider) OpenFile(name string) (fs.File, error) {
	return p.fs.Open(name)
}

// ReadFile reads the named file using the underlying fs.FS implementation.
func (p *FSProvider) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(p.fs, name)
}
