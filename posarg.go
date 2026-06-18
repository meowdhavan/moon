package moon

import "github.com/meowdhavan/moon/converter"

// PosArg represents a positional argument parsed from the command line.
// There are two types of PorArg: required and optional.
// The required PosArgs are always parsed first.
type PosArg struct {
	Variable
}

type posArgCollection struct {
	requiredPosArgs []*PosArg
	optionalPosArgs []*PosArg
}

// String adds a [PosArg] of type string.
//
// Example:
//
//	var source string
//	cmd.PosArgs().String(&source, "source", "Source directory", properties...)
func (c *posArgCollection) String(
	target *string,
	name string,
	about string,
	properties ...variableProperty,
) {
	posArg := &PosArg{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToString(s)
				if err != nil {
					return err
				}

				*target = v
				return nil
			},
		},
	}

	for _, opt := range properties {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

// Bool adds a [PosArg] of type bool.
//
// Example:
//
//	var apply bool
//	cmd.PosArgs().Bool(&apply, "apply", "Apply changes", properties...)
func (c *posArgCollection) Bool(
	target *bool,
	name string,
	about string,
	properties ...variableProperty,
) {
	*target = false

	posArg := &PosArg{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToBool(s)
				if err != nil {
					return err
				}

				*target = v
				return nil
			},
		},
	}

	for _, opt := range properties {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}

// Int adds a [PosArg] of type int.
//
// Example:
//
//	var count int
//	cmd.PosArgs().Int(&count, "count", "Number of items", properties...)
func (c *posArgCollection) Int(
	target *int,
	name string,
	about string,
	properties ...variableProperty,
) {
	posArg := &PosArg{
		Variable: Variable{
			name:    name,
			aliases: []string{},
			about:   about,
			setValue: func(s string) error {
				v, err := converter.ToInt(s)
				if err != nil {
					return err
				}

				*target = v
				return nil
			},
		},
	}

	for _, opt := range properties {
		opt(&posArg.Variable)
	}

	if posArg.isRequired {
		c.requiredPosArgs = append(c.requiredPosArgs, posArg)
	} else {
		c.optionalPosArgs = append(c.optionalPosArgs, posArg)
	}
}
