package model

import "image/color"

//go:generate go-enum --marshal --flag --values --nocomments

// ENUM(WhiteOrTransparent, LightGray, DarkGray, Black)
type Color uint8

func (c Color) Grayscale() uint8 {
	switch c {
	case ColorWhiteOrTransparent:
		return 0xf0
	case ColorLightGray:
		return 0xa0
	case ColorDarkGray:
		return 0x70
	case ColorBlack:
		return 0x30
	}
	return 0x00
}

func (c Color) RGBA() color.RGBA {
	switch c {
	case ColorWhiteOrTransparent:
		return color.RGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff}
	case ColorLightGray:
		return color.RGBA{R: 0xa0, G: 0xa0, B: 0xa0, A: 0xff}
	case ColorDarkGray:
		return color.RGBA{R: 0x70, G: 0x70, B: 0x70, A: 0xff}
	case ColorBlack:
		return color.RGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xff}
	}
	return color.RGBA{}
}
