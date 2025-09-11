package maps

import (
	"strings"
)

// Merge deep merges src into dst while ignoring empty keys and normalizing keys to lower-case.
func Merge(dst, src map[string]any) {
	for k, val := range src {
		normalK := strings.ToLower(strings.TrimSpace(k))
		if normalK == "" {
			continue
		}

		// If both dst[normalK] and val are maps, merge them recursively
		if dv, ok := dst[normalK]; ok {
			if dm, ok1 := dv.(map[string]any); ok1 {
				if sm, ok2 := val.(map[string]any); ok2 {
					Merge(dm, sm)

					continue
				}
			}
		}

		// Otherwise, just overwrite
		dst[normalK] = val
	}
}
