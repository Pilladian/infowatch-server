package infowatch_server

import (
	"errors"
	"testing"
)

func TestProcessData_EmptyID(t *testing.T) {
	got, err_got := processData("", "{'test': 'test'}")
	want := 805
	err_want := errors.New("Unrecognized ID format \"\"")

	if got != want {
		t.Errorf("[ WANTED ] %d\n\t     [ ACTUAL ] %d\n\n", want, got)
	}
	if err_got.Error() != err_want.Error() {
		t.Errorf("[ WANTED ] %s\n\t     [ ACTUAL ] %s", err_want.Error(), err_got.Error())
	}
}
