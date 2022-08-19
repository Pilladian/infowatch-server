package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"

	"github.com/Pilladian/go-helper"
)

func validateID(id string) error {
	re, _ := regexp.Compile(` \w+ `)
	if !re.Match([]byte(fmt.Sprintf(" %s ", id))) {
		return errors.New("ID did not match server requirements")
	}
	return nil
}

func validateData(data string) error {
	if err := helper.ValidateSimpleJSON(data); err != nil {
		return err
	}
	return nil
}

func validateSchema(path string, content string) error {
	schema_b, schema_b_err := ioutil.ReadFile(path + "/schema.json")
	if schema_b_err != nil {
		return errors.New("File schema.json could not be opened")
	}
	var schema map[string]interface{}
	json.Unmarshal(schema_b, &schema)

	var content_json map[string]interface{}
	json.Unmarshal([]byte(content), &content_json)

	for key := range content_json {
		if !(reflect.TypeOf((content_json[key])) == reflect.TypeOf(schema[key])) {
			return errors.New("Provided data does not match current schema")
		}
	}
	return nil
}
