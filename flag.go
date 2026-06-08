package moon

import "github.com/meowdhavan/moon/converter"

type flag struct {
	longNames   []string
	shortName   string
	about       string
	requiresVal bool
	isRequired  bool
	setValue    func(string) error
	isValueSet  bool
}

func (c *Command) AddStringFlag(target *string, longNames []string, shortName string, about string, isRequired bool) {
	c.flags = append(c.flags, flag{
		longNames: longNames,
		shortName: shortName,
		about: about,
		requiresVal: true,
		isRequired: isRequired,
		setValue: func(s string) error {
			v, err := converter.ToString(s)
			if err != nil {
				return err
			}

			*target = v
			return nil
		},
		isValueSet: false,
	})
}

func (c *Command) AddBoolFlag(target *bool, longNames []string, shortName string, about string, isRequired bool) {
	c.flags = append(c.flags, flag{
		longNames: longNames,
		shortName: shortName,
		about: about,
		requiresVal: false,
		isRequired: isRequired,
		setValue: func(s string) error {
			v, err := converter.ToBool(s)
			if err != nil {
				return err
			}

			*target = v
			return nil
		},
		isValueSet: false,
	})
}

func (c *Command) AddIntFlag(target *int, longNames []string, shortName string, about string, isRequired bool) {
	c.flags = append(c.flags, flag{
		longNames: longNames,
		shortName: shortName,
		about: about,
		requiresVal: true,
		isRequired: isRequired,
		setValue: func(s string) error {
			v, err := converter.ToInt(s)
			if err != nil {
				return err
			}

			*target = v
			return nil
		},
		isValueSet: false,
	})
}