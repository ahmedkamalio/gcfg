package providers

import (
	"io/fs"

	"github.com/go-gase/gcfg/internal/sysfs"
)

type FSProvider struct {
	fs fs.FS
}

func NewFSProvider(fs fs.FS) *FSProvider {
	fsOrDefault := fs
	if fs == nil {
		fsOrDefault = sysfs.NewSysFS()
	}
	return &FSProvider{
		fs: fsOrDefault,
	}
}

func (p *FSProvider) SetFS(fs fs.FS) {
	p.fs = fs
}

func (p *FSProvider) OpenFile(name string) (fs.File, error) {
	return p.fs.Open(name)
}

func (p *FSProvider) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(p.fs, name)
}
