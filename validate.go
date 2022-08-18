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
