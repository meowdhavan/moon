package moon

type Command struct {
	Name         string
	Version      string
	Aliases      []string
	AboutShort   string
	AboutLong    string
	Run          func()
	SuppressHelp bool

	subcommands []*Command
	localFlags  flagCollection
	globalFlags flagCollection
	posArgs     posArgCollection
	varArgs     varArgs

	parent *Command
}

func (c *Command) Flags() *flagCollection {
	return &c.localFlags
}

func (c *Command) GlobalFlags() *flagCollection {
	return &c.globalFlags
}

func (c *Command) PosArgs() *posArgCollection {
	return &c.posArgs
}

func (c *Command) VarArgs() *varArgs {
	return &c.varArgs
}

func (c *Command) Subcommand(s *Command) {
	s.parent = c
	c.subcommands = append(c.subcommands, s)
}
