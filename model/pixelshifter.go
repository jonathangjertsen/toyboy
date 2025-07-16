package model

type PixelShifter struct {
	Discard     uint8
	Suspended   bool
	X           uint8
	LastShifted Color

	PPU *PPU
}

func (ps *PixelShifter) fsm() {
	if ps.Suspended {
		return
	}

	if ps.Discard > 0 {
		if _, shifted := ps.PPU.BackgroundFIFO.ShiftOut(); shifted {
			ps.Discard--
		}
		return
	}

	pixel, havePixel := ps.getPixel()
	if !havePixel {
		return
	}

	// Write pixel to LCD
	ps.PPU.FBViewport[ps.PPU.RegLY][ps.X] = pixel.Color
	ps.LastShifted = pixel.Color
	ps.X++
}

func (ps *PixelShifter) getPixel() (Pixel, bool) {
	spritePixel, haveSpritePixel := ps.PPU.SpriteFIFO.ShiftOut()
	bgPixel, haveBGPixel := ps.PPU.BackgroundFIFO.ShiftOut()
	if haveSpritePixel && haveBGPixel {
		if spritePixel.Color == ColorWhiteOrTransparent {
			return bgPixel, true
		} else if spritePixel.BackgroundPriority && bgPixel.Color != ColorWhiteOrTransparent {
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
