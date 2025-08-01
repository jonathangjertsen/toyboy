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
	PPUModeHBlank PPUMode = iota
	PPUModeVBlank
	PPUModeOAMScan
	PPUModePixelDraw
)

var ErrInvalidPPUMode = errors.New("not a valid PPUMode")

const _PPUModeName = "HBlankVBlankOAMScanPixelDraw"

// PPUModeValues returns a list of the values for PPUMode
func PPUModeValues() []PPUMode {
	return []PPUMode{
		PPUModeHBlank,
		PPUModeVBlank,
		PPUModeOAMScan,
		PPUModePixelDraw,
	}
}

var _PPUModeMap = map[PPUMode]string{
	PPUModeHBlank:    _PPUModeName[0:6],
	PPUModeVBlank:    _PPUModeName[6:12],
	PPUModeOAMScan:   _PPUModeName[12:19],
	PPUModePixelDraw: _PPUModeName[19:28],
}

// String implements the Stringer interface.
func (x PPUMode) String() string {
	if str, ok := _PPUModeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("PPUMode(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x PPUMode) IsValid() bool {
	_, ok := _PPUModeMap[x]
	return ok
}

var _PPUModeValue = map[string]PPUMode{
	_PPUModeName[0:6]:   PPUModeHBlank,
	_PPUModeName[6:12]:  PPUModeVBlank,
	_PPUModeName[12:19]: PPUModeOAMScan,
	_PPUModeName[19:28]: PPUModePixelDraw,
}

// ParsePPUMode attempts to convert a string to a PPUMode.
func ParsePPUMode(name string) (PPUMode, error) {
	if x, ok := _PPUModeValue[name]; ok {
		return x, nil
	}
	return PPUMode(0), fmt.Errorf("%s is %w", name, ErrInvalidPPUMode)
}

// MarshalText implements the text marshaller method.
func (x PPUMode) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *PPUMode) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParsePPUMode(name)
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
func (x *PPUMode) AppendText(b []byte) ([]byte, error) {
	return append(b, x.String()...), nil
}

// Set implements the Golang flag.Value interface func.
func (x *PPUMode) Set(val string) error {
	v, err := ParsePPUMode(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *PPUMode) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *PPUMode) Type() string {
	return "PPUMode"
}
