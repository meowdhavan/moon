package converter

// ToString simply returns the string value as is, satisfying the converter pattern.
func ToString(s string) (string, error) {
	return s, nil
}