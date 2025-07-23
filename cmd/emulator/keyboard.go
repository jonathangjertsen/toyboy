package main

import (
	"gioui.org/io/event"
	"gioui.org/io/key"
	"github.com/jonathangjertsen/toyboy/model"
)

//go:generate go-enum --marshal --flag --values --nocomments

type ButtonMapping struct {
	Up     key.Name
	Left   key.Name
	Right  key.Name
	Down   key.Name
	A      key.Name
	B      key.Name
	Start  key.Name
	Select key.Name

	SOCDResolution SOCDResolution

	cachedEventFilters []event.Filter
}

// Resultion to key combinations that are physically impossible on a D-Pad (left+right, up+down)
// List methods suggested here https://www.reddit.com/r/StreetFighter/comments/17gn4zm/comment/llwirmc
// ENUM(Unfiltered, OppositeNeutral, FirstWinsSecondDisabled, OppositeNeutralFirstDisabled)
type SOCDResolution int

func (bm *ButtonMapping) eventFilters() []event.Filter {
	if bm.cachedEventFilters == nil {
		bm.cachedEventFilters = []event.Filter{
			key.Filter{Name: bm.Up},
			key.Filter{Name: bm.Down},
			key.Filter{Name: bm.Left},
			key.Filter{Name: bm.Right},
			key.Filter{Name: bm.A},
			key.Filter{Name: bm.B},
			key.Filter{Name: bm.Start},
			key.Filter{Name: bm.Select},
		}
	}
	return bm.cachedEventFilters
}

func (bm *ButtonMapping) JoypadState(keystate map[key.Name]bool) model.JoypadState {
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

type KeyboardControl struct {
	Pressed       map[key.Name]bool
	ButtonMapping ButtonMapping
}

func NewKeyboardControl(buttonMapping ButtonMapping) *KeyboardControl {
	return &KeyboardControl{
		Pressed:       make(map[key.Name]bool),
		ButtonMapping: buttonMapping,
	}
}

func (ks *KeyboardControl) Frame(gtx C) model.JoypadState {
	for _, filter := range ks.ButtonMapping.eventFilters() {
		if keypress, ok := gtx.Event(filter); ok {
			if kev, ok := keypress.(key.Event); ok {
				pressed := kev.State == key.Press
				ks.Pressed[kev.Name] = pressed
			}
		}
	}
	return ks.ButtonMapping.JoypadState(ks.Pressed)
}
