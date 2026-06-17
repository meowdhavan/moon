package moon

import (
	"errors"
	"fmt"
	"os"
)

type Moon struct {
	RootCmd *Command
	Printer Printer
}

var showHelp bool

func initializeRoot(rootCmd *Command) {
	queue := []*Command{rootCmd}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		cur.Flags().Bool(&showHelp, "help", "h", "Show help message")

		for _, sub := range cur.subcommands {
			queue = append(queue, sub)
		}
	}
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

	initializeRoot(rootCmd)

	return m
}

func (m *Moon) Execute() {
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

func (m *Moon) Validate() []error {
	errs := []error{}

	seen := map[*Command]struct{}{}

	queue := []*Command{m.RootCmd}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		_, found := seen[cur]
		if found {
			err := errors.New("Subcommand loop detected: " + cur.Name)
			errs = append(errs, err)
			continue
		}

		// Check Flags

		// Check Subcommands

		for _, sub := range cur.subcommands {
			queue = append(queue, sub)
		}
	}

	return errs
}
