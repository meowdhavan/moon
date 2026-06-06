package mon

import (
	"bytes"
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
		Names: []string{"app"},
		AboutShort: "short about",
		AboutLong: "Long About Section",
	}

	c.AddStringFlag(nil, []string{"test-flag"}, "t", "Test Flag", false)

	w := CustomWriter{}

	p := NewPrinter(&w)

	p.printIntroLine(&c)

	t.Logf("Intro Line: '%s'", w.String())
}