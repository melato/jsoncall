package jsoncall

import (
	"fmt"
	"reflect"
)

var TraceCalls bool
var TraceInit bool
var TraceDebug bool

// Caller specifies a set of methods and how to call them.
type Caller struct {
	// Names determines what methods are used for marshalling/unmarshalling each method call.
	// If any method's names are missing, default names are used.
	// If you set this, you must do set it before calling SetType or SetTypePointer
	Names Names

	rType   reflect.Type
	methods map[string]*Method
}

func (c *Caller) Type() reflect.Type {
	return c.rType
}

func (c *Caller) SetNames(names Names) {
	c.Names = names
}

func (c *Caller) SetNamesJson(data []byte) error {
	if len(data) == 0 {
		c.SetNames(nil)
		return nil
	}
	names, err := ParseNames(data)
	if err != nil {
		return err
	}
	c.SetNames(names)
	return nil
}

// HasReceiver determines whether the methods of a given Type include a receiver first argument
// It returns false for Interface types, true for other types
func HasReceiver(rtype reflect.Type) bool {
	if rtype.Kind() == reflect.Interface {
		return false
	}
	return true
}

// SetTypePointer is similar to SetType, but infers the type from a provided prototype pointer,
// which can be a pointer to an interface type or a struct type.
//
// Examples:
// 	SetTypePointer((*Demo)(nil))
// or:
// 	var v *Demo
// 	SetTypePointer(v)
// instead of:
// 	var v *Demo
// 	SetType(reflect.TypeOf(v))
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

// SetType specifies the Type whose methods are to be used by the Caller
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
	c.rType = api
	n := api.NumMethod()
	if TraceInit {
		fmt.Printf("api type: %v methods: %d\n", api, n)
	}
	namesMap := make(map[string]*MethodNames)
	for _, m := range c.Names {
		namesMap[m.Method] = m
	}
	c.methods = make(map[string]*Method, n)
	for i := 0; i < n; i++ {
		method := api.Method(i)
		m, err := newMethod(method, hasReceiver, namesMap[method.Name])
		if err != nil {
			return err
		}
		c.methods[method.Name] = m
	}
	return nil
}
