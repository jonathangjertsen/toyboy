package model

import "image/color"

//go:generate go-enum --marshal --flag --values --nocomments

// ENUM(WhiteOrTransparent, LightGray, DarkGray, Black)
type Color uint8

var Grayscale = [4]uint8{0xf0, 0xa0, 0x70, 0x30}
var RGBA = [4]color.RGBA{
	{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff},
	{R: 0xa0, G: 0xa0, B: 0xa0, A: 0xff},
	{R: 0x70, G: 0x70, B: 0x70, A: 0xff},
	{R: 0x30, G: 0x30, B: 0x30, A: 0xff},
}
