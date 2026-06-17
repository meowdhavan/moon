package moon

import (
	"github.com/meowdhavan/moon/converter"
)

// Flag represents a command-line flag (e.g., --xyz or -x) with its parsed value and properties.
type Flag struct {
	Variable
	shortName   string
	requiresVal bool
}

type flagCollection struct {
	flags []*Flag
}

// String adds a [Flag] of type string.
// It sets the target value on encountering the
// flag when parsing the command-line arguments.
func (c *flagCollection) String(target *string, name string, shortName string, about string, properties ...variableProperty) {
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

	for _, opt := range properties {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

// MultiString adds a [Flag] of type string.
// It appends a value to the target slice on encountering
// the flag when parsing the command-line arguments.
func (c *flagCollection) MultiString(target *[]string, name string, shortName string, about string, properties ...variableProperty) {
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

	for _, opt := range properties {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

// Bool adds a [Flag] of type bool.
// It sets the target value to true on encountering
// the flag when parsing the command-line arguments.
func (c *flagCollection) Bool(target *bool, name string, shortName string, about string, properties ...variableProperty) {
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

	for _, opt := range properties {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

// MultiBool adds a [Flag] of type bool.
// It increments the target int by 1 on encountering
// the flag when parsing the command-line arguments.
func (c *flagCollection) MultiBool(target *int, name string, shortName string, about string, properties ...variableProperty) {
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

	for _, opt := range properties {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

// Int adds a [Flag] of type int.
// It sets the target value on encountering the
// flag when parsing the command-line arguments.
func (c *flagCollection) Int(target *int, name string, shortName string, about string, properties ...variableProperty) {
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

	for _, opt := range properties {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}

// MultiInt adds a [Flag] of type int.
// It appends a value to the target slice on encountering
// the flag when parsing the command-line arguments.
func (c *flagCollection) MultiInt(target *[]int, name string, shortName string, about string, properties ...variableProperty) {
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

	for _, opt := range properties {
		opt(&f.Variable)
	}

	c.flags = append(c.flags, f)
}
