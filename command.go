package moon

// Command represents a command/subcommand in the application's CLI.
// It can hold [Flag], and either [PosArg], [VarArgs], or other [Command].
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

// Flags returns a collection to define local [Flag].
// Local flags are only available on this specific command.
func (c *Command) Flags() *flagCollection {
	return &c.localFlags
}

// GlobalFlags returns a collection to define global [Flag].
// Global flags are available on this command and all of its subcommands.
func (c *Command) GlobalFlags() *flagCollection {
	return &c.globalFlags
}

// PosArgs returns a collection to define [PosArg] for this command.
func (c *Command) PosArgs() *posArgCollection {
	return &c.posArgs
}

// VarArgs returns a collection to define [VarArgs] for this command.
func (c *Command) VarArgs() *varArgs {
	return &c.varArgs
}

// Subcommand adds a subcommand to this command.
func (c *Command) Subcommand(s *Command) {
	s.parent = c
	c.subcommands = append(c.subcommands, s)
}
