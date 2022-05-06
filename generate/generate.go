package generate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"

	"melato.org/jsoncall"
)

// Generator - generates client Go code for a web service that uses the jsoncall conventions
type Generator struct {
	Package            string
	Type               string
	Imports            []string
	InternalTypePrefix string
	OutputFile         string
	Names              []*jsoncall.MethodNames `name:"-"`
}

func (t *Generator) Init() error {
	t.InternalTypePrefix = "r"
	return nil
}

func (t *Generator) Configured() error {
	if t.Package == "" {
		return fmt.Errorf("missing package")
	}
	if t.Type == "" {
		return fmt.Errorf("missing type name")
	}
	return nil
}

func (g *Generator) Output(data []byte, err error) error {
	if err != nil {
		return err
	}
	if g.OutputFile != "" {
		return os.WriteFile(g.OutputFile, data, os.FileMode(0644))
	} else {
		fmt.Printf("%s\n", string(data))
	}
	return nil
}

// GenerateP - Same as GenerateType(reflect.TypeOf(v).Elem())
func (g *Generator) GenerateP(v interface{}) ([]byte, error) {
	t := reflect.TypeOf(v).Elem()
	return g.GenerateType(t)
}

func (g *Generator) writeMethodHeader(w io.Writer, m reflect.Method) {
	fmt.Fprintf(w, "\nfunc (t *%s) %s(", g.Type, m.Name)
	numIn := m.Type.NumIn()
	for j := 0; j < numIn; j++ {
		in := m.Type.In(j)
		if j > 0 {
			fmt.Fprintf(w, ", ")
		}
		fmt.Fprintf(w, "p%d %s", j, in.String())
	}
	fmt.Fprintf(w, ") ")
	numOut := m.Type.NumOut()
	if numOut > 1 {
		fmt.Fprintf(w, "(")
	}
	for j := 0; j < numOut; j++ {
		out := m.Type.Out(j)
		if j > 0 {
			fmt.Fprintf(w, ", ")
		}
		fmt.Fprintf(w, "%s", out.String())
	}
	if numOut > 1 {
		fmt.Fprintf(w, ") ")
	}
	fmt.Fprintf(w, "{\n")
}

func (g *Generator) writeMethodInputs(w io.Writer, m reflect.Method) {
	numIn := m.Type.NumIn()
	for j := 0; j < numIn; j++ {
		fmt.Fprintf(w, ", p%d", j)
	}
	fmt.Fprintf(w, ")\n")
}

func (g *Generator) generateMethodStruct(w io.Writer, m reflect.Method, names *jsoncall.MethodNames) {
	var errorp *error
	errorType := reflect.TypeOf(errorp).Elem()
	numOut := m.Type.NumOut()
	structType := g.InternalTypePrefix + m.Name
	if numOut > 0 {
		fmt.Fprintf(w, "\ntype %s struct {\n", structType)
		for j := 0; j < numOut; j++ {
			out := m.Type.Out(j)
			outType := out.String()
			if out != errorType {
				fmt.Fprintf(w, "  P%d %s `json:\"%s\"`\n", j+1, outType, names.Out[j])
			}
			//fmt.Fprintf(w, "  P%d *jsoncall.Error\n", j+1)
		}
		fmt.Fprintf(w, "}\n")
	}
	g.writeMethodHeader(w, m)
	fmt.Fprintf(w, "  var out %s\n", structType)
	fmt.Fprintf(w, "  err := t.Client.CallV(&out, \"%s\"", m.Name)
	g.writeMethodInputs(w, m)

	fmt.Fprintf(w, "  return ")
	for j := 0; j < numOut; j++ {
		if j > 0 {
			fmt.Fprintf(w, ",")
		}
		out := m.Type.Out(j)
		if out == errorType {
			fmt.Fprintf(w, " err")
		} else {
			fmt.Fprintf(w, " out.P%d", j+1)
		}
	}
	fmt.Fprintf(w, "\n}\n")
}

// GenerateType - Generate a type that implements the methods of type <t>
func (g *Generator) GenerateType(t reflect.Type) ([]byte, error) {
	if t.Kind() != reflect.Interface {
		return nil, fmt.Errorf("please use an interface type")
	}
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "package %s\n\n", g.Package)
	fmt.Fprintf(w, "import (\n")
	fmt.Fprintf(w, "  \"melato.org/jsoncall\"\n")
	for _, s := range g.Imports {
		fmt.Fprintf(w, "  \"%s\"\n", s)
	}
	fmt.Fprintf(w, ")\n\n")
	fmt.Fprintf(w, "// %s - Generated client for %s\n", g.Type, t.String())
	fmt.Fprintf(w, "type %s struct {\n", g.Type)
	fmt.Fprintf(w, "  Client   *jsoncall.Client\n")
	fmt.Fprintf(w, "}\n")
	namesMap := make(map[string]*jsoncall.MethodNames)
	for _, m := range g.Names {
		namesMap[m.Method] = m
	}
	n := t.NumMethod()
	for i := 0; i < n; i++ {
		m := t.Method(i)
		if !m.IsExported() {
			continue
		}
		names := namesMap[m.Name]
		if names == nil {
			names = jsoncall.DefaultMethodNames(m, false)
		}
		g.generateMethodStruct(w, m, names)
	}
	return w.Bytes(), nil
}
