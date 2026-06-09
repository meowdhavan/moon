package moon

import "github.com/meowdhavan/moon/converter"

type posArg struct {
	Variable
}

type varLenArg struct {
	name     string
	aliases  []string
	about    string
	addValue func(string) error
}

func (c *Command) AddStringPosArg(target *string, name string, about string, options ...variableOption) {
	posArg := &posArg{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToString(s)
				if err != nil {
					return err
				}

				*target = v
				return nil
			},
		},
	}

	for _, opt := range options {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

func (c *Command) AddStringVarLenArg(target *[]string, name string, about string) {
	*target = []string{}

	v := &varLenArg{
		name:    name,
		aliases: []string{},
		about:   about,
		addValue: func(s string) error {
			v, err := converter.ToString(s)
			if err != nil {
				return err
			}

			*target = append(*target, v)
			return nil
		},
	}

	c.varLenArg = v
}

func (c *Command) AddBoolPosArg(target *bool, name string, about string, options ...variableOption) {
	*target = false

	posArg := &posArg{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToBool(s)
				if err != nil {
					return err
				}

				*target = v
				return nil
			},
		},
	}

	for _, opt := range options {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

func (c *Command) AddIntPosArg(target *int, name string, about string, options ...variableOption) {
	posArg := &posArg{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToInt(s)
				if err != nil {
					return err
				}

				*target = v
				return nil
			},
		},
	}

	for _, opt := range options {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

func (c *Command) AddIntVarLenArg(target *[]int, name string, about string) {
	*target = []int{}

	v := &varLenArg{
		name:    name,
		aliases: []string{},
		about:   about,
		addValue: func(s string) error {
			v, err := converter.ToInt(s)
			if err != nil {
				return err
			}

			*target = append(*target, v)
			return nil
		},
	}

	c.varLenArg = v
}
