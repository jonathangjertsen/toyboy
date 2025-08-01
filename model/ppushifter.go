package model

type Shifter struct {
	Discard     Data8
	Suspended   bool
	X           Data8
	LastShifted Color
}

func (ps *Shifter) fsm(gb *Gameboy, clk *ClockRT) {
	if ps.Suspended {
		return
	}

	if ps.Discard > 0 {
		if _, shifted := gb.PPU.BackgroundFIFO.ShiftOut(); shifted {
			ps.Discard--
		}
	}

	pixel, havePixel := ps.pixelMixer(gb)
	if !havePixel {
		return
	}

	// Write pixel to LCD
	gb.PPU.FBViewport[gb.PPU.RegLY][ps.X] = pixel.Color()
	ps.LastShifted = pixel.Color()
	ps.X++

	gb.Debug.SetX(ps.X, clk)
}

func (ps *Shifter) pixelMixer(gb *Gameboy) (Pixel, bool) {
	spritePixel, haveSpritePixel := gb.PPU.SpriteFIFO.ShiftOut()
	bgPixel, haveBGPixel := gb.PPU.BackgroundFIFO.ShiftOut()
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
