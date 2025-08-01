package model

import "math/bits"

//go:generate go-enum --marshal --flag --values --nocomments

// ENUM(FetchTileNo, FetchTileLSB, FetchTileMSB, PushFIFO)
type FetcherState Data8

type Fetcher struct {
	Cycle       uint64
	State       FetcherState
	X           Data8
	TileIndex   Data8
	TileLSBAddr Addr
	TileLSB     Data8
	TileMSB     Data8
	Suspended   bool
}

type BackgroundFetcher struct {
	Fetcher
	TileIndexAddr Addr

	TileOffsetX Addr
	TileOffsetY Addr

	WindowYReached                  bool
	WindowFetching                  bool
	WindowLineCounter               Data8
	WindowPixelRenderedThisScanline bool
}

func (gb *Gameboy) backgroundFetcherFSM() {
	bgf := &gb.PPU.BackgroundFetcher

	if bgf.Suspended {
		return
	}

	bgfCycle := bgf.Cycle
	bgf.Cycle++

	switch bgf.State {
	case FetcherStateFetchTileNo:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileNo(gb)
		bgf.State = FetcherStateFetchTileLSB
	case FetcherStateFetchTileLSB:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileLSB(gb)
		bgf.State = FetcherStateFetchTileMSB
	case FetcherStateFetchTileMSB:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileMSB(gb)
		bgf.State = FetcherStatePushFIFO
	case FetcherStatePushFIFO:
		if bgf.pushFIFO(gb) {
			bgf.State = FetcherStateFetchTileNo
			bgf.X++
		}
	}
}

func (bgf *BackgroundFetcher) fetchTileNo(gb *Gameboy) {
	// GBEDG: During the first step the fetcher fetches and stores the tile number of the tile which should be used.
	// Which Tilemap is used depends on whether the PPU is currently rendering Background or Window pixels
	// and on the bits 3 and 5 of the LCDC register.
	var addr Addr
	if bgf.WindowFetching {
		addr = gb.PPU.WindowTilemapArea()
		bgf.WindowPixelRenderedThisScanline = true
	} else {
		addr = gb.PPU.BGTilemapArea()
	}

	// GBEDG: Additionally, the address which the tile number is read from is offset by the fetcher-internal X-Position-Counter,
	//        which is incremented each time the last step is completed.
	offsetX := Addr(bgf.X)
	if !bgf.WindowFetching {
		// GBEDG: The value of SCX / 8 is also added if the Fetcher is not fetching Window pixels.
		//        In order to make the wrap-around with SCX work, this offset is ANDed with 0x1f
		offsetX += Addr((gb.PPU.RegSCX / 8))
		offsetX &= 0x1f
	}

	var offsetY Addr
	if !bgf.WindowFetching {
		// GBEDG: An offset of 32 * (((LY + SCY) & 0xFF) / 8) is also added if background pixels are being fetched,
		offs := (gb.PPU.RegLY + gb.PPU.RegSCY) & 0xff
		offs /= 8
		offsetY = 32 * Addr(offs)
	} else {
		// GBEDG: otherwise, if window pixels are being fetched, this offset is determined by 32 * (WINDOW_LINE_COUNTER / 8)
		offsetY = 32 * Addr(bgf.WindowLineCounter/8)
	}

	// GBEDG: Note: The sum of [...] the X-POS+SCX [...] is ANDed with 0x3ff in order to ensure that the address stays within the Tilemap memory regions.
	bgf.TileOffsetX = offsetX
	bgf.TileOffsetY = offsetY / 32
	addr += (offsetX + offsetY) & 0x3ff

	bgf.TileIndexAddr = addr
	bgf.TileIndex = gb.Mem[addr]
}

func (bgf *BackgroundFetcher) fetchTileLSB(gb *Gameboy) {
	idx := bgf.TileIndex
	signedAddressing := gb.PPU.RegLCDC&Bit4 == 0
	var addr Addr
	if signedAddressing {
		if idx < 128 {
			addr = Addr(0x9000 + 16*Addr(idx))
		} else {
			addr = Addr(0x8800 + 16*Addr(idx-128))
		}
	} else {
		addr = Addr(0x8000 + 16*Addr(idx))
	}
	if bgf.WindowFetching {
		addr += 2 * Addr(bgf.WindowLineCounter%8)
	} else {
		addr += 2 * Addr((gb.PPU.RegLY+gb.PPU.RegSCY)%8)
	}
	bgf.TileLSBAddr = addr
	bgf.TileLSB = gb.Mem[addr]
}

func (bgf *BackgroundFetcher) fetchTileMSB(gb *Gameboy) {
	bgf.TileMSB = gb.Mem[bgf.TileLSBAddr+1]
}

func (bgf *BackgroundFetcher) windowReached(gb *Gameboy) bool {
	if !gb.PPU.WindowEnable() {
		return false
	}
	if gb.PPU.RegWY == gb.PPU.RegLY {
		bgf.WindowYReached = true
	}
	if !bgf.WindowYReached {
		return false
	}
	if gb.PPU.Shifter.X < gb.PPU.RegWX-7 {
		return false
	}
	return true
}

func (bgf *BackgroundFetcher) pushFIFO(gb *Gameboy) bool {
	if gb.PPU.BackgroundFIFO.Level > 0 {
		return false
	}
	if gb.PPU.BGWindowEnable() {
		line := DecodeLineBG(bgf.TileMSB, bgf.TileLSB, gb.PPU.BGPalette)
		gb.PPU.BackgroundFIFO.Slots = line
	} else {
		gb.PPU.BackgroundFIFO.Slots = 0
	}
	gb.PPU.BackgroundFIFO.Level = 8
	return true
}

type SpriteFetcher struct {
	Fetcher
	SpriteIDX int
	DoneX     Data8
}

func (gb *Gameboy) spriteFetcherFSM() {
	sf := &gb.PPU.SpriteFetcher

	sfCycle := sf.Cycle
	sf.Cycle++

	if sf.Suspended {
		return
	}
	switch sf.State {
	case FetcherStateFetchTileNo:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileNo(gb)
		sf.State = FetcherStateFetchTileLSB
	case FetcherStateFetchTileLSB:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileLSB(gb)
		sf.State = FetcherStateFetchTileMSB
	case FetcherStateFetchTileMSB:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileMSB(gb)
		sf.State = FetcherStatePushFIFO
	case FetcherStatePushFIFO:
		sf.pushFIFO(gb)
		sf.State = FetcherStateFetchTileNo
		sf.Suspended = true
		sf.DoneX = gb.PPU.Shifter.X
	}
}

func (sf *SpriteFetcher) fetchTileNo(gb *Gameboy) {
	sf.TileIndex = gb.PPU.OAMBuffer.Buffer[sf.SpriteIDX].TileIndex
}

func (sf *SpriteFetcher) fetchTileLSB(gb *Gameboy) {
	obj := gb.PPU.OAMBuffer.Buffer[sf.SpriteIDX]

	screenY := gb.PPU.RegLY + gb.PPU.RegSCY
	offsetInObj := screenY - obj.Y

	addr := Addr(0x8000)
	if gb.PPU.ObjHeight() == 8 {
		addr += 16 * Addr(sf.TileIndex)
		if obj.Attributes&Bit6 != 0 {
			addr -= 2 * Addr(offsetInObj%8)
			addr += 2 * 7 // WHY DOES THAT WORK
		} else {
			addr += 2 * Addr(offsetInObj%8)
		}
	} else {
		if offsetInObj >= 8 {
			addr += 16 * Addr(sf.TileIndex&0xfe)
		} else {
			addr += 16 * Addr(sf.TileIndex|1)
		}
		if obj.Attributes&Bit6 != 0 {
			addr -= 2 * Addr(offsetInObj%16)
			addr += 2 * 14 // WHY DOES THAT WORK
		} else {
			addr += 2 * Addr(offsetInObj%16)
		}
	}

	sf.TileLSBAddr = addr
	sf.TileLSB = gb.Mem[addr]
}

func (sf *SpriteFetcher) fetchTileMSB(gb *Gameboy) {
	sf.TileMSB = gb.Mem[sf.TileLSBAddr+1]
}

func (sf *SpriteFetcher) pushFIFO(gb *Gameboy) {
	obj := gb.PPU.OAMBuffer.Buffer[sf.SpriteIDX]
	palette := gb.PPU.ObjPalette(obj.Attributes)

	line := DecodeLineSprite(sf.TileMSB, sf.TileLSB, palette)
	if obj.Attributes&Bit5 != 0 {
		line = bits.ReverseBytes64(line)
	}
	if !gb.PPU.OBJEnable() {
		line = 0
	}
	if obj.Attributes&Bit7 != 0 {
		line |= PxMaskPriority8
	}

	offset := gb.PPU.Shifter.X + 8 - obj.X
	pixelsToPush := int(8 - offset)
	level := gb.PPU.SpriteFIFO.Level

	existing := gb.PPU.SpriteFIFO.Slots & PXMaskColor8
	newPixels := line >> (offset * 8)

	// Overwrite transparent pixels
	N := min(pixelsToPush, level) * 8
	var mask, replacement uint64
	for i := 0; i < N; i += 8 {
		if (existing>>(i))&0xFF == 0 {
			p := (newPixels >> i) & 0xFF
			mask |= 0xFF << i
			replacement |= p << i
		}
	}
	gb.PPU.SpriteFIFO.Slots = (gb.PPU.SpriteFIFO.Slots & ^mask) | replacement

	// Append new pixels if level < pixelsToPush
	if level < pixelsToPush {
		count := pixelsToPush - level
		appendMask := (uint64(1) << (count * 8)) - 1
		appendData := (line >> (offset * 8)) & appendMask
		gb.PPU.SpriteFIFO.Slots |= appendData << (level * 8)
		gb.PPU.SpriteFIFO.Level = pixelsToPush
	}

}
