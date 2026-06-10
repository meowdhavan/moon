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
	printHelp(*Command)
	printWarnings(*[]error)
	printFullUsage(*Command, *[]error, *[]error)
}

type defaultPrinter struct {
	w                io.Writer
	suppressWarnings bool
	Heading          func(string) string
	Focus            func(string) string
}

func NewDefaultPrinter(w io.Writer, suppressWarnings bool) defaultPrinter {
	return defaultPrinter{
		w: w,
		suppressWarnings: suppressWarnings,
		Heading: func(s string) string {
			return fmt.Sprintf("\x1b[4m%s\x1b[24m", s)
		},
		Focus: func(s string) string {
			return s
		},
	}
}

func (p *defaultPrinter) printErrors(errors *[]error) {
	if len(*errors) == 0 {
		return
	}

	if len(*errors) == 1 {
		fmt.Fprintf(p.w, "%s\n", p.Heading("Error:"))
	} else {
		fmt.Fprintf(p.w, "%s\n", p.Heading("Errors ("+strconv.Itoa(len(*errors))+"):"))
	}

	for _, e := range *errors {
		fmt.Fprintf(p.w, "    - %s\n", e.Error())
	}

	fmt.Println(p.w)
}

func (p *defaultPrinter) printWarnings(warnings *[]error) {
	if p.suppressWarnings || len(*warnings) == 0 {
		return
	}

	if len(*warnings) == 1 {
		fmt.Fprintf(p.w, "%s\n", p.Heading("Warning:"))
	} else {
		fmt.Fprintf(p.w, "%s\n", p.Heading("Warnings ("+strconv.Itoa(len(*warnings))+"):"))
	}

	for _, e := range *warnings {
		fmt.Fprintf(p.w, "    - %s\n", e.Error())
	}

	fmt.Println(p.w)
}

func (p *defaultPrinter) printHelp(c *Command) {
	p.printIntroLine(c)
	fmt.Fprintln(p.w)
	p.printAboutLong(c)
	if c.AboutLong != "" {
		fmt.Fprintln(p.w)
	}
	p.printFullUsage(c, &[]error{}, &[]error{})
}

func (p *defaultPrinter) printIntroLine(c *Command) {
	fmt.Fprint(p.w, p.Focus(c.Name))
	if c.AboutShort != "" {
		fmt.Fprint(p.w, " - ")
		fmt.Fprint(p.w, c.AboutShort)
	}

	fmt.Fprintln(p.w)
}

func (p *defaultPrinter) printFullUsage(c *Command, errors *[]error, warnings *[]error) {
	p.printErrors(errors)
	p.printWarnings(warnings)
	p.printUsage(c)
	fmt.Fprintln(p.w)
	p.printSubcommands(c)
	if len(c.subcommands) > 0 {
		fmt.Fprintln(p.w)
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
		commands = append(commands, cur.Name)
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
			fmt.Fprintf(p.w, " <...%s>", c.varLenArg.name)
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
		fmt.Fprintf(tw, "    %s", p.Focus(s.Name))

		fmt.Fprintf(tw, "\t%s", s.AboutShort)
	}

	fmt.Fprintln(tw)

	tw.Flush()
}

func (p *defaultPrinter) printFlags(c *Command) {
	flags := []*Flag{}

	var cur *Command
	cur = c

	for cur != nil {
		flags = append(flags, cur.flags...)
		cur = cur.parent
	}

	if len(flags) == 0 {
		return
	}

	fmt.Fprintln(p.w, p.Heading("Flags:"))

	tw := tabwriter.NewWriter(p.w, 5, 0, 2, ' ', 0)

	for _, f := range flags {
		fmt.Fprintf(tw, "    %s", p.Focus("--"+f.name))

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
