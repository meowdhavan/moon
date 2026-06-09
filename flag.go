package moon

import (
	"github.com/meowdhavan/moon/converter"
)

type flag struct {
	longNames   []string
	shortName   string
	about       string
	requiresVal bool
	setValue    func(string) error
	isValueSet  bool
	env         *string
	defaultVal  *string
	isRequired  bool
}

type flagOption func(*flag)

func Alias(longName string) flagOption {
	return func(f *flag) {
		f.longNames = append(f.longNames, longName)
	}
}

func About(about string) flagOption {
	return func(f *flag) {
		f.about = about
	}
}

func Env(env string) flagOption {
	return func(f *flag) {
		*f.env = env
	}
}

func Default(defaultVal string) flagOption {
	return func(f *flag) {
		f.defaultVal = &defaultVal
	}
}

func Required() flagOption {
	return func(f *flag) {
		f.isRequired = true
	}
}

func (c *Command) StringFlag(target *string, longName string, shortName string, options ...flagOption) {
	f := &flag{
		longNames: []string{longName},
		shortName: shortName,
		requiresVal: true,
		setValue: func(s string) error {
			v, err := converter.ToString(s)
			if err != nil {
				return err
			}

			*target = v
			return nil
		},
	}

	for _, opt := range options {
		opt(f)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) MultiStringFlag(target *[]string, longName string, shortName string, options ...flagOption) {
	*target = []string{}

	f := &flag{
		longNames:   []string{longName},
		shortName:   shortName,
		requiresVal: true,
		setValue: func(s string) error {
			v, err := converter.ToString(s)
			if err != nil {
				return err
			}

			*target = append(*target, v)
			return nil
		},
	}

	for _, opt := range options {
		opt(f)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) BoolFlag(target *bool, longName string, shortName string, options ...flagOption) {
	*target = false

	f := &flag{
		longNames:   []string{longName},
		shortName:   shortName,
		setValue: func(s string) error {
			v, err := converter.ToBool(s)
			if err != nil {
				return err
			}

			*target = v
			return nil
		},
	}

	for _, opt := range options {
		opt(f)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) MultiBoolFlag(target *int, longName string, shortName string, options ...flagOption) {
	*target = 0

	f := &flag{
		longNames:   []string{longName},
		shortName:   shortName,
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
	}

	for _, opt := range options {
		opt(f)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) IntFlag(target *int, longName string, shortName string, options ...flagOption) {
	f := &flag{
		longNames:   []string{longName},
		shortName:   shortName,
		requiresVal: true,
		setValue: func(s string) error {
			v, err := converter.ToInt(s)
			if err != nil {
				return err
			}

			*target = v
			return nil
		},
	}

	for _, opt := range options {
		opt(f)
	}

	c.flags = append(c.flags, f)
}

func (c *Command) MultiIntFlag(target *[]int, longName string, shortName string, options ...flagOption) {
	f := &flag{
		longNames:   []string{longName},
		shortName:   shortName,
		requiresVal: true,
		setValue: func(s string) error {
			v, err := converter.ToInt(s)
			if err != nil {
				return err
			}

			*target = append(*target, v)
			return nil
		},
	}

	for _, opt := range options {
		opt(f)
	}

	c.flags = append(c.flags, f)
}
