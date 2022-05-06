package generate

import (
	"encoding/json"
	"os"
	"reflect"

	"melato.org/jsoncall"
)

func GenerateMethodNames(v interface{}) []*jsoncall.MethodNames {
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

func WriteMethodNamesJSON(names []*jsoncall.MethodNames, file string) error {
	data, err := json.MarshalIndent(names, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, os.FileMode(0645))
}
