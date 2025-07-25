// Code generated by go-enum DO NOT EDIT.
// Version: 0.9.0
// Revision: 4061a5d82779342c5863a515363feb943fa59455
// Build Date: 2025-07-22T03:42:20Z
// Built By: goreleaser

package model

import (
	"errors"
	"fmt"
)

const (
	IntSourceVBlank IntSource = iota + 1
	IntSourceLCD
	IntSourceTimer
	IntSourceSerial
	IntSourceJoypad
)

var ErrInvalidIntSource = errors.New("not a valid IntSource")

const _IntSourceName = "VBlankLCDTimerSerialJoypad"

// IntSourceValues returns a list of the values for IntSource
func IntSourceValues() []IntSource {
	return []IntSource{
		IntSourceVBlank,
		IntSourceLCD,
		IntSourceTimer,
		IntSourceSerial,
		IntSourceJoypad,
	}
}

var _IntSourceMap = map[IntSource]string{
	IntSourceVBlank: _IntSourceName[0:6],
	IntSourceLCD:    _IntSourceName[6:9],
	IntSourceTimer:  _IntSourceName[9:14],
	IntSourceSerial: _IntSourceName[14:20],
	IntSourceJoypad: _IntSourceName[20:26],
}

// String implements the Stringer interface.
func (x IntSource) String() string {
	if str, ok := _IntSourceMap[x]; ok {
		return str
	}
	return fmt.Sprintf("IntSource(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x IntSource) IsValid() bool {
	_, ok := _IntSourceMap[x]
	return ok
}

var _IntSourceValue = map[string]IntSource{
	_IntSourceName[0:6]:   IntSourceVBlank,
	_IntSourceName[6:9]:   IntSourceLCD,
	_IntSourceName[9:14]:  IntSourceTimer,
	_IntSourceName[14:20]: IntSourceSerial,
	_IntSourceName[20:26]: IntSourceJoypad,
}

// ParseIntSource attempts to convert a string to a IntSource.
func ParseIntSource(name string) (IntSource, error) {
	if x, ok := _IntSourceValue[name]; ok {
		return x, nil
	}
	return IntSource(0), fmt.Errorf("%s is %w", name, ErrInvalidIntSource)
}

// MarshalText implements the text marshaller method.
func (x IntSource) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *IntSource) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseIntSource(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// AppendText appends the textual representation of itself to the end of b
// (allocating a larger slice if necessary) and returns the updated slice.
//
// Implementations must not retain b, nor mutate any bytes within b[:len(b)].
func (x *IntSource) AppendText(b []byte) ([]byte, error) {
	return append(b, x.String()...), nil
}

// Set implements the Golang flag.Value interface func.
func (x *IntSource) Set(val string) error {
	v, err := ParseIntSource(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *IntSource) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *IntSource) Type() string {
	return "IntSource"
}
