package moon

import (
	"fmt"
	"os"
)

// Moon is the main application struct that manages the root command and printing.
type Moon struct {
	RootCmd *Command
	Printer Printer
}

var showHelp bool
var showVersion bool

func initializeRoot(rootCmd *Command) {
	queue := []*Command{rootCmd}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if !cur.SuppressHelp {
			cur.Flags().Bool(&showHelp, "help", "h", "Show help message")
		}

		if cur.Version != "" {
			cur.Flags().Bool(&showVersion, "version", "v", "Show version")
		}

		for _, sub := range cur.subcommands {
			queue = append(queue, sub)
		}
	}
}

// NewMoon initializes a new Moon application with the given root command and the default [Printer].
// It also recursively sets up standard flags like help and version for all subcommands.
func NewMoon(rootCmd *Command) *Moon {
	p := DefaultPrinter{
		SuppressWarnings: false,
		IndentLength:     4,
		HeadingStyle:     []Style{StyleUnderline},
		HelperMaxLength:  80,
	}

	m := &Moon{
		RootCmd: rootCmd,
		Printer: &p,
	}

	initializeRoot(rootCmd)

	return m
}

// Execute parses the command-line arguments and runs the appropriate command.
// It handles help, version output, parsing errors, and command execution.
func (m *Moon) Execute() {
	p := newParser(m.RootCmd, os.Args)
	p.parse()

	cmd := p.currentCmd

	if showVersion {
		fmt.Print(m.Printer.printVersion(m.RootCmd))
		os.Exit(0)
	}

	if !p.unrecognizedSubcommand && (showHelp || cmd.Run == nil) {
		fmt.Print(m.Printer.printHelp(cmd))
		os.Exit(0)
	}

	if len(p.errors) > 0 {
		fmt.Print(m.Printer.printFullUsage(cmd, &p.errors, &p.warnings))
		os.Exit(3)
	}

	if len(p.warnings) > 0 {
		fmt.Print(m.Printer.printWarnings(&p.warnings))
	}

	cmd.Run()
}
