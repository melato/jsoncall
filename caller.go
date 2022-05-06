package jsoncall

import (
	"fmt"
	"reflect"
)

var TraceCalls bool
var TraceInit bool
var TraceDebug bool

type Caller struct {
	Type     reflect.Type
	Methods map[string]*Method
	Names   []*MethodNames
}

func HasReceiver(api reflect.Type) bool {
	if api.Kind() == reflect.Interface {
		return false
	}
	return true
}

func (c *Caller) SetTypePointer(proto interface{}) error {
	pType := reflect.TypeOf(proto)
	switch pType.Kind() {
	case reflect.Interface:
		return fmt.Errorf("cannot use interface{} prototype.  must be pointer or slice")
	case reflect.Pointer:
		eType := pType.Elem()
		if eType.Kind() == reflect.Interface {
			return c.SetType(eType)
		}
	}
	return c.SetType(pType)
}

func (c *Caller) SetType(api reflect.Type) error {
	if api == nil {
		return fmt.Errorf("nil api type")
	}
	hasReceiver := HasReceiver(api)
	switch api.Kind() {
	case reflect.Pointer:
	case reflect.Interface:
	default:
		return fmt.Errorf("unsupported api (%v) kind: %v", api, api.Kind())
	}
	c.Type = api
	n := api.NumMethod()
	if TraceInit {
		fmt.Printf("api type: %v methods: %d\n", api, n)
	}
	namesMap := make(map[string]*MethodNames)
	for _, m := range c.Names {
		namesMap[m.Method] = m
	}
	c.Methods = make(map[string]*Method, n)
	for i := 0; i < n; i++ {
		method := api.Method(i)
		m := newMethod(method, hasReceiver, namesMap[method.Name])
		c.Methods[method.Name] = m
	}
	return nil
}
