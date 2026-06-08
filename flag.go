package moon

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
			*target = s
			return nil
		},
		isValueSet: false,
	})
}
