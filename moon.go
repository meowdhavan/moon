package moon

import (
	"fmt"
	"os"
)

type Moon struct {
	RootCmd *Command
	Printer Printer
}

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

	return m
}

func (m *Moon) Execute() {
	showHelp := false

	m.RootCmd.GlobalFlags().BoolFlag(&showHelp, "help", "h", "Show help message") // TODO: Use local flag for help

	p := newParser(m.RootCmd, os.Args)
	p.parse()

	cmd := p.currentCmd

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
