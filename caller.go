package jsoncall

import (
	"fmt"
	"reflect"
)

var TraceCalls bool
var TraceInit bool

// Caller specifies a set of methods and how to call them.
type Caller struct {
	// Desc determines what methods are used for marshalling/unmarshalling each method call.
	// If any method descriptor is missing, a default descriptor are used.
	// If you set this, you must do set it before calling SetType or SetTypePointer
	Desc ApiDescriptor

	Prefix string

	rType   reflect.Type
	methods map[string]*Method
}

func (c *Caller) Type() reflect.Type {
	return c.rType
}

func (c *Caller) SetDescriptor(desc ApiDescriptor) {
	c.Desc = desc
}

func (c *Caller) SetDescriptorJson(data []byte) error {
	if len(data) == 0 {
		c.SetDescriptor(nil)
		return nil
	}
	desc, err := ParseApiDescriptor(data)
	if err != nil {
		return err
	}
	c.SetDescriptor(desc)
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
	descMap := make(map[string]*MethodDescriptor)
	for _, m := range c.Desc {
		descMap[m.Method] = m
	}
	c.methods = make(map[string]*Method, n)
	for i := 0; i < n; i++ {
		method := api.Method(i)
		m, err := newMethod(method, hasReceiver, descMap[method.Name])
		if err != nil {
			return err
		}
		c.methods[method.Name] = m
	}
	return nil
}

func (c *Caller) MethodByName(name string) *Method {
	return c.methods[name]
}

// NewCaller creates a new caller for the prototype methods
// descJson is an optional JSON representation of ApiDescriptor
func NewCaller(proto interface{}, descJson []byte) (*Caller, error) {
	var c Caller
	err := c.SetDescriptorJson(descJson)
	if err != nil {
		return nil, err
	}
	err = c.SetTypePointer(proto)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
