package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type C = layout.Context
type D = layout.Dimensions
type F = layout.Flex
type FC = layout.FlexChild

var R = layout.Rigid

func Spacer(h, w int) FC {
	return R(layout.Spacer{Height: unit.Dp(h), Width: unit.Dp(w)}.Layout)
}

func Column(gtx layout.Context, children ...FC) D {
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx, children...)
}

func Row(gtx layout.Context, children ...FC) D {
	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceAround,
	}.Layout(gtx, children...)
}
