package generate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"

	"melato.org/jsoncall"
)

// Generator - generates client Go code for a web service that uses the jsoncall conventions
type Generator struct {
	Package            string
	Type               string
	Imports            []string
	InternalTypePrefix string
	OutputFile         string
	inOffset           int
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
		fmt.Printf("generating %s\n", g.OutputFile)
		return os.WriteFile(g.OutputFile, data, os.FileMode(0644))
	} else {
		fmt.Printf("%s\n", string(data))
	}
	return nil
}

var pathRegexp = regexp.MustCompile(`(.*?)([a-z_A-Z]+)\.([^\.]+)`)

func TypeName(rtype reflect.Type, pkg string) string {
	s := rtype.String()
	parts := pathRegexp.FindStringSubmatch(s)
	if len(parts) > 0 && parts[2] == pkg {
		return parts[1] + parts[3]
	}
	return s
}

func (g *Generator) typeName(rtype reflect.Type) string {
	return TypeName(rtype, g.Package)
}

func (g *Generator) writeMethodHeader(w io.Writer, m reflect.Method) {
	fmt.Fprintf(w, "\nfunc (t *%s) %s(", g.Type, m.Name)
	numIn := m.Type.NumIn()
	for j := g.inOffset; j < numIn; j++ {
		in := m.Type.In(j)
		if j > g.inOffset {
			fmt.Fprintf(w, ", ")
		}
		fmt.Fprintf(w, "p%d %s", 1+j-g.inOffset, g.typeName(in))
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
		fmt.Fprintf(w, "%s", g.typeName(out))
	}
	if numOut > 1 {
		fmt.Fprintf(w, ") ")
	}
	fmt.Fprintf(w, "{\n")
}

func (g *Generator) writeMethodInputs(w io.Writer, m reflect.Method) {
	numIn := m.Type.NumIn()
	for j := g.inOffset; j < numIn; j++ {
		fmt.Fprintf(w, ", p%d", 1+j-g.inOffset)
	}
	fmt.Fprintf(w, ")\n")
}

type Field struct {
	Index    int
	Name     string
	Type     string
	JsonName string
}

func (g *Generator) GetOutputFields(m reflect.Method, desc *jsoncall.MethodDescriptor) []Field {
	var fields []Field
	var errorp *error
	errorType := reflect.TypeOf(errorp).Elem()
	numOut := m.Type.NumOut()
	if numOut > 0 {
		for j := 0; j < numOut; j++ {
			out := m.Type.Out(j)
			if out != errorType {
				fields = append(fields, Field{
					Index:    j,
					Name:     fmt.Sprintf("P%d", j+1),
					Type:     g.typeName(out),
					JsonName: desc.Out[j],
				})
			}
		}
	}
	return fields
}

func (g *Generator) generateMethodStruct(w io.Writer, m reflect.Method, desc *jsoncall.MethodDescriptor) {
	fields := g.GetOutputFields(m, desc)
	var errorp *error
	errorType := reflect.TypeOf(errorp).Elem()
	numOut := m.Type.NumOut()
	structName := g.InternalTypePrefix + m.Name
	if len(fields) > 0 {
		fmt.Fprintf(w, "\ntype %s struct {\n", structName)
		for _, f := range fields {
			fmt.Fprintf(w, "  %s %s `json:\"%s\"`\n", f.Name, f.Type, f.JsonName)
		}
		fmt.Fprintf(w, "}\n")
	}
	g.writeMethodHeader(w, m)
	var hasError bool
	for j := 0; j < numOut; j++ {
		out := m.Type.Out(j)
		if out == errorType {
			hasError = true
		}
	}
	var errAssign string
	if hasError {
		errAssign = "err := "
	}
	if len(fields) > 0 {
		fmt.Fprintf(w, "  var out %s\n", structName)
		fmt.Fprintf(w, "  %st.Client.Call(&out, \"%s\"", errAssign, m.Name)
	} else {
		fmt.Fprintf(w, "  %st.Client.Call(nil, \"%s\"", errAssign, m.Name)
	}
	g.writeMethodInputs(w, m)

	if numOut > 0 {
		fmt.Fprintf(w, "  return ")
		fieldNames := make(map[int]string)
		for _, f := range fields {
			fieldNames[f.Index] = f.Name
		}
		for j := 0; j < numOut; j++ {
			if j > 0 {
				fmt.Fprintf(w, ",")
			}
			out := m.Type.Out(j)
			if out == errorType {
				fmt.Fprintf(w, " err")
			} else {
				fmt.Fprintf(w, " out.%s", fieldNames[j])
			}
		}
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "}\n")
}

// GenerateType - Generate a type that implements the methods of type <t>
func (g *Generator) GenerateClient(c *jsoncall.Caller) ([]byte, error) {
	t := c.Type()
	if jsoncall.HasReceiver(t) {
		g.inOffset = 1
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
	fmt.Fprintf(w, "  Client   jsoncall.Client\n")
	fmt.Fprintf(w, "}\n")
	descMap := make(map[string]*jsoncall.MethodDescriptor)
	for _, m := range c.Desc {
		descMap[m.Method] = m
	}
	n := t.NumMethod()
	for i := 0; i < n; i++ {
		m := t.Method(i)
		if !m.IsExported() {
			continue
		}
		desc := descMap[m.Name]
		if desc == nil {
			desc = jsoncall.DefaultMethodDescriptor(m, false)
		}
		g.generateMethodStruct(w, m, desc)
	}
	return w.Bytes(), nil
}
