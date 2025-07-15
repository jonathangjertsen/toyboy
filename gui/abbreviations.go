package gui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type C = layout.Context
type D = layout.Dimensions
type F = layout.Flex

var Rigid = layout.Rigid

func Spacer(h, w int) layout.FlexChild {
	return Rigid(layout.Spacer{Height: unit.Dp(h), Width: unit.Dp(w)}.Layout)
}

func Column(gtx layout.Context, children ...layout.FlexChild) D {
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx, children...)
}

func Row(gtx layout.Context, children ...layout.FlexChild) D {
	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceAround,
	}.Layout(gtx, children...)
}
