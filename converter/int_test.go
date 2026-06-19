package converter

import (
	"testing"
)

func TestIntConverter(t *testing.T) {
	var got int
	var want int
	var err error

	got, err = ToInt("13")
	want = 13

	if got != want || err != nil {
		t.Errorf("TestIntConverter (valid) mismatch; got=%v, want=%v, err=%v", got, want, err)
	}

	got, err = ToInt("gg")
	want = 0
	wantErrMsg := "Value cannot be converted into an integer: gg"

	if err == nil || err.Error() != wantErrMsg {
		t.Errorf("TestIntConverter (invalid) mismatch; err=%v, wantErr=%v", err, wantErrMsg)
	}
}
