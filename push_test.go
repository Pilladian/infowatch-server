package main

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestProcessData_EmptyID(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	initialize()

	got, err_got := processData("", "{\"test\": \"test\"}")
	want := 805
	err_want := errors.New("Unrecognized ID format \"\"")

	if got != want {
		t.Errorf("[ WANTED ] %d\n\t     [ ACTUAL ] %d\n\n", want, got)
	}
	if err_got == nil {
		t.Errorf("[ WANTED ] %s\n\t     [ ACTUAL ] <nil>", err_want.Error())
	} else if err_got.Error() != err_want.Error() {
		t.Errorf("[ WANTED ] %s\n\t     [ ACTUAL ] %s", err_want.Error(), err_got.Error())
	}
}

func TestProcessData_UnknownJsonType(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	initialize()

	got, err_got := processData("test_"+randomString(10), "{\"test\": -}")
	want := 806
	err_want := errors.New("Unknown type for json data")

	if got != want {
		t.Errorf("[ WANTED ] %d\n\t     [ ACTUAL ] %d\n\n", want, got)
	}
	if err_got == nil {
		t.Errorf("[ WANTED ] %s\n\t     [ ACTUAL ] <nil>", err_want.Error())
	} else if err_got.Error() != err_want.Error() {
		t.Errorf("[ WANTED ] %s\n\t     [ ACTUAL ] %s", err_want.Error(), err_got.Error())
	}
}
