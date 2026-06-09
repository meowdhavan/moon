package moon

type Command struct {
	Name       string
	Aliases    []string
	AboutShort string
	AboutLong  string
	Run        func() error

	subcommands     []*Command
	flags           []*Flag
	requiredPosArgs []*PosArg
	optionalPosArgs []*PosArg
	varLenArg       *VarLenArg

	parent *Command
}

func (c *Command) Subcommand(s *Command) {
	s.parent = c
	c.subcommands = append(c.subcommands, s)
}
