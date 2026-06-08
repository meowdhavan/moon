package converter

import "strconv"

func ToInt(s string) (int, error) {
	val, err := strconv.Atoi(s)
	return val, err
}