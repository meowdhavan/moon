package moon

import (
	"errors"
	"fmt"
)

func validateFlag(f *Flag, c *Command) []error {
	errs := []error{}

	if f.requiresVal {
		if f.defaultVal != nil {
			errMsg := fmt.Sprintf("Boolean Flag has a default value for command %s: %s", c.Name, f.name)
			err := errors.New(errMsg)
			errs = append(errs, err)
		}

		if f.isRequired {
			errMsg := fmt.Sprintf("Boolean Flag marked required for command %s: %s", c.Name, f.name)
			err := errors.New(errMsg)
			errs = append(errs, err)
		}

		return errs
	}

	if f.defaultVal != nil && f.isRequired {
		errMsg := fmt.Sprintf("Flag marked required and has a default value for command %s: %s", c.Name, f.name)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	if f.defaultVal != nil {
		err := f.setValue(*f.defaultVal)
		if err != nil {
			errMsg := fmt.Sprintf("Flag does not have a valid default value for command %s: %s", c.Name, f.name)
			err := errors.New(errMsg)
			errs = append(errs, err)
		}
	}

	return errs
}

func validatePosArg(p *PosArg, c *Command) []error {
	errs := []error{}

	if len(p.aliases) > 0 {
		errMsg := fmt.Sprintf("PosArg contains an alias for command %s: %s", c.Name, p.name)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	if p.defaultVal != nil && p.isRequired {
		errMsg := fmt.Sprintf("PosArg marked required and has a default value for command %s: %s", c.Name, p.name)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	if p.env != nil && p.isRequired {
		errMsg := fmt.Sprintf("PosArg marked required and has an env fallback for command %s: %s", c.Name, p.name)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	if p.defaultVal != nil {
		err := p.setValue(*p.defaultVal)
		if err != nil {
			errMsg := fmt.Sprintf("PosArg does not have a valid default value for command %s: %s", c.Name, p.name)
			err := errors.New(errMsg)
			errs = append(errs, err)
		}
	}

	return errs
}

func validateVarArgs(v *VarArgs, c *Command) []error {
	errs := []error{}
	if v == nil {
		return errs
	}

	if len(v.aliases) > 0 {
		errMsg := fmt.Sprintf("VarArgs contains an alias for command %s: %s", c.Name, v.name)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	if v.isRequired {
		errMsg := fmt.Sprintf("VarArgs marked required for command %s: %s", c.Name, v.name)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	return errs
}

func getFlagNames(f *Flag) []string {
	names := []string{}

	if f.name != "" {
		names = append(names, "--"+f.name)
	}

	for _, alias := range f.aliases {
		names = append(names, "--"+alias)
	}

	if f.shortName != "" {
		names = append(names, "-"+f.shortName)
	}

	return names
}

func validateCommand(c *Command, globalFlagsSeen map[string]struct{}) []error {
	errs := []error{}

	// Check Local Flags

	localFlagNames := []string{}
	for _, f := range c.localFlags.flags {
		localFlagNames = append(localFlagNames, getFlagNames(f)...)
	}

	localFlagSeen := map[string]struct{}{}
	for _, name := range localFlagNames {
		// fmt.Printf("%s, ", name)

		_, found := localFlagSeen[name]
		if found {
			errMsg := fmt.Sprintf("Conflicting local flag names present for command %v: %s", c.Name, name)
			err := errors.New(errMsg)
			errs = append(errs, err)
		}

		_, found = globalFlagsSeen[name]
		if found {
			errMsg := fmt.Sprintf("Conflicting local flag name with global flag present for command %v: %s", c.Name, name)
			err := errors.New(errMsg)
			errs = append(errs, err)
		}

		localFlagSeen[name] = struct{}{}
	}

	for _, f := range c.localFlags.flags {
		errs = append(errs, validateFlag(f, c)...)
	}

	posArgsPresent := len(c.posArgs.optionalPosArgs) > 0 || len(c.posArgs.requiredPosArgs) > 0 || c.varArgs.varArg != nil

	if len(c.subcommands) > 0 && posArgsPresent {
		errMsg := fmt.Sprintf("Command contains both subcommands and posArgs: %v", c.Name)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	if len(c.posArgs.optionalPosArgs) > 0 && c.varArgs.varArg != nil {
		errMsg := fmt.Sprintf("Command contains both optional posArgs and varArgs: %v", c.Name)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	for _, p := range c.posArgs.optionalPosArgs {
		errs = append(errs, validatePosArg(p, c)...)
	}

	for _, p := range c.posArgs.requiredPosArgs {
		errs = append(errs, validatePosArg(p, c)...)
	}

	errs = append(errs, validateVarArgs(c.varArgs.varArg, c)...)

	return errs
}

func (m *Moon) Validate() []error {
	errs := []error{}

	cmdSeen := map[*Command]struct{}{}
	globalFlagSeen := map[string]struct{}{}

	queue := []*Command{m.RootCmd}

	for len(queue) > 0 {
		curGlobalFlagNames := []string{}

		cur := queue[0]
		queue = queue[1:]

		_, found := cmdSeen[cur]
		if found {
			errMsg := fmt.Sprintf("Subcommand loop present: %v", cur)
			err := errors.New(errMsg)
			errs = append(errs, err)

			continue
		}

		cmdSeen[cur] = struct{}{}

		// Check Global Flags

		for _, f := range cur.globalFlags.flags {
			curGlobalFlagNames = append(curGlobalFlagNames, getFlagNames(f)...)
		}

		for _, name := range curGlobalFlagNames {
			_, found := globalFlagSeen[name]
			if found {
				errMsg := fmt.Sprintf("Conflicting global flag name present for command %v: %s", cur.Name, name)
				err := errors.New(errMsg)
				errs = append(errs, err)
			}

			globalFlagSeen[name] = struct{}{}
		}

		errs = append(errs, validateCommand(cur, globalFlagSeen)...)

		for _, sub := range cur.subcommands {
			queue = append(queue, sub)
		}
	}

	return errs
}
