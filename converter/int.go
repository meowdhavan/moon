package converter

import "strconv"

// ToInt converts a string representation to an integer value.
func ToInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	return val, err
}
