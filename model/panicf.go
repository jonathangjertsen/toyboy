package model

import "fmt"

func panicf(s string, v ...any) {
	panic(fmt.Sprintf(s, v...))
}

func panicv(v any) {
	panicf("%+v", v)
}
