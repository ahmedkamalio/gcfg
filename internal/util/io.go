//go:build !gcfg_unsafe

package util

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const maxConfigFileSize = 1 << 20 // 1 MB

// SafeReadFile ensures the file path is safe and reads the file content.
func SafeReadFile(filePath string) ([]byte, error) {
	// Get current working directory as baseDir
	baseDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Clean the input path
	cleanPath := filepath.Clean(filePath)

	// Make it absolute
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return nil, err
	}

	// Ensure the absolute path is within the baseDir
	if !strings.HasPrefix(absPath, baseDir+string(os.PathSeparator)) && absPath != baseDir {
		return nil, errors.New("unsafe file path: outside allowed directory")
	}

	// Ensure file is not a symlink
	info, err := os.Lstat(absPath)
	if err != nil {
		return nil, err
	}

	if info.Mode()&fs.ModeSymlink != 0 {
		return nil, errors.New("unsafe file path: symlink detected")
	}

	// Enforce size limit
	if info.Size() > maxConfigFileSize {
		return nil, errors.New("config file too large")
	}

	return os.ReadFile(absPath)
}
