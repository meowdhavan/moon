package moon

// Printer defines the interface for formatting and printing command outputs, such as versions, help
// menus, warnings, and usage details. This interface can be implemented to provide custom
// formatting for the CLI.
type Printer interface {
	printVersion(*Command) string
	printHelp(*Command) string
	printWarnings(*[]error) string
	printFullUsage(*Command, *[]error, *[]error) string
}
