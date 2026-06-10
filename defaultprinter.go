package moon

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Style int

const (
	StyleBold = iota
	StyleUnderline
	StyleUppercase
)

type DefaultPrinter struct {
	Writer           io.Writer
	SuppressWarnings bool
	IndentLength     int
	HeadingStyle     []Style
	FocusStyle       []Style
}

func formatText(s string, styles *[]Style) string {
	for _, sty := range *styles {
		switch sty {
		case StyleBold:
			s = fmt.Sprintf("\x1b[4m%s\x1b[24m", s)
		case StyleUnderline:
			s = fmt.Sprintf("\033[4m%s\033[0m", s)
		case StyleUppercase:
			s = strings.ToUpper(s)
		}
	}

	return s
}

func (p *DefaultPrinter) Heading(s string) string {
	return formatText(s, &p.HeadingStyle)
}

func (p *DefaultPrinter) Focus(s string) string {
	return formatText(s, &p.FocusStyle)
}

func (p *DefaultPrinter) getIndent() string {
	return strings.Repeat(" ", p.IndentLength)
}

func (p *DefaultPrinter) printErrors(errors *[]error) {
	if len(*errors) == 0 {
		return
	}

	if len(*errors) == 1 {
		fmt.Fprintf(p.Writer, "%s\n", p.Heading("Error:"))
	} else {
		fmt.Fprintf(p.Writer, "%s\n", p.Heading("Errors ("+strconv.Itoa(len(*errors))+"):"))
	}

	for _, e := range *errors {
		fmt.Fprintf(p.Writer, "%s- %s\n", p.getIndent(), e.Error())
	}

	fmt.Println(p.Writer)
}

func (p *DefaultPrinter) printWarnings(warnings *[]error) {
	if p.SuppressWarnings || len(*warnings) == 0 {
		return
	}

	if len(*warnings) == 1 {
		fmt.Fprintf(p.Writer, "%s\n", p.Heading("Warning:"))
	} else {
		fmt.Fprintf(p.Writer, "%s\n", p.Heading("Warnings ("+strconv.Itoa(len(*warnings))+"):"))
	}

	for _, e := range *warnings {
		fmt.Fprintf(p.Writer, "%s- %s\n", p.getIndent(), e.Error())
	}

	fmt.Println(p.Writer)
}

func (p *DefaultPrinter) printHelp(c *Command) {
	p.printIntroLine(c)
	fmt.Fprintln(p.Writer)
	p.printAboutLong(c)
	if c.AboutLong != "" {
		fmt.Fprintln(p.Writer)
	}
	p.printFullUsage(c, &[]error{}, &[]error{})
}

func (p *DefaultPrinter) printIntroLine(c *Command) {
	fmt.Fprint(p.Writer, p.Focus(c.Name))
	if c.AboutShort != "" {
		fmt.Fprint(p.Writer, " - ")
		fmt.Fprint(p.Writer, c.AboutShort)
	}

	fmt.Fprintln(p.Writer)
}

func (p *DefaultPrinter) printFullUsage(c *Command, errors *[]error, warnings *[]error) {
	p.printErrors(errors)
	p.printWarnings(warnings)
	p.printUsage(c)
	fmt.Fprintln(p.Writer)
	p.printSubcommands(c)
	if len(c.subcommands) > 0 {
		fmt.Fprintln(p.Writer)
	}
	p.printFlags(c)
}

func (p *DefaultPrinter) printAboutLong(c *Command) {
	if c.AboutLong == "" {
		return
	}

	fmt.Fprintln(p.Writer, c.AboutLong)
}

func (p *DefaultPrinter) printUsage(c *Command) {
	fmt.Fprintln(p.Writer, p.Heading("Usage:"))

	fmt.Fprint(p.Writer, p.getIndent())

	p.printCommand(c)

	if len(c.globalFlags.flags) > 0 || len(c.localFlags.flags) > 0 {
		fmt.Fprint(p.Writer, " [FLAGS]")
	}

	if len(c.subcommands) > 0 {
		fmt.Fprint(p.Writer, " <COMMAND>")
	} else {
		for _, a := range c.requiredPosArgs {
			fmt.Fprintf(p.Writer, " <%s>", a.name)
		}

		for _, a := range c.optionalPosArgs {
			fmt.Fprintf(p.Writer, " <%s>", a.name)
		}

		if c.varLenArg != nil {
			fmt.Fprintf(p.Writer, " <...%s>", c.varLenArg.name)
		}
	}

	fmt.Fprintln(p.Writer)
}

func (p *DefaultPrinter) printCommand(c *Command) {
	var cur *Command
	cur = c

	commands := []string{}

	for cur != nil {
		commands = append(commands, cur.Name)
		cur = cur.parent
	}

	slices.Reverse(commands)

	fmt.Fprintf(p.Writer, "%s", strings.Join(commands, " "))
}

func (p *DefaultPrinter) printSubcommands(c *Command) {
	if len(c.subcommands) == 0 {
		return
	}

	fmt.Fprintln(p.Writer, p.Heading("Commands:"))

	tw := tabwriter.NewWriter(p.Writer, 5, 0, 2, ' ', 0)

	for _, s := range c.subcommands {
		fmt.Fprintf(tw, "%s%s", p.getIndent(), p.Focus(s.Name))

		fmt.Fprintf(tw, "\t%s", s.AboutShort)
	}

	fmt.Fprintln(tw)

	tw.Flush()
}

func (p *DefaultPrinter) printFlagsUtil(flags []*Flag) {
	tw := tabwriter.NewWriter(p.Writer, 5, 0, 2, ' ', 0)

	for _, f := range flags {
		fmt.Fprint(tw, p.getIndent())

		if f.shortName != "" {
			fmt.Fprintf(tw, "%s, ", p.Focus("-"+f.shortName))
		}

		fmt.Fprintf(tw, "%s", p.Focus("--"+f.name))

		fmt.Fprintf(tw, "\t%s", f.about)
		fmt.Fprintln(tw)
	}

	tw.Flush()
}

func (p *DefaultPrinter) printFlags(c *Command) {
	globalFlags := []*Flag{}

	var cur *Command
	cur = c

	for cur != nil {
		globalFlags = append(globalFlags, cur.globalFlags.flags...)
		cur = cur.parent
	}

	if len(c.localFlags.flags) > 0 {
		fmt.Fprintln(p.Writer, p.Heading("Flags:"))
		p.printFlagsUtil(c.localFlags.flags)
	}

	if len(globalFlags) > 0 {
		if len(c.localFlags.flags) > 0 {
			fmt.Fprintln(p.Writer)
		}

		fmt.Fprintln(p.Writer, p.Heading("Global Flags:"))
		p.printFlagsUtil(globalFlags)
	}
}
