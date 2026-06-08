package moon

type posArg struct {
	name     string
	about    string
	setValue func(string) error
}

type varLenArg struct {
	name     string
	about    string
	addValue func(string) error
}

func (c *Command) AddStringPosArg(target *string, name string, about string, isRequired bool) {
	posArg := posArg{
		name:  name,
		about: about,
		setValue: func(s string) error {
			*target = s
			return nil
		},
	}

	if isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

func (c *Command) AddStringVarLenArg(target []string, name string, about string) {
	v := varLenArg{
		name:  name,
		about: about,
		addValue: func(s string) error {
			target = append(target, s)
			return nil
		},
	}

	c.varLenArg = &v
}