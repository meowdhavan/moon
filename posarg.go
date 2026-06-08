package moon

import "github.com/meowdhavan/moon/converter"

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
			v, err := converter.ToString(s)
			if err != nil {
				return err
			}

			*target = v
			return nil
		},
	}

	if isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

func (c *Command) AddBoolPosArg(target *bool, name string, about string, isRequired bool) {
	posArg := posArg{
		name:  name,
		about: about,
		setValue: func(s string) error {
			v, err := converter.ToBool(s)
			if err != nil {
				return err
			}

			*target = v
			return nil
		},
	}

	if isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

func (c *Command) AddIntPosArg(target *int, name string, about string, isRequired bool) {
	posArg := posArg{
		name:  name,
		about: about,
		setValue: func(s string) error {
			v, err := converter.ToInt(s)
			if err != nil {
				return err
			}

			*target = v
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
			v, err := converter.ToString(s)
			if err != nil {
				return err
			}

			target = append(target, v)
			return nil
		},
	}

	c.varLenArg = &v
}

func (c *Command) AddIntVarLenArg(target []int, name string, about string) {
	v := varLenArg{
		name:  name,
		about: about,
		addValue: func(s string) error {
			v, err := converter.ToInt(s)
			if err != nil {
				return err
			}

			target = append(target, v)
			return nil
		},
	}

	c.varLenArg = &v
}