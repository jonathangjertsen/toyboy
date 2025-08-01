package model

import (
	"reflect"
	"unicode"
)

func mustTriviallySerialize(v any) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("not a struct")
	}
	checkType(t)
}

func checkType(t reflect.Type) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !unicode.IsUpper(rune(f.Name[0])) {
				panic("unexported field: " + f.Name)
			}
			if f.Type.Kind() == reflect.Ptr {
				elem := f.Type.Elem()
				if elem.Kind() != reflect.Slice && elem.Kind() != reflect.Map {
					panic("invalid pointer field: " + f.Name)
				}
			}
			checkType(f.Type)
		}
	case reflect.Array, reflect.Slice:
		checkType(t.Elem())
	case reflect.Map:
		checkType(t.Key())
		checkType(t.Elem())
	case reflect.Interface, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panic("unsupported type: " + t.Kind().String())
	}
}
