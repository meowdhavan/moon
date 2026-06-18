package converter

import "strings"

// ToBool converts a string representation to a boolean value.
func ToBool(s string) (bool, error) {
	if strings.ToLower(s) == "true" {
		return true, nil
	}

	return false, nil
}
