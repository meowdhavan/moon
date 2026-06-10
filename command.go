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
	requiredPosArgs []*PosArg
	optionalPosArgs []*PosArg
	varLenArg       *VarLenArg

	parent *Command
}

func (c *Command) GlobalFlags() *flagCollection {
	return &c.globalFlags
}

func (c *Command) Flags() *flagCollection {
	return &c.localFlags
}

func (c *Command) Subcommand(s *Command) {
	s.parent = c
	c.subcommands = append(c.subcommands, s)
}
