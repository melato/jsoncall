package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"melato.org/jsoncall"
)

func ReadDescriptor(file string) (jsoncall.ApiDescriptor, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return jsoncall.ParseApiDescriptor(data)
}

func WriteDescriptor(desc jsoncall.ApiDescriptor, file string) error {
	data, err := json.MarshalIndent(desc, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, os.FileMode(0645))
}

func GenerateDescriptorT(t reflect.Type) jsoncall.ApiDescriptor {
	hasReceiver := jsoncall.HasReceiver(t)
	var result []*jsoncall.MethodDescriptor
	numMethods := t.NumMethod()
	for i := 0; i < numMethods; i++ {
		m := t.Method(i)
		if m.IsExported() {
			result = append(result, jsoncall.DefaultMethodDescriptor(m, hasReceiver))
		}
	}
	return result
}

func GenerateDescriptor(v interface{}) jsoncall.ApiDescriptor {
	t := reflect.TypeOf(v).Elem()
	return GenerateDescriptorT(t)
}

func UpdateDescriptorT(t reflect.Type, file string) error {
	desc := GenerateDescriptorT(t)
	_, err := os.Stat(file)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var existingDescriptor jsoncall.ApiDescriptor
	err = json.Unmarshal(data, &existingDescriptor)
	if err != nil {
		return err
	}
	fmt.Printf("%s: merging %d methods with %d\n", file, len(desc), len(existingDescriptor))
	err = desc.Merge(existingDescriptor)
	if err != nil {
		return err
	}

	return WriteDescriptor(desc, file)
}

func UpdateDescriptor(v interface{}, file string) error {
	t := reflect.TypeOf(v).Elem()
	return UpdateDescriptorT(t, file)
}
