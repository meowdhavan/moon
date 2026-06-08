package converter

import "strings"

func ToBool(s string) (bool, error) {
	if s == "" || strings.ToLower(s) == "true" {
		return true, nil
	}

	return false, nil
}