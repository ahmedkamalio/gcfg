//go:build gcfg_unsafe

package util

import (
	"os"
)

// SafeReadFile an unsafe version of the original SafeReadFile.
func SafeReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
