package moon

import "os"

type Moon struct {
	rootCmd          *Command
	SuppressWarnings bool
	Printer          Printer
}

func NewMoon(rootCmd *Command) *Moon {
	p := newDefaultPrinter(os.Stdout)

	return &Moon{
		rootCmd: rootCmd,
		Printer: &p,
	}
}

func (m *Moon) Execute() {
	showHelp := false

	m.rootCmd.BoolFlag(&showHelp, "help", "h", About("Show help message"))

	parser := newParser(m.rootCmd, os.Args)
	parser.parseFlags()

	cmd := parser.currentCmd

	if showHelp || cmd.Run == nil {
		m.Printer.printHelp(cmd)
		os.Exit(0)
	}

	m.Printer.printError(&parser)
	if len(parser.errors) > 0 {
		m.Printer.newLine()
	}

	if !m.SuppressWarnings {
		m.Printer.printWarning(&parser)
		if len(parser.warnings) > 0 {
			m.Printer.newLine()
		}
	}

	if len(parser.errors) > 0 {
		m.Printer.printFullUsage(cmd)
		os.Exit(3)
	}

	cmd.Run()
}
