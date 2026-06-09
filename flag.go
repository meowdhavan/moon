package moon

import (
	"github.com/meowdhavan/moon/converter"
)

type Flag struct {
	Variable
	shortName   string
	requiresVal bool
}

func (c *Command) StringFlag(target *string, name string, shortName string, about string, options ...variableOption) {
	f := &Flag{
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
		shortName:   shortName,
		requiresVal: true,
	}

	for _, opt := range options {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) MultiStringFlag(target *[]string, name string, shortName string, about string, options ...variableOption) {
	*target = []string{}

	f := &Flag{
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
		shortName:   shortName,
		requiresVal: true,
	}

	for _, opt := range options {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) BoolFlag(target *bool, name string, shortName string, about string, options ...variableOption) {
	*target = false

	f := &Flag{
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
		shortName: shortName,
	}

	for _, opt := range options {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) MultiBoolFlag(target *int, name string, shortName string, about string, options ...variableOption) {
	*target = 0

	f := &Flag{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToBool(s)
				if err != nil {
					return err
				}

				if v {
					*target++
				}

				return nil
			},
		},
		shortName: shortName,
	}

	for _, opt := range options {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) IntFlag(target *int, name string, shortName string, about string, options ...variableOption) {
	f := &Flag{
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
		shortName:   shortName,
		requiresVal: true,
	}

	for _, opt := range options {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) MultiIntFlag(target *[]int, name string, shortName string, about string, options ...variableOption) {
	f := &Flag{
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
		shortName:   shortName,
		requiresVal: true,
	}

	for _, opt := range options {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}
