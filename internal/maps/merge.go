package maps

import (
	"strings"
)

// Merge deep merges src into dst while ignoring empty keys and normalizing keys to lower-case.
func Merge(dst, src map[string]any) {
	for k, v := range src {
		nk := strings.ToLower(strings.TrimSpace(k))
		if nk == "" {
			continue
		}

		// If both dst[nk] and v are maps, merge them recursively
		if dv, ok := dst[nk]; ok {
			if dm, ok1 := dv.(map[string]any); ok1 {
				if sm, ok2 := v.(map[string]any); ok2 {
					Merge(dm, sm)

					continue
				}
			}
		}

		// Otherwise, just overwrite
		dst[nk] = v
	}
}
