package model_test

import (
	"reflect"
	"testing"

	"github.com/jonathangjertsen/toyboy/model"
)

func TestSUB(t *testing.T) {
	type args struct {
		a     uint8
		b     uint8
		carry bool
	}
	tests := []struct {
		name string
		args args
		want model.ALUResult
	}{
		{name: "allzero", args: args{a: 0, b: 0, carry: false}, want: model.ALUResult{
			Value: 0,
			N:     true,
		}},
		{name: "allzero-carry", args: args{a: 0, b: 0, carry: true}, want: model.ALUResult{
			Value: 0xff,
			N:     true,
			C:     true,
			H:     true,
		}},
		{name: "unbricked-fail", args: args{a: 0x35, b: 0x3f, carry: false}, want: model.ALUResult{
			Value: 0xf6,
			N:     true,
			C:     true,
			H:     true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := model.SUB(tt.args.a, tt.args.b, tt.args.carry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SUB() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
