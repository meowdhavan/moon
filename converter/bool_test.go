package converter

import (
	"testing"
)

func TestBoolConverter(t *testing.T) {
	var got bool
	var want bool
	var err error

	got, err = ToBool("false")
	want = false

	if got != want || err != nil {
		t.Errorf("TestBoolConverter (false 1) mismatch; got=%v, want=%v, err=%v", got, want, err)
	}

	got, err = ToBool("no")
	want = false

	if got != want || err != nil {
		t.Errorf("TestBoolConverter (false 2) mismatch; got=%v, want=%v, err=%v", got, want, err)
	}

	got, err = ToBool("yes")
	want = false

	if got != want || err != nil {
		t.Errorf("TestBoolConverter (false 3) mismatch; got=%v, want=%v, err=%v", got, want, err)
	}

	got, err = ToBool("true")
	want = true

	if got != want || err != nil {
		t.Errorf("TestBoolConverter (true 1) mismatch; got=%v, want=%v, err=%v", got, want, err)
	}

	got, err = ToBool("TRUE")
	want = true

	if got != want || err != nil {
		t.Errorf("TestBoolConverter (true 2) mismatch; got=%v, want=%v, err=%v", got, want, err)
	}

	got, err = ToBool("tRuE")
	want = true

	if got != want || err != nil {
		t.Errorf("TestBoolConverter (true 3) mismatch; got=%v, want=%v, err=%v", got, want, err)
	}
}
