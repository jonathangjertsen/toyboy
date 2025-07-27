package model_test

import (
	"testing"

	"github.com/jonathangjertsen/toyboy/model"
)

func BenchmarkBitMethod(b *testing.B) {
	for range b.N {
		for i := range 256 {
			for j := range uint(8) {
				a := model.Data8(i)
				_ = a.Bit(j)
			}
		}
	}
}
