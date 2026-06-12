package moon

import (
	"fmt"
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
	SuppressWarnings bool
	IndentLength     int
	HelperMaxLength  int
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

func (p *DefaultPrinter) printErrors(errors *[]error) string {
	var b strings.Builder

	if len(*errors) == 0 {
		return b.String()
	}

	if len(*errors) == 1 {
		b.WriteString(fmt.Sprintf("%s\n", p.Heading("Error:")))
	} else {
		b.WriteString(fmt.Sprintf("%s\n", p.Heading("Errors ("+strconv.Itoa(len(*errors))+"):")))
	}

	for _, e := range *errors {
		b.WriteString(fmt.Sprintf("%s- %s\n", p.getIndent(), e.Error()))
	}

	b.WriteString("\n")

	return b.String()
}

func (p *DefaultPrinter) printWarnings(warnings *[]error) string {
	var b strings.Builder

	if p.SuppressWarnings || len(*warnings) == 0 {
		return b.String()
	}

	if len(*warnings) == 1 {
		b.WriteString(fmt.Sprintf("%s\n", p.Heading("Warning:")))
	} else {
		b.WriteString(fmt.Sprintf("%s\n", p.Heading("Warnings ("+strconv.Itoa(len(*warnings))+"):")))
	}

	for _, e := range *warnings {
		b.WriteString(fmt.Sprintf("%s- %s\n", p.getIndent(), e.Error()))
	}

	b.WriteString("\n")

	return b.String()
}

func (p *DefaultPrinter) printHelp(c *Command) string {
	var b strings.Builder

	b.WriteString(p.printIntroLine(c))
	b.WriteString("\n")
	b.WriteString(p.printAboutLong(c))
	if c.AboutLong != "" {
		b.WriteString("\n")
	}
	b.WriteString(p.printFullUsage(c, &[]error{}, &[]error{}))

	return b.String()
}

func (p *DefaultPrinter) printIntroLine(c *Command) string {
	var b strings.Builder

	b.WriteString(p.Focus(c.Name))
	if c.AboutShort != "" {
		b.WriteString(fmt.Sprintf(" - %s", c.AboutShort))
	}

	b.WriteString("\n")

	return b.String()
}

func (p *DefaultPrinter) printFullUsage(c *Command, errors *[]error, warnings *[]error) string {
	var b strings.Builder

	b.WriteString(p.printErrors(errors))
	b.WriteString(p.printWarnings(warnings))
	b.WriteString(p.printUsage(c))
	b.WriteString("\n")
	b.WriteString(p.printSubcommands(c))
	if len(c.subcommands) > 0 {
		b.WriteString("\n")
	}
	b.WriteString(p.printFlags(c))

	return b.String()
}

func (p *DefaultPrinter) printAboutLong(c *Command) string {
	var b strings.Builder

	if c.AboutLong == "" {
		return b.String()
	}

	b.WriteString(c.AboutLong)
	b.WriteString("\n")

	return b.String()
}

func (p *DefaultPrinter) printUsage(c *Command) string {
	var b strings.Builder

	b.WriteString(p.Heading("Usage:"))
	b.WriteString("\n")

	b.WriteString(p.getIndent())

	b.WriteString(p.printCommand(c))

	if len(c.globalFlags.flags) > 0 || len(c.localFlags.flags) > 0 {
		b.WriteString(" [FLAGS]")
	}

	if len(c.subcommands) > 0 {
		b.WriteString(" <COMMAND>")
	} else {
		for _, a := range c.requiredPosArgs {
			b.WriteString(" <")
			b.WriteString(a.name)
			b.WriteString(">")
		}

		for _, a := range c.optionalPosArgs {
			b.WriteString(" <")
			b.WriteString(a.name)
			b.WriteString(">")
		}

		if c.varLenArg != nil {
			b.WriteString(" <...")
			b.WriteString(c.varLenArg.name)
			b.WriteString(">")
		}
	}

	b.WriteString("\n")

	return b.String()
}

func (p *DefaultPrinter) printCommand(c *Command) string {
	var b strings.Builder

	var cur *Command
	cur = c

	commands := []string{}

	for cur != nil {
		commands = append(commands, cur.Name)
		cur = cur.parent
	}

	slices.Reverse(commands)

	b.WriteString(strings.Join(commands, " "))

	return b.String()
}

func (p *DefaultPrinter) printSubcommands(c *Command) string {
	var b strings.Builder

	if len(c.subcommands) == 0 {
		return b.String()
	}

	b.WriteString(p.Heading("Commands:"))
	b.WriteString("\n")

	tw := tabwriter.NewWriter(&b, 5, 0, 2, ' ', 0)

	for _, s := range c.subcommands {
		fmt.Fprintf(tw, "%s%s", p.getIndent(), p.Focus(s.Name))
		fmt.Fprintf(tw, "\t%s", s.AboutShort)
	}

	fmt.Fprintln(tw)

	tw.Flush()

	return b.String()
}

func (p *DefaultPrinter) printFlagLine(tw *tabwriter.Writer, f *Flag, initialIndent bool, splitHelperLines bool) {
	fmt.Fprint(tw, p.getIndent())

	if f.shortName != "" {
		fmt.Fprintf(tw, "%s, ", p.Focus("-"+f.shortName))
	} else if initialIndent {
		fmt.Fprint(tw, strings.Repeat(" ", 4))
	}

	fmt.Fprintf(tw, "%s", p.Focus("--"+f.name))

	fmt.Fprintf(tw, "\t%s", f.about)

	helpers := []string{}

	if f.isRequired {
		helpers = append(helpers, "[required]")
	} else if f.defaultVal != nil {
		defaultVal := getDefault(&f.Variable)
		helpers = append(helpers, fmt.Sprintf("[default: %s]", *defaultVal))
	}

	if f.env != nil {
		var b strings.Builder

		b.WriteString("[env: ")
		b.WriteString(*f.env)

		envVal := getFromEnv(&f.Variable)
		if envVal != nil {
			b.WriteString("=")
			b.WriteString(*envVal)
		}

		b.WriteString("]")
		helpers = append(helpers, b.String())
	}

	if len(f.aliases) > 0 {
		helpers = append(helpers, fmt.Sprintf("[aliases: %s]", strings.Join(f.aliases, ", ")))
	}

	if len(helpers) > 0 {
		helperSep := " "

		if splitHelperLines {
			helperSep = "\n\t"
		}

		fmt.Fprintf(tw, "%s%s", helperSep, strings.Join(helpers, helperSep))
	}

	fmt.Fprintln(tw)
}

func (p *DefaultPrinter) getInitialIndent(flags []*Flag) bool {
	for _, f := range flags {
		if f.shortName != "" {
			return true
		}
	}

	return false
}

func (p *DefaultPrinter) printFlagsUtil(flags []*Flag, splitHelperLines bool) string {
	var b strings.Builder

	initialIndent := p.getInitialIndent(flags)

	tw := tabwriter.NewWriter(&b, 5, 0, 2, ' ', 0)

	for _, f := range flags {
		p.printFlagLine(tw, f, initialIndent, splitHelperLines)
	}

	tw.Flush()

	return b.String()
}

func (p *DefaultPrinter) getMaxHelperLength(flags []*Flag) int {
	flagLines := p.printFlagsUtil(flags, false)

	lines := strings.Split(flagLines, "\n")

	maxLen := 0

	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}

	return maxLen
}

func (p *DefaultPrinter) printFlags(c *Command) string {
	var b strings.Builder

	globalFlags := []*Flag{}

	var cur *Command
	cur = c

	for cur != nil {
		globalFlags = append(globalFlags, cur.globalFlags.flags...)
		cur = cur.parent
	}

	if len(c.localFlags.flags) > 0 {
		maxHelperLen := p.getMaxHelperLength(c.localFlags.flags)
		splitHelperLines := maxHelperLen > p.HelperMaxLength

		b.WriteString(p.Heading("Flags:"))
		b.WriteString("\n")

		b.WriteString(p.printFlagsUtil(c.localFlags.flags, splitHelperLines))
	}

	if len(globalFlags) > 0 {
		maxHelperLen := p.getMaxHelperLength(globalFlags)
		splitHelperLines := maxHelperLen > p.HelperMaxLength

		if len(c.localFlags.flags) > 0 {
			b.WriteString("\n")
		}

		b.WriteString(p.Heading("Global Flags:"))
		b.WriteString("\n")

		b.WriteString(p.printFlagsUtil(globalFlags, splitHelperLines))
	}

	return b.String()
}
