package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"melato.org/jsoncall"
)

type Names []*jsoncall.MethodNames

func GenerateMethodNames(v interface{}) Names {
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
	var existingNames Names
	err = json.Unmarshal(data, &existingNames)
	if err != nil {
		return err
	}
	err = names.Merge(existingNames)
	if err != nil {
		return err
	}

	return names.Write(file)
}

// Merge overrides a with b
// Typically a would be the default method names, and b would be user-defined
// If a method exists in both a and b, it uses the specification from b
// If the two specifications of a method have different number of inputs or outputs, it returns an error
// If a method exists in b but not in a, it is ignored
func (a Names) Merge(b Names) error {
	fmt.Printf("merging %d with %d names\n", len(a), len(b))
	bmap := make(map[string]*jsoncall.MethodNames)
	for _, names := range b {
		bmap[names.Method] = names
	}
	for _, x := range a {
		y, exists := bmap[x.Method]
		if exists {
			if len(x.In) != len(y.In) {
				return fmt.Errorf("the inputs of method %s have changed.  Please update them.", x.Method)
			}
			if len(x.Out) != len(y.Out) {
				return fmt.Errorf("the outputs of method %s have changed.  Please update them.", x.Method)
			}
			x.In = y.In
			x.Out = y.Out
			x.Path = y.Path
		}
	}
	return nil
}

func (names Names) Write(file string) error {
	data, err := json.MarshalIndent(names, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, os.FileMode(0645))
}
