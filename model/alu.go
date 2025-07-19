package model

import "fmt"

type Bit uint8

func (b Bit) Bool() bool {
	switch b {
	case 0:
		return false
	case 1:
		return true
	}
	panicf("unknown value for Bit: %v", b)
	return false
}

func NewBit(b bool) Bit {
	if b {
		return 1
	}
	return 0
}

type Bits [8]Bit

func (bi Bits) Pack() uint8 {
	var out uint8
	for i := range 8 {
		out |= uint8(bi[i]) << i
	}
	return out
}

type ALUResult struct {
	Value uint8
	Z0    bool
	N     bool
	H     bool
	C     bool
}

func (res ALUResult) Z() bool {
	if res.Z0 {
		return false
	}
	return res.Value == 0
}

func Unpack(b8 uint8) Bits {
	var out Bits
	for i := range 8 {
		out[i] = Bit((b8 & (1 << i)) >> i)
	}
	return out
}

func ADD(a, b uint8, carryIn bool) ALUResult {
	carry := NewBit(carryIn)
	aBits := Unpack(a)
	bBits := Unpack(b)
	var dBits Bits
	var result ALUResult
	for i := range 8 {
		ai := aBits[i]
		bi := bBits[i]

		// full adder circuit
		di := ai ^ bi ^ carry
		nextCarry := (ai & bi) | ((ai ^ bi) & carry)

		dBits[i] = di
		if i == 3 {
			result.H = nextCarry.Bool()
		} else if i == 7 {
			result.C = nextCarry.Bool()
		}
		carry = nextCarry
	}
	result.Value = dBits.Pack()
	if expected := a + b + uint8(NewBit(carryIn)); result.Value != expected {
		panic(fmt.Sprintf("invalid calc: a=%v b=%v carry=%v, a+b+carry=%v, got %v", a, b, carryIn, expected, result.Value))
	}
	return result
}

func SUB(a, b uint8, carry bool) ALUResult {
	result := ADD(a, ^b, !carry)
	result.N = true
	result.C = !result.C
	result.H = !result.H
	if result.Value != uint8(int(a)-int(b)-int(NewBit(carry))) {
		panic("invalid calc")
	}
	return result
}

func OR(a, b uint8) ALUResult {
	return ALUResult{Value: a | b}
}

func AND(a, b uint8) ALUResult {
	return ALUResult{Value: a & b, H: true}
}

func XOR(a, b uint8) ALUResult {
	return ALUResult{Value: a ^ b}
}

func RL(a uint8, carry bool) ALUResult {
	var mask uint8
	if carry {
		mask = 0x01
	}
	return ALUResult{Value: a<<1 | mask, C: a&0x80 != 0}
}

func SRL(a uint8) ALUResult {
	return ALUResult{Value: a >> 1, C: a&0x01 != 0}
}

func RLA(a uint8, carry bool) ALUResult {
	res := RL(a, carry)
	res.Z0 = true
	return res
}

func RR(a uint8, carry bool) ALUResult {
	var mask uint8
	if carry {
		mask = 0x80
	}
	return ALUResult{Value: a>>1 | mask, C: a&1 != 0}
}

func RRA(a uint8, carry bool) ALUResult {
	res := RR(a, carry)
	res.Z0 = true
	return res
}

func RLCA(a uint8) ALUResult {
	var mask uint8
	if a&80 != 0 {
		mask = 1
	}
	return ALUResult{Value: a<<1 | mask, C: a&0x80 != 0, Z0: true}
}

func RRCA(a uint8) ALUResult {
	var mask uint8
	if a&1 != 0 {
		mask = 0x80
	}
	return ALUResult{Value: a>>1 | mask, C: a&1 != 0, Z0: true}
}

func SWAP(a uint8) ALUResult {
	return ALUResult{Value: ((a & 0x0f) << 4) | ((a & 0xf0) >> 4)}
}
