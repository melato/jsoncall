package generate

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
)

// Generator - generates client Go code for a web service that uses the jsoncall conventions
type Generator struct {
	Package    string
	Type       string
	Imports    []string
	OutputFile string
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

// GenerateType - Generate a type that implements the methods of type <t>
func (g *Generator) GenerateType(t reflect.Type) ([]byte, error) {
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
	n := t.NumMethod()
	var errorp *error
	errorType := reflect.TypeOf(errorp).Elem()
	for i := 0; i < n; i++ {
		m := t.Method(i)
		if !m.IsExported() {
			continue
		}
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
		fmt.Fprintf(w, "  result := t.Client.Call(\"%s\")\n", m.Name)

		for j := 0; j < numOut; j++ {
			out := m.Type.Out(j)
			if out == errorType {
				fmt.Fprintf(w, `  var x%d error
  if result[%d] != nil {
	 x%d = result[%d].(error)
  }
`, j, j, j, j)
			} else {
				fmt.Fprintf(w, "  var x%d %s = result[%d].(%s)\n", j, out.String(), j, out.String())
			}
		}
		fmt.Fprintf(w, "  return")
		for j := 0; j < numOut; j++ {
			if j > 0 {
				fmt.Fprintf(w, ",")
			}
			fmt.Fprintf(w, " x%d", j)
		}
		fmt.Fprintf(w, "\n}\n")
	}
	return w.Bytes(), nil
}
