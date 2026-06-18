package moon

// Variable represents an entity that can contain a value.
// It is the common base for [Flag], [PosArg], and [VarArgs].
// It holds metadata such as name, aliases, description, and fallback values.
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
//
// Example:
//
//	cmd.Flags().String(&name, "name", "n", "Your name", moon.Alias("first-name"))
func Alias(alias string) variableProperty {
	return func(v *Variable) {
		v.aliases = append(v.aliases, alias)
	}
}

// Env specifies an environment variable as a fallback to the [Variable].
// If no flag is passed, the value in the provided
// environment variable will be used if present.
//
// Example:
//
//	cmd.Flags().String(&token, "token", "t", "API Token", moon.Env("API_TOKEN"))
func Env(env string) variableProperty {
	return func(v *Variable) {
		v.env = &env
	}
}

// Default specifies a default value if the [Variable] is not provided.
// Due to current limitations, this value must be provided as a string.
//
// Example:
//
//	cmd.Flags().Int(&port, "port", "p", "Port to listen on", moon.Default("8080"))
func Default(defaultVal string) variableProperty {
	return func(v *Variable) {
		v.defaultVal = &defaultVal
	}
}

// Required marks the [Variable] as mandatory.
// If the value for this variable is not supplied and there
// are no fallback values, an error will be reported.
//
// Example:
//
//	cmd.Flags().String(&config, "file", "f", "File to compile", moon.Required())
func Required() variableProperty {
	return func(v *Variable) {
		v.isRequired = true
	}
}
