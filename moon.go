package moon

import "os"

type Moon struct {
	rootCmd *Command
	printer Printer
}

type moonOption func(*Moon)

func WithPrinter(p Printer) moonOption {
	return func(m *Moon) {
		m.printer = p
	}
}

func NewMoon(rootCmd *Command, options ...moonOption) *Moon {
	m := &Moon{
		rootCmd: rootCmd,
	}

	for _, opt:= range options {
		opt(m)
	}

	if m.printer == nil {
		p := NewDefaultPrinter(os.Stdout, false)
		m.printer = &p
	}

	return m
}

func (m *Moon) Execute() {
	showHelp := false

	m.rootCmd.BoolFlag(&showHelp, "help", "h", "Show help message")

	parser := newParser(m.rootCmd, os.Args)
	parser.parse()

	cmd := parser.currentCmd

	if !parser.unrecognizedSubcommand && (showHelp || cmd.Run == nil) {
		m.printer.printHelp(cmd)
		os.Exit(0)
	}

	if parser.unrecognizedSubcommand || len(parser.errors) > 0 {
		m.printer.printFullUsage(cmd, &parser.errors, &parser.warnings)
		os.Exit(3)
	}

	cmd.Run()
}
