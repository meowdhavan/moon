package moon

import "github.com/meowdhavan/moon/converter"

type VarArgs struct {
	Variable
}

type varArgs struct {
	varArg *VarArgs
}

func (a *varArgs) String(target *[]string, name string, about string, options ...variableOption) {
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

	for _, opt := range options {
		opt(&v.Variable)
	}

	a.varArg = v
}

func (a *varArgs) Int(target *[]int, name string, about string, options ...variableOption) {
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

	for _, opt := range options {
		opt(&v.Variable)
	}

	a.varArg = v
}
