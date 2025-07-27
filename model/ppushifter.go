package model

type Shifter struct {
	Discard     Data8
	Suspended   bool
	X           Data8
	LastShifted Color

	PPU *PPU
}

func (ps *Shifter) fsm() {
	if ps.Suspended {
		return
	}

	if ps.Discard > 0 {
		if _, shifted := ps.PPU.BackgroundFIFO.ShiftOut(); shifted {
			ps.Discard--
		}
	}

	pixel, havePixel := ps.pixelMixer()
	if !havePixel {
		return
	}

	// Write pixel to LCD
	ps.PPU.FBViewport[ps.PPU.RegLY][ps.X] = pixel.Color()
	ps.LastShifted = pixel.Color()
	ps.X++

	ps.PPU.Debug.SetX(ps.X)
}

func (ps *Shifter) pixelMixer() (Pixel, bool) {
	spritePixel, haveSpritePixel := ps.PPU.SpriteFIFO.ShiftOut()
	bgPixel, haveBGPixel := ps.PPU.BackgroundFIFO.ShiftOut()
	if haveSpritePixel && haveBGPixel {
		if spritePixel.ColorIDX() == 0 {
			return bgPixel, true
		} else if spritePixel.BGPriority() && bgPixel.ColorIDXBGPriority > 0 { // strictly speaking should be bgPixel.ColorIDX() > 0, but BGPriority is never set on background pixels
			return bgPixel, true
		} else {
			return spritePixel, true
		}
	} else if haveSpritePixel {
		return spritePixel, true
	} else if haveBGPixel {
		return bgPixel, true
	}
	return Pixel{}, false
}
