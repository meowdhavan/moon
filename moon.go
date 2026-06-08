package moon

import "os"

type Command struct {
	Names      []string
	AboutShort string
	AboutLong  string
	Run        func() error

	subcommands map[string]*Command
	flags       []flag
	posArgs     []posArg
	errors      []error

	parent *Command
}

func (c *Command) Execute() {
	p := newParser(os.Args)
	p.parseFlags(c)
}