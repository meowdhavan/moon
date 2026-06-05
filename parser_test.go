package mon

import (
	"testing"
)

func TestLongStringFlagParse(t *testing.T) {
	var target1 string
	var target2 string

	c := Command{}
	c.AddStringFlag(&target1, []string{"test-flag-1"}, "", "", false)
	c.AddStringFlag(&target2, []string{"test-flag-2"}, "", "", false)

	p := newParser([]string{"app", "--test-flag-1", "target_value_1", "--test-flag-2", "target_value_2"})

	p.parseFlags(&c)

	if target1 != "target_value_1" {
		t.Errorf("target1=%s; want %s", target1, "target_value_1")
	}

	if target2 != "target_value_2" {
		t.Errorf("target2=%s; want %s", target2, "target_value_2")
	}
}

func TestShortStringFlagParse(t *testing.T) {
	var targetA string
	var targetB string

	c := Command{}
	c.AddStringFlag(&targetA, []string{}, "a", "", false)
	c.AddStringFlag(&targetB, []string{}, "b", "", false)

	p := newParser([]string{"app", "-a", "target_value_1", "-btarget_value_2"})

	p.parseFlags(&c)

	if targetA != "target_value_1" {
		t.Errorf("targetA=%s; want %s", targetA, "target_value_1")
	}

	if targetB != "target_value_2" {
		t.Errorf("targetB=%s; want %s", targetB, "target_value_2")
	}
}
