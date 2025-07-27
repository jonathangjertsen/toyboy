package main

import (
	"github.com/jonathangjertsen/toyboy/model"
)

//go:generate go-enum --marshal --flag --values --nocomments

type ButtonMapping struct {
	Up     string
	Left   string
	Right  string
	Down   string
	A      string
	B      string
	Start  string
	Select string

	SOCDResolution SOCDResolution
}

// Resultion to key combinations that are physically impossible on a D-Pad (left+right, up+down)
// List methods suggested here https://www.reddit.com/r/StreetFighter/comments/17gn4zm/comment/llwirmc
// ENUM(Unfiltered, OppositeNeutral, FirstWinsSecondDisabled, OppositeNeutralFirstDisabled)
type SOCDResolution int

func (bm *ButtonMapping) JoypadState(keystate map[string]bool) model.JoypadState {
	jps := model.JoypadState{}
	for key, pressed := range keystate {
		if pressed {
			switch key {
			case bm.Up:
				jps.Up = true
			case bm.Left:
				jps.Left = true
			case bm.Right:
				jps.Right = true
			case bm.Down:
				jps.Down = true
			case bm.A:
				jps.A = true
			case bm.B:
				jps.B = true
			case bm.Start:
				jps.Start = true
			case bm.Select:
				jps.Select = true
			}
		}
	}

	switch bm.SOCDResolution {
	case SOCDResolutionUnfiltered:
	case SOCDResolutionOppositeNeutral:
		if jps.Right && jps.Left {
			jps.Right = false
			jps.Left = false
		}
		if jps.Up && jps.Down {
			jps.Up = false
			jps.Down = false
		}
	case SOCDResolutionFirstWinsSecondDisabled:
		panic("not implemented")
	case SOCDResolutionOppositeNeutralFirstDisabled:
		panic("not implemented")
	}

	return jps
}
