package moon

import (
	"slices"
	"testing"
)

func TestLongStringFlagParse(t *testing.T) {
	var targetA string
	var targetB string

	c := Command{}
	c.Flags().String(&targetA, "test-flag-a", "", "")
	c.Flags().String(&targetB, "test-flag-b", "", "")

	p := newParser(&c, []string{"app", "--test-flag-a", "target_value_1", "--test-flag-b", "target_value_2"})
	p.parse()

	var wantString string

	wantString = "target_value_1"
	if targetA != wantString {
		t.Errorf("targetA mismatch; got=%s, want %s", targetA, wantString)
	}

	wantString = "target_value_2"
	if targetB != wantString {
		t.Errorf("targetB mismatch; got=%s, want %s", targetB, wantString)
	}
}

func TestShortStringFlagParse(t *testing.T) {
	var targetA string
	var targetB string

	c := Command{}
	c.Flags().String(&targetA, "", "a", "")
	c.Flags().String(&targetB, "", "b", "")

	p := newParser(&c, []string{"app", "-a", "target_value_1", "-btarget_value_2"})
	p.parse()

	var wantString string

	wantString = "target_value_1"
	if targetA != wantString {
		t.Errorf("targetA mismatch; got=%s, want %s", targetA, wantString)
	}

	wantString = "target_value_2"
	if targetB != wantString {
		t.Errorf("targetB mismatch; got=%s, want %s", targetB, wantString)
	}
}

func TestStringPosArgParse(t *testing.T) {
	var targetA string
	var targetB string

	c := Command{}
	c.PosArgs().String(&targetA, "a", "")
	c.PosArgs().String(&targetB, "b", "")

	p := newParser(&c, []string{"app", "target_value_1", "target_value_2"})
	p.parse()

	var wantString string

	wantString = "target_value_1"
	if targetA != wantString {
		t.Errorf("targetA mismatch; got=%s, want %s", targetA, wantString)
	}

	wantString = "target_value_2"
	if targetB != wantString {
		t.Errorf("targetB mismatch; got=%s, want %s", targetB, wantString)
	}
}

func TestMultitypePosArgParse(t *testing.T) {
	var targetA string
	var targetB int
	var targetSlice []int

	c := Command{}
	c.PosArgs().String(&targetA, "a", "", Required())
	c.PosArgs().Int(&targetB, "b", "")
	c.VarArgs().Int(&targetSlice, "vla", "")

	p := newParser(&c, []string{"app", "target_value_1", "123", "10", "20", "30"})
	p.parse()

	gotString := "target_value_1"
	if targetA != gotString {
		t.Errorf("targetA mismatch; got=%s; want %s", targetA, gotString)
	}

	gotInt := 123
	if targetB != gotInt {
		t.Errorf("targetB mismatch; got=%d; want %d", targetB, gotInt)
	}

	wantIntSlice := []int{10, 20, 30}
	if !slices.Equal(targetSlice, wantIntSlice) {
		t.Errorf("targetSlice mismatch; got=%v; want %v", targetSlice, wantIntSlice)
	}
}

func TestMultiStringFlagParse(t *testing.T) {
	var targetSlice []string

	c := Command{}
	c.Flags().MultiString(&targetSlice, "vla", "v", "")

	p := newParser(&c, []string{"app", "--vla", "a", "-v", "b", "--vla", "c"})
	p.parse()

	wantIntSlice := []string{"a", "b", "c"}
	if !slices.Equal(targetSlice, wantIntSlice) {
		t.Errorf("targetSlice mismatch; got=%v; want %v", targetSlice, wantIntSlice)
	}
}

func TestMultiBoolFlagParse(t *testing.T) {
	var target int

	c := Command{}
	c.Flags().MultiBool(&target, "vla", "v", "")

	p := newParser(&c, []string{"app", "--vla", "-v", "--vla", "-v", "--vla"})
	p.parse()

	want := 5
	if target != want {
		t.Errorf("target mismatch; got=%v; want %v", target, want)
	}
}

func TestStringFlagDefaultValueParse(t *testing.T) {
	var targetA string

	c := Command{}
	c.Flags().String(&targetA, "test-flag", "", "", Default("target_value"))

	p := newParser(&c, []string{"app"})
	p.parse()

	var wantString string

	wantString = "target_value"
	if targetA != wantString {
		t.Errorf("targetA mismatch; got=%s, want %s", targetA, wantString)
	}
}

func TestSubcommandParse(t *testing.T) {
	rootCmd := Command{Name: "root"}
	subCmd := Command{Name: "sub"}

	rootCmd.Subcommand(&subCmd)

	p := newParser(&rootCmd, []string{"app", "sub"})
	p.parse()

	got := p.currentCmd.Name
	want := subCmd.Name

	if got != want {
		t.Errorf("subCmd mismatch; got=%s, want %s", got, want)
	}
}

func TestInvalidSubcommandParse(t *testing.T) {
	rootCmd := Command{Name: "root"}
	subCmd := Command{Name: "sub"}

	rootCmd.Subcommand(&subCmd)

	p := newParser(&rootCmd, []string{"app", "incorrect", "sub"})
	p.parse()

	gotName := p.currentCmd.Name
	wantName := rootCmd.Name

	if gotName != wantName {
		t.Errorf("subCmd mismatch; got=%s, want %s", gotName, wantName)
	}

	if !p.unrecognizedSubcommand {
		t.Errorf("p.unrecognizedSubcommand mismatch; got=%v, want %v", p.unrecognizedSubcommand, true)
	}
}
