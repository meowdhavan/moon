package moon

type Printer interface {
	printVersion(*Command) string
	printHelp(*Command) string
	printWarnings(*[]error) string
	printFullUsage(*Command, *[]error, *[]error) string
}
