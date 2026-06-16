package moon

import "github.com/meowdhavan/moon/converter"

type PosArg struct {
	Variable
}

type posArgCollection struct {
	requiredPosArgs []*PosArg
	optionalPosArgs []*PosArg
}

func (c *posArgCollection) String(target *string, name string, about string, properties ...variableProperty) {
	posArg := &PosArg{
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

	for _, opt := range properties {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

func (c *posArgCollection) Bool(target *bool, name string, about string, properties ...variableProperty) {
	*target = false

	posArg := &PosArg{
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

	for _, opt := range properties {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

func (c *posArgCollection) Int(target *int, name string, about string, properties ...variableProperty) {
	posArg := &PosArg{
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

	for _, opt := range properties {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}
