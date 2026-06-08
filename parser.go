package moon

import (
	"errors"
	"strings"
)

type parser struct {
	currentCmd *Command
	flagMap  map[string]*flag
	tokens   []string
	tokenIdx int
	errors []error
	warnings []error
}

func newParser(rootCmd *Command, tokens []string) parser {
	return parser{
		currentCmd: rootCmd,
		flagMap:  make(map[string]*flag),
		tokens: tokens,
		tokenIdx: 1,
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
		for _, l := range f.longNames {
			if l != "" {
				p.flagMap["--"+l] = &f
			}
		}

		if f.shortName != "" {
			p.flagMap["-"+f.shortName] = &f
		}
	}
}

// Sets the value of a flag, and indicates that there has been an attempt to set a value.
// `isValueSet` must be set to true even if there was an attempt to set an invalid value.
// Not doing so will result in an additional error for an unset value.
func (p *parser) setValue(f *flag, val string) error {
	err := f.setValue(val)
	f.isValueSet = true

	return err
}

func (p *parser) setNextTokenAsValue(f *flag) error {
	if p.tokenIdx+1 < len(p.tokens) && !isFlag(p.tokens[p.tokenIdx+1]) {
		err := p.setValue(f, p.tokens[p.tokenIdx+1])
		p.tokenIdx++

		if err != nil {
			return err
		}
	} else {
		return errors.New("No value supplied for flag")
	}

	return nil
}

func (p *parser) parseFlags() {
	p.fillFlagMap()

	for ; p.tokenIdx < len(p.tokens); p.tokenIdx++ {
		token := p.tokens[p.tokenIdx]

		if isLongFlag(token) {
			f, found := p.flagMap[token]
			if !found {
				p.warnings = append(p.warnings, errors.New("Unrecognized flag: " + token))
				continue
			}

			if f.requiresVal {
				err := p.setNextTokenAsValue(f)
				if err != nil {
					p.errors = append(p.errors, err)
				}
			} else {
				err := p.setValue(f, "true")
				if err != nil {
					p.errors = append(p.errors, err)
				}
			}
		} else if isShortFlag(token) {
			for i, ch := range token[1:] {
				f, found := p.flagMap["-"+string(ch)]
				if !found {
					p.warnings = append(p.warnings, errors.New("Unrecognized flag: -" + string(ch)))
					continue
				}

				if f.requiresVal {
					if i + 2 < len(token) {
						err := p.setValue(f, token[i+2:])
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
					err := p.setValue(f, "true")
					if err != nil {
						p.errors = append(p.errors, err)
					}
				}
			}
		} else {
			s, found := p.currentCmd.subcommands[token]
			if !found {
				p.warnings = append(p.warnings, errors.New("Unrecognized subcommand: " + token))
				continue
			}

			p.currentCmd = s
			p.fillFlagMap()
		}
	}

	for _, f := range p.flagMap {
		if f.isRequired && !f.isValueSet {
			p.errors = append(p.errors, errors.New("No value supplied for required flag: --" + f.longNames[0]))
		}
	}
}
