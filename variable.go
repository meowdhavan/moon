package moon

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

func Alias(alias string) variableProperty {
	return func(v *Variable) {
		v.aliases = append(v.aliases, alias)
	}
}

func Env(env string) variableProperty {
	return func(v *Variable) {
		v.env = &env
	}
}

func Default(defaultVal string) variableProperty {
	return func(v *Variable) {
		v.defaultVal = &defaultVal
	}
}

func Required() variableProperty {
	return func(v *Variable) {
		v.isRequired = true
	}
}
