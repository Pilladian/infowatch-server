package main

import (
	"errors"
	"fmt"
	"regexp"
)

func validateID(id string) error {
	re, _ := regexp.Compile(` \w+ `)
	if !re.Match([]byte(fmt.Sprintf(" %s ", id))) {
		return errors.New("ID did not match server requirements")
	}
	return nil
}

func validateData(data string) error {
	re, _ := regexp.Compile(` \{((\d+|"\w+") *: *(\d+|"\w+") *, *)*(\d+|"\w+") *: *(\d+|"\w+")\} `)
	if !re.Match([]byte(fmt.Sprintf(" %s ", data))) {
		return errors.New("data did not match json format or contained invalid characters")
	}
	return nil
}
