package moon

// Command represents a command or subcommand in the application's CLI. It acts as the core building
// block for your CLI application. Each command can hold its own set of flags ([Flag]), and either
// positional arguments ([PosArg]), variadic arguments ([VarArgs]), or subcommands ([Command]).
//
// Example:
//
//	rootCmd := &moon.Command{
//		Name:       "app",
//		Version:    "1.0.0",
//		AboutShort: "A sample CLI application",
//		Run: func() {
//			fmt.Println("Running app...")
//		},
//	}
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

// Flags returns a collection to define local flags ([Flag]). Local flags are only available on this
// specific command.
//
// Example:
//
//	var verbose bool
//	cmd.Flags().Bool(&verbose, "verbose", "v", "Enable verbose output")
func (c *Command) Flags() *flagCollection {
	return &c.localFlags
}

// GlobalFlags returns a collection to define global flags ([Flag]).
// Global flags are available on this command and all of its subcommands.
//
// Example:
//
//	var configPath string
//	cmd.GlobalFlags().String(&configPath, "config", "c", "Path to config file")
func (c *Command) GlobalFlags() *flagCollection {
	return &c.globalFlags
}

// PosArgs returns a collection to define positional arguments ([PosArg]) for this command.
// Positional arguments are parsed in the order they are defined. However, the required ones are
// always parsed before the optional ones, regardless of their order of definition.
//
// If a command has positional arguments ([PosArg]), it should not have subcommands ([Command]). If
// a command has optional positional arguments, it should not have variadic arguments ([VarArgs]).
// Variadic positional arguments are allowed with required positional arguments.
//
// Example:
//
//	var input string
//	cmd.PosArgs().String(&input, "input", "Input file path", moon.Required())
func (c *Command) PosArgs() *posArgCollection {
	return &c.posArgs
}

// VarArgs returns a collection to define variadic arguments ([VarArgs]) for this command. Variadic
// arguments capture all remaining positional arguments provided.
//
// If a command has variadic arguments, it should not have subcommands ([Command]) or optional
// positional arguments ([PosArg]). Required positional arguments are allowed.
//
// Example:
//
//	var files []string
//	cmd.VarArgs().String(&files, "files", "List of files to process")
func (c *Command) VarArgs() *varArgs {
	return &c.varArgs
}

// Subcommand adds a subcommand to this command.
//
// If a command has subcommands, it should not have positional arguments ([PosArg]) or variadic
// arguments ([VarArgs]).
//
// Example:
//
//	sub := &moon.Command{Name: "serve"}
//	rootCmd.Subcommand(sub)
func (c *Command) Subcommand(s *Command) {
	s.parent = c
	c.subcommands = append(c.subcommands, s)
}
