package model

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

func (bi Bits) Pack() Data8 {
	var out Data8
	for i := range 8 {
		out |= Data8(bi[i]) << i
	}
	return out
}

type ALUResult struct {
	Value Data8
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

func Unpack(b8 Data8) Bits {
	var out Bits
	for i := range 8 {
		out[i] = Bit((b8 & (1 << i)) >> i)
	}
	return out
}

func ADD(a, b Data8, carry bool) ALUResult {
	cint := 0
	if carry {
		cint = 1
	}
	v := int(a) + int(b) + cint
	hv := int(a&0xf) + int(b&0xf) + cint
	return ALUResult{Value: Data8(v), C: v > 0xff, H: hv > 0xf}
}

func SUB(a, b Data8, carry bool) ALUResult {
	cint := 0
	if carry {
		cint = 1
	}
	v := int(a) - int(b) - cint
	hv := int(a&0xf) - int(b&0xf) - cint
	return ALUResult{Value: Data8(v), C: v < 0, H: hv < 0, N: true}
}

func OR(a, b Data8) ALUResult {
	return ALUResult{Value: a | b}
}

func AND(a, b Data8) ALUResult {
	return ALUResult{Value: a & b, H: true}
}

func XOR(a, b Data8) ALUResult {
	return ALUResult{Value: a ^ b}
}

func RL(a Data8, carry bool) ALUResult {
	var mask Data8
	if carry {
		mask = 0x01
	}
	return ALUResult{Value: a<<1 | mask, C: a.Bit(7)}
}

func SLA(a Data8) ALUResult {
	return ALUResult{Value: a << 1, C: a.Bit(7)}
}

func SRA(a Data8) ALUResult {
	var mask Data8
	if a.Bit(7) {
		mask = 0x80
	}
	return ALUResult{Value: a>>1 | mask, C: a.Bit(0)}
}

func RLC(a Data8) ALUResult {
	b7 := a.Bit(7)
	var mask Data8
	if b7 {
		mask = 0x01
	}
	return ALUResult{Value: a<<1 | mask, C: b7}
}

func SRL(a Data8) ALUResult {
	return ALUResult{Value: a >> 1, C: a.Bit(0)}
}

func RLA(a Data8, carry bool) ALUResult {
	res := RL(a, carry)
	res.Z0 = true
	return res
}

func RR(a Data8, carry bool) ALUResult {
	var mask Data8
	if carry {
		mask = 0x80
	}
	return ALUResult{Value: a>>1 | mask, C: a&1 != 0}
}

func RRC(a Data8) ALUResult {
	b0 := a.Bit(0)
	var mask Data8
	if b0 {
		mask = 0x80
	}
	return ALUResult{Value: a>>1 | mask, C: b0}
}

func RRA(a Data8, carry bool) ALUResult {
	res := RR(a, carry)
	res.Z0 = true
	return res
}

func RLCA(a Data8) ALUResult {
	var mask Data8
	if a&80 != 0 {
		mask = 1
	}
	return ALUResult{Value: a<<1 | mask, C: a&0x80 != 0, Z0: true}
}

func RRCA(a Data8) ALUResult {
	var mask Data8
	if a&1 != 0 {
		mask = 0x80
	}
	return ALUResult{Value: a>>1 | mask, C: a&1 != 0, Z0: true}
}

func SWAP(a Data8) ALUResult {
	return ALUResult{Value: ((a & 0x0f) << 4) | ((a & 0xf0) >> 4)}
}

func DAA(a Data8, c, n, h bool) ALUResult {
	va := int(a)
	if !n {
		if h || ((va & 0xf) > 9) {
			va += 0x06
		}
		if c || (va > 0x9f) {
			va += 0x60
		}
	} else {
		if h {
			va = (va - 6) & 0xff
		}
		if c {
			va -= 0x60
		}
	}
	if va&0x100 == 0x100 {
		c = true
	}
	return ALUResult{Value: Data8(va & 0xff), C: c}

}
