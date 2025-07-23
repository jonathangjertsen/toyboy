package model_test

import (
	"testing"

	"github.com/jonathangjertsen/toyboy/model"
)

func TestDAA(t *testing.T) {
	if have, want := model.DAA(0x00, false, false, false), (model.ALUResult{Value: 0}); have != want {
		t.Errorf("want %+v have %+v", want, have)
	}
	if have, want := model.DAA(0x0a, false, false, true), (model.ALUResult{Value: 0x10}); have != want {
		t.Errorf("want %+v have %+v", want, have)
	}
	if have, want := model.DAA(0x0f, false, true, true), (model.ALUResult{Value: 0x09, N: true}); have != want {
		t.Errorf("want %+v have %+v", want, have)
	}
	if have, want := model.DAA(0xff, false, true, true), (model.ALUResult{Value: 0x99, N: true}); have != want {
		t.Errorf("want %+v have %+v", want, have)
	}
}
