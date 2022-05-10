package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"melato.org/jsoncall"
)

func GenerateMethodNames(v interface{}) jsoncall.Names {
	t := reflect.TypeOf(v).Elem()
	hasReceiver := jsoncall.HasReceiver(t)
	var result []*jsoncall.MethodNames
	numMethods := t.NumMethod()
	for i := 0; i < numMethods; i++ {
		m := t.Method(i)
		if m.IsExported() {
			result = append(result, jsoncall.DefaultMethodNames(m, hasReceiver))
		}
	}
	return result
}

func ReadNames(file string) (jsoncall.Names, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return jsoncall.ParseNames(data)
}

func WriteNames(names jsoncall.Names, file string) error {
	data, err := json.MarshalIndent(names, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, os.FileMode(0645))
}

func UpdateMethodNames(v interface{}, file string) error {
	names := GenerateMethodNames(v)
	_, err := os.Stat(file)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var existingNames jsoncall.Names
	err = json.Unmarshal(data, &existingNames)
	if err != nil {
		return err
	}
	fmt.Printf("%s: merging %d names with %d\n", file, len(names), len(existingNames))
	err = names.Merge(existingNames)
	if err != nil {
		return err
	}

	return WriteNames(names, file)
}
