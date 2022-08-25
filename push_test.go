package main

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Pilladian/go-helper"
)

// ------------------- Template -------------------
// func TestProcessData_EmptyID(t *testing.T) {
// 	rand.Seed(time.Now().UnixNano())
// 	initialize()

// 	// Input
// 	input := make(map[string]string)
// 	input["id"] = ""
// 	input["data"] = "{\"test\": \"test\"}"

// 	// Expectation
// 	ret_expect_01 := EXPECTATION_01
// 	ret_expect_02 := EXPECTATION_02

// 	// Run test
// 	ret_actual_01, ret_actual_02 := TESTING_PROCEDURE(input["id"], input["data"])

// 	// Evaluate test
// 	if ret_expect_01 == nil && ret_actual_01 == nil {
// 		// OK
// 	} else if ret_expect_01 == nil && ret_actual_01 != nil {
// 		t.Errorf(helper.Prettify("<nil>", ret_actual_01))
// 	} else if ret_expect_01 != nil && ret_actual_01 == nil {
// 		t.Errorf(helper.Prettify(ret_expect_01, "<nil>"))
// 	} else if ret_expect_01 != ret_actual_01 {
// 		t.Errorf(helper.Prettify(ret_expect_01, ret_actual_01))
// 	}

// 	if ret_expect_02 == nil && ret_actual_02 == nil {
// 		// OK
// 	} else if ret_expect_02 == nil && ret_actual_02 != nil {
// 		t.Errorf(helper.Prettify("<nil>", ret_actual_02))
// 	} else if ret_expect_02 != nil && ret_actual_02 == nil {
// 		t.Errorf(helper.Prettify(ret_expect_02, "<nil>"))
// 	} else if ret_expect_02 != ret_actual_02 {
// 		t.Errorf(helper.Prettify(ret_expect_02, ret_actual_02))
// 	}
// }

func TestProcessData_EmptyID(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	initialize()

	// Input
	input := make(map[string]string)
	input["pid"] = ""
	input["data"] = "{\"test\": \"test\"}"

	// Expectation
	ret_expect_01 := 801
	ret_expect_02 := errors.New("Unrecognized pid format \"\"")

	// Run test
	ret_actual_01, ret_actual_02 := processData(input["pid"], input["data"])

	// Evaluate test
	if ret_expect_01 != ret_actual_01 {
		t.Errorf(helper.Prettify(fmt.Sprint(ret_expect_01), fmt.Sprint(ret_actual_01)))
	}

	if ret_expect_02 == nil && ret_actual_02 == nil {
		// OK
	} else if ret_expect_02 == nil && ret_actual_02 != nil {
		t.Errorf(helper.Prettify("<nil>", ret_actual_02.Error()))
	} else if ret_expect_02 != nil && ret_actual_02 == nil {
		t.Errorf(helper.Prettify(ret_expect_02.Error(), "<nil>"))
	} else if ret_expect_02.Error() != ret_actual_02.Error() {
		t.Errorf(helper.Prettify(ret_expect_02.Error(), ret_actual_02.Error()))
	}
}

func TestProcessData_JSONParseError_01(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	initialize()

	// Input
	input := make(map[string]string)
	input["id"] = "test_" + helper.RandomString(10)
	input["data"] = "{\"test\": -}"

	// Expectation
	ret_expect_01 := 807
	ret_expect_02 := errors.New("Provided json data {\"test\": -} could not be parsed : invalid character '}' in numeric literal")

	// Run test
	ret_actual_01, ret_actual_02 := processData(input["id"], input["data"])

	// Evaluate test
	if ret_expect_01 != ret_actual_01 {
		t.Errorf(helper.Prettify(fmt.Sprint(ret_expect_01), fmt.Sprint(ret_actual_01)))
	}

	if ret_expect_02 == nil && ret_actual_02 == nil {
		// OK
	} else if ret_expect_02 == nil && ret_actual_02 != nil {
		t.Errorf(helper.Prettify("<nil>", ret_actual_02.Error()))
	} else if ret_expect_02 != nil && ret_actual_02 == nil {
		t.Errorf(helper.Prettify(ret_expect_02.Error(), "<nil>"))
	} else if ret_expect_02.Error() != ret_actual_02.Error() {
		t.Errorf(helper.Prettify(ret_expect_02.Error(), ret_actual_02.Error()))
	}
}

func TestProcessData_JSONParseError_02(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	initialize()

	// Input
	input := make(map[string]string)
	input["id"] = "test_" + helper.RandomString(10)
	input["data"] = "{\"test\": -, \"test2\": 444}"

	// Expectation
	ret_expect_01 := 807
	ret_expect_02 := errors.New("Provided json data {\"test\": -, \"test2\": 444} could not be parsed : invalid character ',' in numeric literal")

	// Run test
	ret_actual_01, ret_actual_02 := processData(input["id"], input["data"])

	// Evaluate test
	if ret_expect_01 != ret_actual_01 {
		t.Errorf(helper.Prettify(fmt.Sprint(ret_expect_01), fmt.Sprint(ret_actual_01)))
	}

	if ret_expect_02 == nil && ret_actual_02 == nil {
		// OK
	} else if ret_expect_02 == nil && ret_actual_02 != nil {
		t.Errorf(helper.Prettify("<nil>", ret_actual_02.Error()))
	} else if ret_expect_02 != nil && ret_actual_02 == nil {
		t.Errorf(helper.Prettify(ret_expect_02.Error(), "<nil>"))
	} else if ret_expect_02.Error() != ret_actual_02.Error() {
		t.Errorf(helper.Prettify(ret_expect_02.Error(), ret_actual_02.Error()))
	}
}

func TestProcessData_JSONParseError_03(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	initialize()

	// Input
	input := make(map[string]string)
	input["id"] = "test_" + helper.RandomString(10)
	input["data"] = "{\"test\": \"-\", \"test2\": 444}"

	// Expectation
	ret_expect_01 := 0
	var ret_expect_02 error

	// Run test
	ret_actual_01, ret_actual_02 := processData(input["id"], input["data"])

	// Evaluate test
	if ret_expect_01 != ret_actual_01 {
		t.Errorf(helper.Prettify(fmt.Sprint(ret_expect_01), fmt.Sprint(ret_actual_01)))
	}

	if ret_expect_02 == nil && ret_actual_02 == nil {
		// OK
	} else if ret_expect_02 == nil && ret_actual_02 != nil {
		t.Errorf(helper.Prettify("<nil>", ret_actual_02.Error()))
	} else if ret_expect_02 != nil && ret_actual_02 == nil {
		t.Errorf(helper.Prettify(ret_expect_02.Error(), "<nil>"))
	} else if ret_expect_02.Error() != ret_actual_02.Error() {
		t.Errorf(helper.Prettify(ret_expect_02.Error(), ret_actual_02.Error()))
	}
}
