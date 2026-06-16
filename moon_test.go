package moon

import (
	"testing"
)

func TestHelpPrint(t *testing.T) {
	c := &Command{
		Name:       "app",
		AboutShort: "short about",
		AboutLong:  "Long About Section",
	}

	c.Flags().String(nil, "test-flag", "t", "Test Flag")

	NewMoon(c)

	parser := newParser(c, []string{"app"})
	parser.parse()

	got := showHelp
	want := false

	if got != want {
		t.Errorf("Moon Help False Flag mismatch; got=%v, want=%v", got, want)
	}

	parser = newParser(c, []string{"app", "--help"})
	parser.parse()

	got = showHelp
	want = true

	if got != want {
		t.Errorf("Moon Help True Flag mismatch; got=%v, want=%v", got, want)
	}

}
