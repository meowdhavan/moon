package moon

import (
	"slices"
	"testing"
)

func TestLongStringFlagParse(t *testing.T) {
	var targetA string
	var targetB string

	c := Command{}
	c.StringFlag(&targetA, "test-flag-a", "")
	c.StringFlag(&targetB, "test-flag-b", "")

	p := newParser(&c, []string{"app", "--test-flag-a", "target_value_1", "--test-flag-b", "target_value_2"})
	p.parseFlags()

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
	c.StringFlag(&targetA, "", "a")
	c.StringFlag(&targetB, "", "b")

	p := newParser(&c, []string{"app", "-a", "target_value_1", "-btarget_value_2"})
	p.parseFlags()

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
	c.AddStringPosArg(&targetA, "a", "", false)
	c.AddStringPosArg(&targetB, "b", "", false)

	p := newParser(&c, []string{"app", "target_value_1", "target_value_2"})
	p.parseFlags()

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
	c.AddStringPosArg(&targetA, "a", "", true)
	c.AddIntPosArg(&targetB, "b", "", false)
	c.AddIntVarLenArg(&targetSlice, "vla", "")

	p := newParser(&c, []string{"app", "target_value_1", "123", "10", "20", "30"})
	p.parseFlags()

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
	c.MultiStringFlag(&targetSlice, "vla", "v")

	p := newParser(&c, []string{"app", "--vla", "a", "-v", "b", "--vla", "c"})
	p.parseFlags()

	wantIntSlice := []string{"a", "b", "c"}
	if !slices.Equal(targetSlice, wantIntSlice) {
		t.Errorf("targetSlice mismatch; got=%v; want %v", targetSlice, wantIntSlice)
	}
}

func TestMultiBoolFlagParse(t *testing.T) {
	var target int

	c := Command{}
	c.MultiBoolFlag(&target, "vla", "v")

	p := newParser(&c, []string{"app", "--vla", "-v", "--vla", "-v", "--vla"})
	p.parseFlags()

	want := 5
	if target != want {
		t.Errorf("target mismatch; got=%v; want %v", target, want)
	}
}

func TestStringFlagDefaultValueParse(t *testing.T) {
	var targetA string

	c := Command{}
	c.StringFlag(&targetA, "test-flag", "", Default("target_value"))

	p := newParser(&c, []string{"app"})
	p.parseFlags()

	var wantString string

	wantString = "target_value"
	if targetA != wantString {
		t.Errorf("targetA mismatch; got=%s, want %s", targetA, wantString)
	}
}
