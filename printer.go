package moon

import (
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"
)

type printer struct {
	w       io.Writer
	Heading func(string) string
	Focus   func(string) string
}

func newPrinter(w io.Writer) printer {
	return printer{
		w: w,
		Heading: func(s string) string {
			return underlineText(s)
		},
		Focus: func(s string) string {
			return s
		},
	}
}

func boldText(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

func underlineText(s string) string {
	return fmt.Sprintf("\x1b[4m%s\x1b[24m", s)
}

func (p *printer) newLine() {
	fmt.Fprintln(p.w)
}

func (p *printer) printError(parser *parser) {
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

func (p *printer) printWarning(parser *parser) {
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

func (p *printer) printHelp(c *Command) {
	p.printIntroLine(c)
	p.newLine()
	p.printAboutLong(c)
	p.newLine()
	p.printFullUsage(c)
}

func (p *printer) printIntroLine(c *Command) {
	fmt.Fprint(p.w, p.Focus(c.Names[0]))
	if c.AboutShort != "" {
		fmt.Fprint(p.w, " - ")
		fmt.Fprint(p.w, c.AboutShort)
	}

	fmt.Fprintln(p.w)
}

func (p *printer) printFullUsage(c *Command) {
	p.printUsage(c)
	p.newLine()
	p.printSubcommands(c)
	p.newLine()
	p.printFlags(c)
}

func (p *printer) printAboutLong(c *Command) {
	if c.AboutLong == "" {
		return
	}

	fmt.Fprintln(p.w, c.AboutLong)
}

func (p *printer) printUsage(c *Command) {
	fmt.Fprintln(p.w, p.Heading("Usage:"))

	fmt.Fprint(p.w, "    ")

	fmt.Fprintf(p.w, "%s", c.Names[0])

	if len(c.flags) > 0 {
		fmt.Fprint(p.w, " [FLAGS]")
	}

	if len(c.subcommands) > 0 {
		fmt.Fprint(p.w, " <COMMAND>")
	} else {
		fmt.Fprint(p.w, " ")
		// TODO: Pos Args
	}

	fmt.Fprintln(p.w)
}

func (p *printer) printSubcommands(c *Command) {
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

func (p *printer) printFlags(c *Command) {
	if len(c.flags) == 0 {
		return
	}

	fmt.Fprintln(p.w, p.Heading("Flags:"))

	tw := tabwriter.NewWriter(p.w, 5, 0, 2, ' ', 0)

	for _, f := range c.flags {
		fmt.Fprintf(tw, "    %s", p.Focus("--"+f.longNames[0]))

		if f.shortName != "" {
			fmt.Fprintf(tw, "\t%s", p.Focus("-"+f.shortName))
		} else {
			fmt.Fprintf(tw, "\t")
		}

		fmt.Fprintf(tw, "\t%s", f.about)
		fmt.Fprintln(tw)
	}

	tw.Flush()
}
