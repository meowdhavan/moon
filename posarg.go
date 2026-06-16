package moon

import "github.com/meowdhavan/moon/converter"

type PosArg struct {
	Variable
}

type posArgCollection struct {
	posArgs []*PosArg
}

func (c *posArgCollection) String(target *string, name string, about string, options ...variableOption) {
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

	for _, opt := range options {
		opt(&posArg.Variable)
	}

	c.posArgs = append(c.posArgs, posArg)
}

func (c *posArgCollection) Bool(target *bool, name string, about string, options ...variableOption) {
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

	for _, opt := range options {
		opt(&posArg.Variable)
	}

	c.posArgs = append(c.posArgs, posArg)
}

func (c *posArgCollection) Int(target *int, name string, about string, options ...variableOption) {
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

	for _, opt := range options {
		opt(&posArg.Variable)
	}

	c.posArgs = append(c.posArgs, posArg)
}
