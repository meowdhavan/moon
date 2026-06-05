package mon

import (
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
				if p.tokenIdx+1 < len(p.tokens) && !isFlag(p.tokens[p.tokenIdx+1]) {
					f.setValue(p.tokens[p.tokenIdx+1])
					f.isValueSet = true
					p.tokenIdx++
				} else {
					// Error: No value supplied for flag
				}
			} else {
				f.setValue("true")
				f.isValueSet = true
			}
		} else if isShortFlag(token) {
			for i, c := range token[1:] {
				f, found := p.flagMap["-"+string(c)]
				if !found {
					// Warning: Unrecognized flag
					continue
				}

				if f.requiresVal {
					if i + 2 < len(token) {
						f.setValue(token[i+2:])
						f.isValueSet = true
					} else {
						if p.tokenIdx+1 < len(p.tokens) && !isFlag(p.tokens[p.tokenIdx+1]) {
							f.setValue(p.tokens[p.tokenIdx+1])
							f.isValueSet = true
							p.tokenIdx++
						} else {
							// Error: No value supplied for flag
						}
					}

					break
				} else {
					f.setValue("true")
					f.isValueSet = true
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
