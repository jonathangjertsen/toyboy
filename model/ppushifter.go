package model

type Shifter struct {
	Discard     Data8
	Suspended   bool
	X           Data8
	LastShifted Color
}

func (gb *Gameboy) shifterFSM(clk *ClockRT) {
	ps := &gb.PPU.Shifter

	if ps.Suspended {
		return
	}

	gb.shiftDiscarded()
	pixel, havePixel := ps.pixelMixer(gb)
	if !havePixel {
		return
	}
	gb.writePixelToLCD(pixel)
	gb.Debug.SetX(ps.X, clk)
}

func (gb *Gameboy) shiftDiscarded() {
	ps := &gb.PPU.Shifter

	if ps.Discard > 0 {
		if _, shifted := gb.PPU.BackgroundFIFO.ShiftOut(); shifted {
			ps.Discard--
		}
	}
}

func (gb *Gameboy) writePixelToLCD(pixel Pixel) {
	ps := &gb.PPU.Shifter

	gb.PPU.FBViewport[gb.PPU.RegLY][ps.X] = pixel.Color()
	ps.LastShifted = pixel.Color()
	ps.X++
}

func (ps *Shifter) pixelMixer(gb *Gameboy) (Pixel, bool) {
	spritePixel, haveSpritePixel := gb.PPU.SpriteFIFO.ShiftOut()
	bgPixel, haveBGPixel := gb.PPU.BackgroundFIFO.ShiftOut()
	if haveSpritePixel && haveBGPixel {
		if (spritePixel & 0x3) == 0 {
			return bgPixel, true
		} else if (spritePixel&PxMaskPriority != 0) && (bgPixel&0x3 != 0) {
			return bgPixel, true
		} else {
			return spritePixel, true
		}
	} else if haveSpritePixel {
		return spritePixel, true
	} else if haveBGPixel {
		return bgPixel, true
	}
	return 0, false
}
