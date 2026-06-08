package moon

import "os"

func Execute(rootCmd *Command) {
	showHelp := false

	rootCmd.AddBoolFlag(&showHelp, []string{"help"}, "h", "Show help message", false)

	parser := newParser(rootCmd, os.Args)
	parser.parseFlags()

	cmd := parser.currentCmd

	printer := newPrinter(os.Stdout)

	if showHelp {
		printer.printHelp(cmd)
		os.Exit(0)
	}

	printer.printError(&parser)
	if len(parser.errors) > 0 {
		printer.newLine()
	}

	printer.printWarning(&parser)
	if len(parser.warnings) > 0 {
		printer.newLine()
	}

	if len(parser.errors) > 0 {
		printer.printFullUsage(cmd)
		os.Exit(3)
	}

	cmd.Run()
}
