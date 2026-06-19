package moon

import (
	"errors"
	"testing"
)

func TestDuplicateFlagValidation(t *testing.T) {
	c1 := &Command{Name: "C1"}
	c1.GlobalFlags().String(nil, "test", "", "")

	c2 := &Command{Name: "C2"}
	c2.Flags().String(nil, "test", "", "")
	c1.Subcommand(c2)

	m := NewMoon(c1)

	gotErrs := m.Validate()
	wantErrs := []error{errors.New("Conflicting local flag name with global flag present for command C2: --test")}

	if len(gotErrs) != 1 || gotErrs[0].Error() != wantErrs[0].Error() {
		t.Errorf("errs mismatch; got=%v, want %v", gotErrs, wantErrs)
	}
}
