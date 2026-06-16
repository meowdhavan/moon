package moon

import "github.com/meowdhavan/moon/converter"

type PosArg struct {
	Variable
}

type posArgCollection struct {
	posArgs []*PosArg
}

type VarLenArg struct {
	Variable
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

func (c *Command) StringVarLenArg(target *[]string, name string, about string, options ...variableOption) {
	*target = []string{}

	v := &VarLenArg{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToString(s)
				if err != nil {
					return err
				}

				*target = append(*target, v)
				return nil
			},
		},
	}

	for _, opt := range options {
		opt(&v.Variable)
	}

	c.varLenArg = v
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

func (c *Command) IntVarLenArg(target *[]int, name string, about string, options ...variableOption) {
	*target = []int{}

	v := &VarLenArg{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToInt(s)
				if err != nil {
					return err
				}

				*target = append(*target, v)
				return nil
			},
		},
	}

	for _, opt := range options {
		opt(&v.Variable)
	}

	c.varLenArg = v
}
