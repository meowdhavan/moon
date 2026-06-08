package moon

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
