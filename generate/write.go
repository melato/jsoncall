// Copyright 2024 Alex Athanasopoulos.
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"fmt"
	"io"
	"reflect"
)

type Writer struct {
	Writer  io.Writer
	Package string
}

func (w *Writer) typeName(rtype reflect.Type) string {
	return TypeName(rtype, w.Package)
}

func (w *Writer) InputName(i int) string {
	return fmt.Sprintf("p%d", 1+i)
}

func (w *Writer) WriteMethodSignature(receiverType string, m reflect.Method, inOffset int) {
	fmt.Fprintf(w.Writer, "func")
	if receiverType != "" {
		fmt.Fprintf(w.Writer, " (t *%s)", receiverType)
	}
	fmt.Fprintf(w.Writer, " %s(", m.Name)
	numIn := m.Type.NumIn()
	for j := inOffset; j < numIn; j++ {
		in := m.Type.In(j)
		if j > inOffset {
			fmt.Fprintf(w.Writer, ", ")
		}
		fmt.Fprintf(w.Writer, "%s %s", w.InputName(j-inOffset), w.typeName(in))
	}
	fmt.Fprintf(w.Writer, ") ")
	numOut := m.Type.NumOut()
	if numOut > 1 {
		fmt.Fprintf(w.Writer, "(")
	}
	for j := 0; j < numOut; j++ {
		out := m.Type.Out(j)
		if j > 0 {
			fmt.Fprintf(w.Writer, ", ")
		}
		fmt.Fprintf(w.Writer, "%s", TypeName(out, w.Package))
	}
	if numOut > 1 {
		fmt.Fprintf(w.Writer, ") ")
	}
}
