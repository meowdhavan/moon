package moon

import (
	"bytes"
	"fmt"
	"testing"
)

type CustomWriter struct {
	buf bytes.Buffer
}

func (c *CustomWriter) Write(p []byte) (n int, err error) {
	return c.buf.Write(p)
}

func (c *CustomWriter) String() string {
	return c.buf.String()
}

func TestIntroLinePrint(t *testing.T) {
	c := Command{
		Name:       "app",
		AboutShort: "short about",
		AboutLong:  "Long About Section",
	}

	c.Flags().StringFlag(nil, "test-flag", "t", "Test Flag")

	w := CustomWriter{}

	p := DefaultPrinter{
		Writer:           &w,
		SuppressWarnings: false,
	}

	p.printIntroLine(&c)

	got := w.String()
	want := "app - short about\n"

	if got != want {
		t.Errorf("Intro Line Print mismatch; got='%s', want='%s'", got, want)
	}
}

func TestFullHelpPrint(t *testing.T) {
	rootCmd := Command{
		Name:       "app",
		AboutShort: "short about for rootCmd",
		AboutLong:  "Long About Section",
	}

	rootCmd.Flags().StringFlag(nil, "test-flag", "t", "Test Flag")
	rootCmd.StringPosArg(nil, "TEST_ARG", "")

	subCmd := Command{
		Name:       "sub",
		AboutShort: "short about subCmd",
	}

	rootCmd.Subcommand(&subCmd)

	w := CustomWriter{}

	p := DefaultPrinter{
		Writer:           &w,
		SuppressWarnings: false,
	}

	p.printHelp(&rootCmd)

	got := w.String()
	want := `app - short about for rootCmd

Long About Section

Usage:
app [FLAGS] <COMMAND>

Commands:
sub  short about subCmd

Flags:
-t, --test-flag  Test Flag
`

	if got != want {
		t.Errorf("Full Help Print mismatch; got='%s', want='%s'", got, want)
	}
}

func TestIndentPrint(t *testing.T) {
	rootCmd := Command{
		Name:       "app",
		AboutShort: "short about for rootCmd",
		AboutLong:  "Long About Section",
	}

	rootCmd.Flags().StringFlag(nil, "test-flag", "t", "Test Flag")
	rootCmd.StringPosArg(nil, "TEST_ARG", "")

	subCmd := Command{
		Name:       "sub",
		AboutShort: "short about subCmd",
	}

	rootCmd.Subcommand(&subCmd)

	w := CustomWriter{}

	p := DefaultPrinter{
		Writer:           &w,
		SuppressWarnings: false,
		IndentLength:     4,
	}

	p.printHelp(&rootCmd)

	got := w.String()
	want := `app - short about for rootCmd

Long About Section

Usage:
    app [FLAGS] <COMMAND>

Commands:
    sub  short about subCmd

Flags:
    -t, --test-flag  Test Flag
`

	if got != want {
		t.Errorf("Indent Print mismatch; got='%s', want='%s'", got, want)
	}
}

func TestLocalFlagsPrint(t *testing.T) {
	rootCmd := Command{
		Name:       "app",
		AboutShort: "short about for rootCmd",
		AboutLong:  "Long About Section",
	}

	localFlagCount := 2

	for i := range localFlagCount {
		rootCmd.Flags().StringFlag(nil, fmt.Sprintf("local-flag-%d", i+1), "t", fmt.Sprintf("Local Flag %d", i+1))
	}

	w := CustomWriter{}

	p := DefaultPrinter{
		Writer:           &w,
		SuppressWarnings: false,
	}

	p.printHelp(&rootCmd)

	got := w.String()
	want := `app - short about for rootCmd

Long About Section

Usage:
app [FLAGS]

Flags:
-t, --local-flag-1  Local Flag 1
-t, --local-flag-2  Local Flag 2
`

	if got != want {
		t.Errorf("Local Flags Print mismatch; got='%s', want='%s'", got, want)
	}
}

func TestInitialIndentPrint(t *testing.T) {
	rootCmd := Command{
		Name:       "app",
		AboutShort: "short about for rootCmd",
	}

	rootCmd.Flags().StringFlag(nil, "local-flag-1", "", "Local Flag 1")
	rootCmd.Flags().StringFlag(nil, "local-flag-2", "", "Local Flag 2")
	rootCmd.Flags().StringFlag(nil, "local-flag-3", "", "Local Flag 3")

	w := CustomWriter{}

	p := DefaultPrinter{
		Writer:           &w,
		SuppressWarnings: false,
	}

	p.printHelp(&rootCmd)

	got := w.String()
	want := `app - short about for rootCmd

Usage:
app [FLAGS]

Flags:
--local-flag-1  Local Flag 1
--local-flag-2  Local Flag 2
--local-flag-3  Local Flag 3
`

	if got != want {
		t.Errorf("Initial Indent Print mismatch; got='%s', want='%s'", got, want)
	}
}

func TestGlobalFlagsPrint(t *testing.T) {
	rootCmd := Command{
		Name:       "app",
		AboutShort: "short about for rootCmd",
		AboutLong:  "Long About Section",
	}

	globalFlagCount := 3

	for i := range globalFlagCount {
		rootCmd.GlobalFlags().StringFlag(nil, fmt.Sprintf("global-flag-%d", i+1), "t", fmt.Sprintf("Global Flag %d", i+1))
	}

	w := CustomWriter{}

	p := DefaultPrinter{
		Writer:           &w,
		SuppressWarnings: false,
	}

	p.printHelp(&rootCmd)

	got := w.String()
	want := `app - short about for rootCmd

Long About Section

Usage:
app [FLAGS]

Global Flags:
-t, --global-flag-1  Global Flag 1
-t, --global-flag-2  Global Flag 2
-t, --global-flag-3  Global Flag 3
`

	if got != want {
		t.Errorf("Global Flags Print mismatch; got='%s', want='%s'", got, want)
	}
}


func TestLocalAndGlobalFlagsPrint(t *testing.T) {
	rootCmd := Command{
		Name:       "app",
		AboutShort: "short about for rootCmd",
		AboutLong:  "Long About Section",
	}

	localFlagCount := 2

	for i := range localFlagCount {
		rootCmd.Flags().StringFlag(nil, fmt.Sprintf("local-flag-%d", i+1), "t", fmt.Sprintf("Local Flag %d", i+1))
	}

	globalFlagCount := 3

	for i := range globalFlagCount {
		rootCmd.GlobalFlags().StringFlag(nil, fmt.Sprintf("global-flag-%d", i+1), "t", fmt.Sprintf("Global Flag %d", i+1))
	}

	w := CustomWriter{}

	p := DefaultPrinter{
		Writer:           &w,
		SuppressWarnings: false,
	}

	p.printHelp(&rootCmd)

	got := w.String()
	want := `app - short about for rootCmd

Long About Section

Usage:
app [FLAGS]

Flags:
-t, --local-flag-1  Local Flag 1
-t, --local-flag-2  Local Flag 2

Global Flags:
-t, --global-flag-1  Global Flag 1
-t, --global-flag-2  Global Flag 2
-t, --global-flag-3  Global Flag 3
`

	if got != want {
		t.Errorf("Local and Global Flags Print mismatch; got='%s', want='%s'", got, want)
	}
}

func TestFlagFallbackPrint(t *testing.T) {
	rootCmd := Command{
		Name:       "app",
		AboutShort: "short about for rootCmd",
		AboutLong:  "Long About Section",
	}

	rootCmd.Flags().StringFlag(nil, "test-flag-1", "a", "Test Flag 1", Required())
	rootCmd.Flags().StringFlag(nil, "test-flag-2", "b", "Test Flag 2", Env("TEST_ENV_VAR"))
	rootCmd.Flags().StringFlag(nil, "test-flag-3", "c", "Test Flag 3", Env("TEST_ENV_VAR"), Required())
	rootCmd.Flags().StringFlag(nil, "test-flag-4", "d", "Test Flag 4", Default("DEF"))
	rootCmd.Flags().StringFlag(nil, "test-flag-5", "e", "Test Flag 5", Env("TEST_ENV_VAR"), Default("DEF"))

	w := CustomWriter{}

	p := DefaultPrinter{
		Writer:           &w,
		SuppressWarnings: false,
	}

	p.printHelp(&rootCmd)

	got := w.String()
	want := `app - short about for rootCmd

Long About Section

Usage:
app [FLAGS]

Flags:
-a, --test-flag-1  Test Flag 1 (Required)
-b, --test-flag-2  Test Flag 2 [$TEST_ENV_VAR]
-c, --test-flag-3  Test Flag 3 (Required) [$TEST_ENV_VAR]
-d, --test-flag-4  Test Flag 4 (default DEF)
-e, --test-flag-5  Test Flag 5 (default DEF) [$TEST_ENV_VAR]
`

	if got != want {
		t.Errorf("Flag Fallback Print mismatch; got='%s', want='%s'", got, want)
	}
}
