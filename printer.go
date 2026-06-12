package moon

type Printer interface {
	printHelp(*Command) string
	printWarnings(*[]error) string
	printFullUsage(*Command, *[]error, *[]error) string
}
