package dotenv

import (
	"bufio"
	"bytes"
	"strings"
)

// Parse parses dotenv-style configuration and returns a map of key->value.
// It supports quoted values and multi-line continuations inside quotes.
func Parse(data []byte) (map[string]string, error) {
	env := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(data))

	var (
		key          string
		valueBuilder strings.Builder
		inMultiline  bool
		quoteChar    rune
	)

	for scanner.Scan() {
		line := scanner.Text()

		// If we're in a multiline quoted value, keep appending
		if inMultiline {
			valueBuilder.WriteString("\n")
			valueBuilder.WriteString(line)

			// Check if this line ends the quoted value
			if strings.HasSuffix(strings.TrimRight(line, " \t"), string(quoteChar)) {
				value := strings.Trim(valueBuilder.String(), string(quoteChar))
				env[key] = value
				inMultiline = false

				valueBuilder.Reset()
			}

			continue
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key = strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// If value starts with quote but doesn’t close on same line → multiline
		if (strings.HasPrefix(value, `"`) && !strings.HasSuffix(strings.TrimRight(value, " \t"), `"`)) ||
			(strings.HasPrefix(value, `'`) && !strings.HasSuffix(strings.TrimRight(value, " \t"), `'`)) {
			inMultiline = true
			quoteChar = rune(value[0])
			valueBuilder.WriteString(strings.TrimPrefix(value, string(quoteChar)))

			continue
		}

		// Handle quoted values on one line
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			value = strings.Trim(value, `"`)
		} else if strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
			value = strings.Trim(value, `'`)
		}

		env[key] = value
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return env, nil
}
