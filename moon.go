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

	cmdSeen := map[*Command]struct{}{}
	globalFlagSeen := map[string]struct{}{}

	queue := []*Command{m.RootCmd}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		_, found := cmdSeen[cur]
		if found {
			errMsg := fmt.Sprintf("Subcommand loop present: %v", cur)
			err := errors.New(errMsg)
			errs = append(errs, err)

			continue
		}

		cmdSeen[cur] = struct{}{}

		// Check Global Flag

		for _, f := range cur.globalFlags.flags {
			names := []string{}

			if f.name != "" {
				names = append(names, "--"+f.name)
			}

			for _, alias := range f.aliases {
				names = append(names, "--"+alias)
			}

			if f.shortName != "" {
				names = append(names, "-"+f.shortName)
			}

			for _, name := range names {
				_, found := globalFlagSeen[name]
				if found {
					errMsg := fmt.Sprintf("Conflicting global flag name present for command %v and flag %v: %s", cur, f, name)
					err := errors.New(errMsg)
					errs = append(errs, err)
				}

				globalFlagSeen[name] = struct{}{}
			}
		}

		// Check Local Flags

		localFlagSeen := map[string]struct{}{}

		for _, f := range cur.localFlags.flags {
			names := []string{}

			if f.name != "" {
				names = append(names, "--"+f.name)
			}

			for _, alias := range f.aliases {
				names = append(names, "--"+alias)
			}

			if f.shortName != "" {
				names = append(names, "-"+f.shortName)
			}

			for _, name := range names {
				_, found := localFlagSeen[name]
				if found {
					errMsg := fmt.Sprintf("Conflicting local flag name present for command %v and flag %v: %s", cur, f, name)
					err := errors.New(errMsg)
					errs = append(errs, err)
				}

				localFlagSeen[name] = struct{}{}
			}
		}

		// Check Subcommands

		for _, sub := range cur.subcommands {
			queue = append(queue, sub)
		}
	}

	return errs
}
