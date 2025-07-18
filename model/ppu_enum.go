// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

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

const (
	PixelFetcherStateFetchTileNo PixelFetcherState = iota
	PixelFetcherStateFetchTileLSB
	PixelFetcherStateFetchTileMSB
	PixelFetcherStatePushFIFO
)

var ErrInvalidPixelFetcherState = errors.New("not a valid PixelFetcherState")

const _PixelFetcherStateName = "FetchTileNoFetchTileLSBFetchTileMSBPushFIFO"

// PixelFetcherStateValues returns a list of the values for PixelFetcherState
func PixelFetcherStateValues() []PixelFetcherState {
	return []PixelFetcherState{
		PixelFetcherStateFetchTileNo,
		PixelFetcherStateFetchTileLSB,
		PixelFetcherStateFetchTileMSB,
		PixelFetcherStatePushFIFO,
	}
}

var _PixelFetcherStateMap = map[PixelFetcherState]string{
	PixelFetcherStateFetchTileNo:  _PixelFetcherStateName[0:11],
	PixelFetcherStateFetchTileLSB: _PixelFetcherStateName[11:23],
	PixelFetcherStateFetchTileMSB: _PixelFetcherStateName[23:35],
	PixelFetcherStatePushFIFO:     _PixelFetcherStateName[35:43],
}

// String implements the Stringer interface.
func (x PixelFetcherState) String() string {
	if str, ok := _PixelFetcherStateMap[x]; ok {
		return str
	}
	return fmt.Sprintf("PixelFetcherState(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x PixelFetcherState) IsValid() bool {
	_, ok := _PixelFetcherStateMap[x]
	return ok
}

var _PixelFetcherStateValue = map[string]PixelFetcherState{
	_PixelFetcherStateName[0:11]:  PixelFetcherStateFetchTileNo,
	_PixelFetcherStateName[11:23]: PixelFetcherStateFetchTileLSB,
	_PixelFetcherStateName[23:35]: PixelFetcherStateFetchTileMSB,
	_PixelFetcherStateName[35:43]: PixelFetcherStatePushFIFO,
}

// ParsePixelFetcherState attempts to convert a string to a PixelFetcherState.
func ParsePixelFetcherState(name string) (PixelFetcherState, error) {
	if x, ok := _PixelFetcherStateValue[name]; ok {
		return x, nil
	}
	return PixelFetcherState(0), fmt.Errorf("%s is %w", name, ErrInvalidPixelFetcherState)
}

// MarshalText implements the text marshaller method.
func (x PixelFetcherState) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *PixelFetcherState) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParsePixelFetcherState(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// Set implements the Golang flag.Value interface func.
func (x *PixelFetcherState) Set(val string) error {
	v, err := ParsePixelFetcherState(val)
	*x = v
	return err
}

// Get implements the Golang flag.Getter interface func.
func (x *PixelFetcherState) Get() interface{} {
	return *x
}

// Type implements the github.com/spf13/pFlag Value interface.
func (x *PixelFetcherState) Type() string {
	return "PixelFetcherState"
}
