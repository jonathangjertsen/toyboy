package model

import (
	"reflect"
	"strings"
	"unicode"
)

type stackT struct {
	buf []string
	l   int
}

func (st *stackT) push(x string) {
	st.buf = append(st.buf, x)
	st.l++
}

func (st *stackT) pop() {
	st.l--
	st.buf = st.buf[:st.l]
}

func (st *stackT) join() string {
	return strings.Join(st.buf, ".")[1:]
}

func mustTriviallySerialize(v any) {
	t := reflect.TypeOf(v)
	var stack stackT
	stack.buf = make([]string, 0, 10)
	stack.push(t.Name())
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("not a struct")
	}
	checkType(t, &stack)
}

func checkType(t reflect.Type, stack *stackT) {
	stack.push(t.Name())
	defer func() {
		stack.pop()
	}()
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !unicode.IsUpper(rune(f.Name[0])) {
				panic("unexported field in " + stack.join() + ": " + f.Name)
			}
			if f.Type.Kind() == reflect.Ptr {
				elem := f.Type.Elem()
				if elem.Kind() != reflect.Slice && elem.Kind() != reflect.Map {
					panic("invalid pointer field in " + stack.join() + ": " + f.Name)
				}
			}
			checkType(f.Type, stack)
		}
	case reflect.Array, reflect.Slice:
		checkType(t.Elem(), stack)
	case reflect.Map:
		checkType(t.Key(), stack)
		checkType(t.Elem(), stack)
	case reflect.Interface, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panic("unsupported type: " + t.Kind().String())
	}
}
