package moon

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Printer interface {
	newLine()
	printError(*parser) 
	printWarning(*parser) 
	printHelp(*Command) 
	printFullUsage(*Command) 
}

type defaultPrinter struct {
	w       io.Writer
	Heading func(string) string
	Focus   func(string) string
}

func newDefaultPrinter(w io.Writer) defaultPrinter {
	return defaultPrinter{
		w: w,
		Heading: func(s string) string {
			return fmt.Sprintf("\x1b[4m%s\x1b[24m", s)
		},
		Focus: func(s string) string {
			return s
		},
	}
}

func (p *defaultPrinter) newLine() {
	fmt.Fprintln(p.w)
}

func (p *defaultPrinter) printError(parser *parser) {
	if len(parser.errors) == 0 {
		return
	}

	if len(parser.errors) == 1 {
		fmt.Fprintf(p.w, "%s\n", p.Heading("Error:"))
	} else {
		fmt.Fprintf(p.w, "%s\n", p.Heading("Errors (" + strconv.Itoa(len(parser.errors)) + "):"))
	}

	for _, e := range parser.errors {
		fmt.Fprintf(p.w, "    - %s\n", e.Error())
	}
}

func (p *defaultPrinter) printWarning(parser *parser) {
	if len(parser.warnings) == 0 {
		return
	}

	if len(parser.warnings) == 1 {
		fmt.Fprintf(p.w, "%s\n", p.Heading("Warning:"))
	} else {
		fmt.Fprintf(p.w, "%s\n", p.Heading("Warnings (" + strconv.Itoa(len(parser.warnings)) + "):"))
	}

	for _, e := range parser.warnings {
		fmt.Fprintf(p.w, "    - %s\n", e.Error())
	}
}

func (p *defaultPrinter) printHelp(c *Command) {
	p.printIntroLine(c)
	p.newLine()
	p.printAboutLong(c)
	if c.AboutLong != "" {
		p.newLine()
	}
	p.printFullUsage(c)
}

func (p *defaultPrinter) printIntroLine(c *Command) {
	fmt.Fprint(p.w, p.Focus(c.Names[0]))
	if c.AboutShort != "" {
		fmt.Fprint(p.w, " - ")
		fmt.Fprint(p.w, c.AboutShort)
	}

	fmt.Fprintln(p.w)
}

func (p *defaultPrinter) printFullUsage(c *Command) {
	p.printUsage(c)
	p.newLine()
	p.printSubcommands(c)
	if len(c.subcommands) > 0 {
		p.newLine()
	}
	p.printFlags(c)
}

func (p *defaultPrinter) printAboutLong(c *Command) {
	if c.AboutLong == "" {
		return
	}

	fmt.Fprintln(p.w, c.AboutLong)
}

func (p *defaultPrinter) printUsage(c *Command) {
	fmt.Fprintln(p.w, p.Heading("Usage:"))

	fmt.Fprint(p.w, "    ")

	var cur *Command
	cur = c

	commands := []string{}

	for cur != nil {
		commands = append(commands, cur.Names[0])
		cur = cur.parent
	}

	slices.Reverse(commands)

	fmt.Fprintf(p.w, "%s", strings.Join(commands, " "))

	if len(c.flags) > 0 {
		fmt.Fprint(p.w, " [FLAGS]")
	}

	if len(c.subcommands) > 0 {
		fmt.Fprint(p.w, " <COMMAND>")
	} else {
		for _, a := range c.requiredPosArgs {
			fmt.Fprintf(p.w, " <%s>", a.name)
		}

		for _, a := range c.optionalPosArgs {
			fmt.Fprintf(p.w, " <%s>", a.name)
		}

		if c.varLenArg != nil {
			fmt.Fprintf(p.w, " ...<%s>", c.varLenArg.name)
		}
	}

	fmt.Fprintln(p.w)
}

func (p *defaultPrinter) printSubcommands(c *Command) {
	if len(c.subcommands) == 0 {
		return
	}

	fmt.Fprintln(p.w, p.Heading("Commands:"))

	tw := tabwriter.NewWriter(p.w, 5, 0, 2, ' ', 0)

	for _, s := range c.subcommands {
		fmt.Fprintf(tw, "    %s", p.Focus(s.Names[0]))

		fmt.Fprintf(tw, "\t%s", s.AboutShort)
	}

	fmt.Fprintln(tw)

	tw.Flush()
}

func (p *defaultPrinter) printFlags(c *Command) {
	if len(c.flags) == 0 {
		return
	}

	fmt.Fprintln(p.w, p.Heading("Flags:"))

	tw := tabwriter.NewWriter(p.w, 5, 0, 2, ' ', 0)

	var cur *Command
	cur = c

	for cur != nil {
		for _, f := range cur.flags {
			fmt.Fprintf(tw, "    %s", p.Focus("--"+f.name))

			if f.shortName != "" {
				fmt.Fprintf(tw, "\t%s", p.Focus("-"+f.shortName))
			} else {
				fmt.Fprintf(tw, "\t")
			}

			fmt.Fprintf(tw, "\t%s", f.about)
			fmt.Fprintln(tw)
		}

		cur = cur.parent
	}

	tw.Flush()
}
