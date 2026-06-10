package moon

import (
	"github.com/meowdhavan/moon/converter"
)

type Flag struct {
	Variable
	shortName   string
	requiresVal bool
}

type flagCollection struct {
	flags []*Flag
}

func (c *flagCollection) StringFlag(target *string, name string, shortName string, about string, options ...variableOption) {
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

func (c *flagCollection) MultiStringFlag(target *[]string, name string, shortName string, about string, options ...variableOption) {
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

func (c *flagCollection) BoolFlag(target *bool, name string, shortName string, about string, options ...variableOption) {
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

func (c *flagCollection) MultiBoolFlag(target *int, name string, shortName string, about string, options ...variableOption) {
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

func (c *flagCollection) IntFlag(target *int, name string, shortName string, about string, options ...variableOption) {
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

func (c *flagCollection) MultiIntFlag(target *[]int, name string, shortName string, about string, options ...variableOption) {
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
