package moon

import "os"

func Execute(rootCmd *Command) {
	parser := newParser(rootCmd, os.Args)
	parser.parseFlags()

	cmd := parser.currentCmd

	if len(parser.errors) > 0 {
		// TODO
	}

	if len(parser.warnings) > 0 {
		// TODO
	}

	cmd.Run()
}
