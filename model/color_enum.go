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
	ColorWhiteOrTransparent Color = iota
	ColorLightGray
	ColorDarkGray
	ColorBlack
)

var ErrInvalidColor = errors.New("not a valid Color")

const _ColorName = "WhiteOrTransparentLightGrayDarkGrayBlack"

// ColorValues returns a list of the values for Color
func ColorValues() []Color {
	return []Color{
		ColorWhiteOrTransparent,
		ColorLightGray,
		ColorDarkGray,
		ColorBlack,
	}
}

var _ColorMap = map[Color]string{
	ColorWhiteOrTransparent: _ColorName[0:18],
	ColorLightGray:          _ColorName[18:27],
	ColorDarkGray:           _ColorName[27:35],
	ColorBlack:              _ColorName[35:40],
}

// String implements the Stringer interface.
func (x Color) String() string {
	if str, ok := _ColorMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Color(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Color) IsValid() bool {
	_, ok := _ColorMap[x]
	return ok
}

var _ColorValue = map[string]Color{
	_ColorName[0:18]:  ColorWhiteOrTransparent,
	_ColorName[18:27]: ColorLightGray,
	_ColorName[27:35]: ColorDarkGray,
	_ColorName[35:40]: ColorBlack,
}

// ParseColor attempts to convert a string to a Color.
func ParseColor(name string) (Color, error) {
	if x, ok := _ColorValue[name]; ok {
		return x, nil
	}
	return Color(0), fmt.Errorf("%s is %w", name, ErrInvalidColor)
}

// MarshalText implements the text marshaller method.
func (x Color) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Color) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseColor(name)
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
func (x *Color) AppendText(b []byte) ([]byte, error) {
	return append(b, x.String()...), nil
}

// Set implements the Golang flag.Value interface func.
func (x *Color) Set(val string) error {
	v, err := ParseColor(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *Color) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *Color) Type() string {
	return "Color"
}
