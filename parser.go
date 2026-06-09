package moon

import (
	"errors"
	"os"
	"strings"
)

type parser struct {
	tokens   []string
	tokenIdx int
	flagMap  map[string]*Flag
	subcommandsMap map[string]*Command
	currentCmd *Command
	requiredPosArgIdx int
	optionalPosArgIdx int
	errors []error
	warnings []error
}

func newParser(rootCmd *Command, tokens []string) parser {
	return parser{
		currentCmd: rootCmd,
		flagMap:  make(map[string]*Flag),
		subcommandsMap: make(map[string]*Command),
		tokens: tokens,
		tokenIdx: 1,
		requiredPosArgIdx: 0,
		optionalPosArgIdx: 0,
		errors: make([]error, 0),
		warnings: make([]error, 0),
	}
}

func isLongFlag(s string) bool {
	return len(s) > 2 && strings.HasPrefix(s, "--")
}

func isShortFlag(s string) bool {
	return len(s) > 1 && !isLongFlag(s) && strings.HasPrefix(s, "-")
}

func isFlag(s string) bool {
	return isLongFlag(s) || isShortFlag(s)
}

func (p *parser) fillFlagMap() {
	for _, f := range p.currentCmd.flags {
		if f.name != "" {
			p.flagMap["--"+f.name] = f
		}

		for _, l := range f.aliases {
			if l != "" {
				p.flagMap["--"+l] = f
			}
		}

		if f.shortName != "" {
			p.flagMap["-"+f.shortName] = f
		}
	}
}

func (p *parser) fillSubcommandsMap() {
	for k := range p.subcommandsMap {
		delete(p.subcommandsMap, k)
	}

	for _, s := range p.currentCmd.subcommands {
		for _, name := range s.Names {
			if name != "" {
				p.subcommandsMap[name] = s
			}
		}
	}
}

// Sets the value of a flag, and indicates that there has been an attempt to set a value.
// `isValueSet` must be set to true even if there was an attempt to set an invalid value.
// Not doing so will result in an additional error for an unset value.
func (p *parser) setValue(f *Variable, val string) error {
	err := f.setValue(val)
	f.isValueSet = true

	return err
}

func (p *parser) setNextTokenAsValue(f *Flag) error {
	if p.tokenIdx+1 < len(p.tokens) && !isFlag(p.tokens[p.tokenIdx+1]) {
		err := p.setValue(&f.Variable, p.tokens[p.tokenIdx+1])
		p.tokenIdx++

		if err != nil {
			return errors.New("Invalid value supplied for flag: --" + f.name)
		}
	} else {
		f.isValueSet = true
		return errors.New("No value supplied for flag: --" + f.name)
	}

	return nil
}

func (p *parser) parseFlags() {
	p.fillFlagMap()
	p.fillSubcommandsMap()

	for ; p.tokenIdx < len(p.tokens); p.tokenIdx++ {
		token := p.tokens[p.tokenIdx]

		if isLongFlag(token) {
			f, found := p.flagMap[token]
			if !found {
				warning := errors.New("Unrecognized flag: " + token)
				p.warnings = append(p.warnings, warning)
				continue
			}

			if f.requiresVal {
				err := p.setNextTokenAsValue(f)
				if err != nil {
					p.errors = append(p.errors, err)
				}
			} else {
				err := p.setValue(&f.Variable, "true")
				if err != nil {
					p.errors = append(p.errors, err)
				}
			}
		} else if isShortFlag(token) {
			for i, ch := range token[1:] {
				f, found := p.flagMap["-"+string(ch)]
				if !found {
					warning := errors.New("Unrecognized flag: -" + string(ch))
					p.warnings = append(p.warnings, warning)
					continue
				}

				if f.requiresVal {
					if i + 2 < len(token) {
						err := p.setValue(&f.Variable, token[i+2:])
						if err != nil {
							p.errors = append(p.errors, err)
						}
					} else {
						err := p.setNextTokenAsValue(f)
						if err != nil {
							p.errors = append(p.errors, err)
						}
					}

					break
				} else {
					err := p.setValue(&f.Variable, "true")
					if err != nil {
						p.errors = append(p.errors, err)
					}
				}
			}
		} else {
			if len(p.currentCmd.subcommands) > 0 {
				s, found := p.subcommandsMap[token]
				if !found {
					warning := errors.New("Unrecognized subcommand: " + token)
					p.warnings = append(p.warnings, warning)
					continue
				}

				p.currentCmd = s
				p.fillFlagMap()
				p.fillSubcommandsMap()
			} else {
				if p.requiredPosArgIdx < len(p.currentCmd.requiredPosArgs) { // Required PosArg
					a := p.currentCmd.requiredPosArgs[p.requiredPosArgIdx]
					err := a.setValue(token)
					if err != nil {
						p.errors = append(p.errors, err)
					}

					p.requiredPosArgIdx++
				} else if p.optionalPosArgIdx < len(p.currentCmd.optionalPosArgs) { // Optional PosArg
					a := p.currentCmd.optionalPosArgs[p.optionalPosArgIdx]
					err := a.setValue(token)
					if err != nil {
						p.errors = append(p.errors, err)
					}

					p.optionalPosArgIdx++
				} else { // VarLenArg
					v := p.currentCmd.varLenArg
					if v == nil {
						warning := errors.New("Unrecognized argument: " + token)
						p.warnings = append(p.warnings, warning)
						continue
					}
					v.setValue(token)
				}
			}
		}
	}

	for _, f := range p.flagMap {
		if !f.isValueSet {
			p.setFromFallback(&f.Variable)
		}

		if !f.isValueSet && f.isRequired {
			err := errors.New("Missing value for required flag: " + f.name)
			p.errors = append(p.errors, err)
		}
	}

	for _, a := range p.currentCmd.optionalPosArgs[p.optionalPosArgIdx:] {
		p.setFromFallback(&a.Variable)
	}
}

func (p *parser) setFromFallback(f *Variable) {
	getFromEnv := func() *string {
		if f.env == nil {
			return nil
		}

		val := os.Getenv(*f.env)
		if val == "" {
			return nil
		}
		
		return &val
	}

	getDefault := func() *string {
		if f.defaultVal == nil {
			return nil
		}

		return f.defaultVal
	}

	fallbacks := []func() *string {getFromEnv, getDefault}

	for _, fallback := range fallbacks {
		s := fallback()

		if s != nil {
			err := p.setValue(f, *s)
			if err != nil {
				p.errors = append(p.errors, err)
			}

			return
		}
	}
}
