package model

import "fmt"

const (
	Bit0 = 1 << iota
	Bit1
	Bit2
	Bit3
	Bit4
	Bit5
	Bit6
	Bit7
	Bit8
	Bit9
	Bit10
	Bit11
	Bit12
	Bit13
	Bit14
	Bit15

	SignBit8  = Bit7
	SignBit16 = Bit15
)

type Data16 uint16

func (a Data16) Bit(i int) bool {
	return a&(1<<i) != 0
}

func (a Data16) Bit4() bool {
	return a&Bit4 != 0
}

func (a Data16) MSB() Data8 {
	return Data8(a >> 8)
}

func (a Data16) LSB() Data8 {
	return Data8(a)
}

func (a Data16) Hex() string {
	return Hex16(uint16(a))
}

func (a Data16) Dec() string {
	return fmt.Sprintf("%dd", a)
}

func (a Data16) Split() (Data8, Data8) {
	return a.MSB(), a.LSB()
}

type Data8 uint8

func (a Data8) Bit(i uint) bool {
	return a&(1<<i) != 0
}

func (a Data8) Bit0() bool {
	return a&Bit0 != 0
}

func (a Data8) Bit1() bool {
	return a&Bit1 != 0
}

func (a Data8) Bit2() bool {
	return a&Bit2 != 0
}

func (a Data8) Bit3() bool {
	return a&Bit3 != 0
}

func (a Data8) Bit4() bool {
	return a&Bit4 != 0
}

func (a Data8) Bit5() bool {
	return a&Bit5 != 0
}

func (a Data8) Bit6() bool {
	return a&Bit6 != 0
}

func (a Data8) Bit7() bool {
	return a&Bit7 != 0
}

func (a Data8) SignBit() bool {
	return a&Bit7 != 0
}

func (a Data8) SignedOffset() Offset8 {
	return Offset8(int8(uint8(a)))
}
func (a Data8) SignedAbs() Data8 {
	if a&SignBit8 != 0 {
		return Data8(-int8(a))
	}
	return a
}

func (a Data8) Hex() string {
	return Hex8(uint8(a))
}

func (a Data8) Dec() string {
	return fmt.Sprintf("%d", a)
}

func Hex16(x uint16) string {
	return fmt.Sprintf("%04x", x)
}

func Hex8(x uint8) string {
	return fmt.Sprintf("%02x", x)
}
