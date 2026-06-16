package moon

type Command struct {
	Name       string
	Aliases    []string
	AboutShort string
	AboutLong  string
	Run        func()

	subcommands     []*Command
	localFlags      flagCollection
	globalFlags     flagCollection
	requiredPosArgs posArgCollection
	optionalPosArgs posArgCollection
	varArgs         varArgs

	parent *Command
}

func (c *Command) Flags() *flagCollection {
	return &c.localFlags
}

func (c *Command) GlobalFlags() *flagCollection {
	return &c.globalFlags
}

func (c *Command) RequiredPosArgs() *posArgCollection {
	return &c.requiredPosArgs
}

func (c *Command) OptionalPosArgs() *posArgCollection {
	return &c.optionalPosArgs
}

func (c *Command) VarArgs() *varArgs {
	return &c.varArgs
}

func (c *Command) Subcommand(s *Command) {
	s.parent = c
	c.subcommands = append(c.subcommands, s)
}
