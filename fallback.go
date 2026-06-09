package moon

import "os"

func getFromEnv(v *Variable) *string {
	if v.env == nil {
		return nil
	}

	val := os.Getenv(*v.env)
	if val == "" {
		return nil
	}

	return &val
}

func getDefault(v *Variable) *string {
	if v.defaultVal == nil {
		return nil
	}

	return v.defaultVal
}
