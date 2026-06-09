package moon

type Command struct {
	Names      []string
	AboutShort string
	AboutLong  string
	Run        func() error

	subcommands     []*Command
	flags           []*Flag
	requiredPosArgs []*posArg
	optionalPosArgs []*posArg
	varLenArg       *varLenArg

	parent *Command
}

func (c *Command) AddSubcommand(s *Command) {
	s.parent = c
	c.subcommands = append(c.subcommands, s)
}