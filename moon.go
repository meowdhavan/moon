package moon

import "os"

type Command struct {
	Names      []string
	AboutShort string
	AboutLong  string
	Run        func() error

	subcommands     map[string]*Command
	flags           []flag
	requiredPosArgs []posArg
	optionalPosArgs []posArg
	varLenArg       *varLenArg
	errors          []error

	parent *Command
}

func (c *Command) Execute() {
	p := newParser(c, os.Args)
	p.parseFlags()
}
