package moon

import "github.com/meowdhavan/moon/converter"

// VarArgs represents the variadic arguments parsed from the command line.
type VarArgs struct {
	Variable
}

type varArgs struct {
	varArg *VarArgs
}

// String ensures that the application accepts [VarArgs] of type string.
func (a *varArgs) String(target *[]string, name string, about string, properties ...variableProperty) {
	*target = []string{}

	v := &VarArgs{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToString(s)
				if err != nil {
					return err
				}

				*target = append(*target, v)
				return nil
			},
		},
	}

	for _, opt := range properties {
		opt(&v.Variable)
	}

	a.varArg = v
}

// Int ensures that the application accepts [Int] of type string.
func (a *varArgs) Int(target *[]int, name string, about string, properties ...variableProperty) {
	*target = []int{}

	v := &VarArgs{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToInt(s)
				if err != nil {
					return err
				}

				*target = append(*target, v)
				return nil
			},
		},
	}

	for _, opt := range properties {
		opt(&v.Variable)
	}

	a.varArg = v
}
