package env

import (
	"strings"
)

const (
	objSep = ":"
	envSep = "_"
)

func ParseVariables(vars map[string]string, pre, sep string, normalizeKey bool) map[string]any {
	data := make(map[string]any)

	for key, value := range vars {
		// Filter out unsafe variables
		if IsUnsafeVar(key) {
			continue
		}

		normalizedKey := key

		if pre != "" {
			if strings.HasPrefix(key, pre) {
				normalizedKey = strings.TrimPrefix(key, pre)
			} else {
				continue // Skip if doesn't match prefix
			}
		}

		// Convert to lowercase and replace separator with dots for nested structure
		normalizedKey = strings.ToLower(normalizedKey)
		normalizedKey = strings.ReplaceAll(normalizedKey, sep, objSep)

		if normalizeKey {
			// Convert "snake_case_key" to "snakecasekey", this can be accessed later as "snakeCaseKey" or "SnakeCaseKey"
			normalizedKey = strings.ReplaceAll(normalizedKey, envSep, "")
		}

		if sep != "" {
			// Build nested map structure
			BuildNestedMap(data, normalizedKey, value, objSep)
		} else {
			data[normalizedKey] = value
		}
	}

	return data
}
