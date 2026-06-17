package moon

// Printer defines the interface for formatting and printing command outputs,
// such as versions, help menus, warnings, and usage details.
type Printer interface {
	printVersion(*Command) string
	printHelp(*Command) string
	printWarnings(*[]error) string
	printFullUsage(*Command, *[]error, *[]error) string
}
