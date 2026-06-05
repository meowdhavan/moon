package mon

import (
	"errors"
	"strings"
)

type parser struct {
	flagMap  map[string]*flag
	tokens   []string
	tokenIdx int
}

func newParser(tokens []string) parser {
	return parser{
		flagMap:  make(map[string]*flag),
		tokens: tokens,
		tokenIdx: 1,
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

func (p *parser) fillFlagMap(c *Command) {
	for _, f := range c.flags {
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

func (p *parser) setValue(f *flag, val string) error {
	err := f.setValue(val)
	f.isValueSet = true // We set it to true even if there is an attempt to set an invalid value
	if err != nil {
		return err
	}

	return nil
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

func (p *parser) parseFlags(c *Command) {
	p.fillFlagMap(c)

	for ; p.tokenIdx < len(p.tokens); p.tokenIdx++ {
		token := p.tokens[p.tokenIdx]

		if isLongFlag(token) {
			f, found := p.flagMap[token]
			if !found {
				// Warning: Unrecognized flag
				continue
			}

			if f.requiresVal {
				err := p.setNextTokenAsValue(f)
				if err != nil {
					c.errors = append(c.errors, err)
				}
			} else {
				err := p.setValue(f, "true")
				if err != nil {
					c.errors = append(c.errors, err)
				}
			}
		} else if isShortFlag(token) {
			for i, ch := range token[1:] {
				f, found := p.flagMap["-"+string(ch)]
				if !found {
					// Warning: Unrecognized flag
					continue
				}

				if f.requiresVal {
					if i + 2 < len(token) {
						err := p.setValue(f, token[i+2:])
						if err != nil {
							c.errors = append(c.errors, err)
						}
					} else {
						err := p.setNextTokenAsValue(f)
						if err != nil {
							c.errors = append(c.errors, err)
						}
					}

					break
				} else {
					err := p.setValue(f, "true")
					if err != nil {
						c.errors = append(c.errors, err)
					}
				}
			}
		} else {
			for _, f := range c.flags {
				if f.isRequired && !f.isValueSet {
					// Error: No value supplied for Required Flag
				}
			}

			// TODO: Subcommand
		}
	}
}
