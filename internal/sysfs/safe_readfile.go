package sysfs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrUnsafeFilePathOutsideDirectory indicates the file path is outside the allowed directory.
	ErrUnsafeFilePathOutsideDirectory = errors.New("unsafe file path: outside allowed directory")
	// ErrUnsafeFilePathSymlink indicates the file path is a symlink which is not allowed.
	ErrUnsafeFilePathSymlink = errors.New("unsafe file path: symlink detected")
	// ErrConfigFileTooLarge indicates the config file exceeds the maximum allowed size.
	ErrConfigFileTooLarge = errors.New("config file too large")
)

const maxConfigFileSize = 1 << 20 // 1 MB

// SafeOpen ensures the file path is safe and opens the file.
func SafeOpen(filePath string) (*os.File, error) {
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
		return nil, ErrUnsafeFilePathOutsideDirectory
	}

	// Ensure file is not a symlink
	info, err := os.Lstat(absPath)
	if err != nil {
		return nil, err
	}

	if info.Mode()&fs.ModeSymlink != 0 {
		return nil, ErrUnsafeFilePathSymlink
	}

	// Enforce size limit
	if info.Size() > maxConfigFileSize {
		return nil, ErrConfigFileTooLarge
	}

	//nolint:gosec
	return os.Open(absPath)
}
