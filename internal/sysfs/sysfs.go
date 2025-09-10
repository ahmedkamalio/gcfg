package sysfs

import "io/fs"

type SysFS struct {
}

var _ fs.FS = (*SysFS)(nil)

func NewSysFS() *SysFS {
	return &SysFS{}
}

func (s SysFS) Open(name string) (fs.File, error) {
	return SafeOpen(name)
}
