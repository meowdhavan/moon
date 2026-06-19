package moon

import (
	"errors"
	"fmt"
)

func validateFlag(f *Flag) []error {
	errs := []error{}

	// TODO

	return errs
}

func validatePosArg(p *PosArg) []error {
	errs := []error{}

	// TODO

	return errs
}

func validateVarArgs(v *VarArgs) []error {
	errs := []error{}
	if v == nil {
		return errs
	}

	// TODO

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

	localFlagSeen := map[string]struct{}{}

	localFlagNames := []string{}
	for _, f := range c.localFlags.flags {
		localFlagNames = append(localFlagNames, getFlagNames(f)...)

		for _, name := range localFlagNames {
			_, found := localFlagSeen[name]
			if found {
				errMsg := fmt.Sprintf("Conflicting local flag name present for command %v: %s", c, name)
				err := errors.New(errMsg)
				errs = append(errs, err)
			}

			_, found = globalFlagsSeen[name]
			if found {
				errMsg := fmt.Sprintf("Conflicting flag name with global flag present for command %v: %s", c, name)
				err := errors.New(errMsg)
				errs = append(errs, err)
			}

			localFlagSeen[name] = struct{}{}
		}

		errs = append(errs, validateFlag(f)...)
	}

	posArgsPresent := len(c.posArgs.optionalPosArgs) > 0 || len(c.posArgs.requiredPosArgs) > 0 || c.varArgs.varArg != nil

	if len(c.subcommands) > 0 && posArgsPresent {
		errMsg := fmt.Sprintf("Command contains both subcommands and posArgs: %v", c)
		err := errors.New(errMsg)
		errs = append(errs, err)
	}

	for _, p := range c.posArgs.optionalPosArgs {
		errs = append(errs, validatePosArg(p)...)
	}

	for _, p := range c.posArgs.requiredPosArgs {
		errs = append(errs, validatePosArg(p)...)
	}

	errs = append(errs, validateVarArgs(c.varArgs.varArg)...)

	return errs
}

func (m *Moon) Validate() []error {
	errs := []error{}

	cmdSeen := map[*Command]struct{}{}
	globalFlagSeen := map[string]struct{}{}

	queue := []*Command{m.RootCmd}

	globalFlagNames := []string{}

	for len(queue) > 0 {
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
			globalFlagNames = append(globalFlagNames, getFlagNames(f)...)
		}

		for _, name := range globalFlagNames {
			_, found := globalFlagSeen[name]
			if found {
				errMsg := fmt.Sprintf("Conflicting global flag name present for command %v: %s", cur, name)
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
