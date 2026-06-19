package converter

import (
	"errors"
	"strconv"
)

// ToInt converts a string representation to an integer value.
func ToInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("Value cannot be converted into an integer: " + s)
	}

	return val, nil
}
