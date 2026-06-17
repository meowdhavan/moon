package moon

// Variable represents an entity that can contain a value.
// It is common base for [Flag], [PosArg], and [VarArgs].
type Variable struct {
	name       string
	aliases    []string
	about      string
	setValue   func(string) error
	env        *string
	defaultVal *string
	isRequired bool
	isValueSet bool
}

type variableProperty func(*Variable)

// Alias adds an alias name to the [Variable].
func Alias(alias string) variableProperty {
	return func(v *Variable) {
		v.aliases = append(v.aliases, alias)
	}
}

// Env specifies an environment variable as a fallback to the [Variable].
// If no flag is passed, the value in the provided
// environment variable will be used if present.
func Env(env string) variableProperty {
	return func(v *Variable) {
		v.env = &env
	}
}

// Default specifies a default value if the [Variable] is not provided.
// Due to current limitations, this value must be provided as a string.
func Default(defaultVal string) variableProperty {
	return func(v *Variable) {
		v.defaultVal = &defaultVal
	}
}

// Required marks the [Variable] as mandatory.
// If the value for this variable is not supplied and there
// are no fallback values, an error will be reported.
func Required() variableProperty {
	return func(v *Variable) {
		v.isRequired = true
	}
}
