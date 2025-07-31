package model

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

func (bgf *BackgroundFetcher) fsm(ppu *PPU, mem []Data8) {
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
		bgf.fetchTileNo(ppu, mem)
		bgf.State = FetcherStateFetchTileLSB
	case FetcherStateFetchTileLSB:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileLSB(ppu, mem)
		bgf.State = FetcherStateFetchTileMSB
	case FetcherStateFetchTileMSB:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileMSB(ppu, mem)
		bgf.State = FetcherStatePushFIFO
	case FetcherStatePushFIFO:
		if bgf.pushFIFO(ppu) {
			bgf.State = FetcherStateFetchTileNo
			bgf.X++
		}
	}
}

func (bgf *BackgroundFetcher) fetchTileNo(ppu *PPU, mem []Data8) {
	// GBEDG: During the first step the fetcher fetches and stores the tile number of the tile which should be used.
	// Which Tilemap is used depends on whether the PPU is currently rendering Background or Window pixels
	// and on the bits 3 and 5 of the LCDC register.
	var addr Addr
	if bgf.WindowFetching {
		addr = ppu.WindowTilemapArea()
		bgf.WindowPixelRenderedThisScanline = true
	} else {
		addr = ppu.BGTilemapArea()
	}

	// GBEDG: Additionally, the address which the tile number is read from is offset by the fetcher-internal X-Position-Counter,
	//        which is incremented each time the last step is completed.
	offsetX := Addr(bgf.X)
	if !bgf.WindowFetching {
		// GBEDG: The value of SCX / 8 is also added if the Fetcher is not fetching Window pixels.
		//        In order to make the wrap-around with SCX work, this offset is ANDed with 0x1f
		offsetX += Addr((ppu.RegSCX / 8))
		offsetX &= 0x1f
	}

	var offsetY Addr
	if !bgf.WindowFetching {
		// GBEDG: An offset of 32 * (((LY + SCY) & 0xFF) / 8) is also added if background pixels are being fetched,
		offs := (ppu.RegLY + ppu.RegSCY) & 0xff
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
	bgf.TileIndex = mem[addr]
}

func (bgf *BackgroundFetcher) fetchTileLSB(ppu *PPU, mem []Data8) {
	idx := bgf.TileIndex
	signedAddressing := ppu.RegLCDC&Bit4 == 0
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
		addr += 2 * Addr((ppu.RegLY+ppu.RegSCY)%8)
	}
	bgf.TileLSBAddr = addr
	bgf.TileLSB = mem[addr]
}

func (bgf *BackgroundFetcher) fetchTileMSB(ppu *PPU, mem []Data8) {
	bgf.TileMSB = mem[bgf.TileLSBAddr+1]
}

func (bgf *BackgroundFetcher) windowReached(ppu *PPU) bool {
	if !ppu.WindowEnable() {
		return false
	}
	if ppu.RegWY == ppu.RegLY {
		bgf.WindowYReached = true
	}
	if !bgf.WindowYReached {
		return false
	}
	if ppu.Shifter.X < ppu.RegWX-7 {
		return false
	}
	return true
}

func (bgf *BackgroundFetcher) pushFIFO(ppu *PPU) bool {
	if ppu.BackgroundFIFO.Level > 0 {
		return false
	}
	if ppu.BGWindowEnable() {
		line := DecodeTileLine(bgf.TileMSB, bgf.TileLSB)
		for i := range 8 {
			line[i].Palette = ppu.BGPalette
		}
		ppu.BackgroundFIFO.Write8(line)
	} else {
		ppu.BackgroundFIFO.Write8(TileLine{})
	}
	return true
}

type SpriteFetcher struct {
	Fetcher
	SpriteIDX int
	DoneX     Data8
}

func (sf *SpriteFetcher) fsm(ppu *PPU, mem []Data8) {
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
		sf.fetchTileNo(ppu)
		sf.State = FetcherStateFetchTileLSB
	case FetcherStateFetchTileLSB:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileLSB(ppu, mem)
		sf.State = FetcherStateFetchTileMSB
	case FetcherStateFetchTileMSB:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileMSB(ppu, mem)
		sf.State = FetcherStatePushFIFO
	case FetcherStatePushFIFO:
		if sf.pushFIFO(ppu) {
			sf.State = FetcherStateFetchTileNo
			sf.Suspended = true
			sf.DoneX = ppu.Shifter.X
		}
	}
}

func (sf *SpriteFetcher) fetchTileNo(ppu *PPU) {
	sf.TileIndex = ppu.OAMBuffer.Buffer[sf.SpriteIDX].TileIndex
}

func (sf *SpriteFetcher) fetchTileLSB(ppu *PPU, mem []Data8) {
	obj := ppu.OAMBuffer.Buffer[sf.SpriteIDX]

	screenY := ppu.RegLY + ppu.RegSCY
	offsetInObj := screenY - obj.Y

	addr := Addr(0x8000)
	if ppu.ObjHeight() == 8 {
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
	sf.TileLSB = mem[addr]
}

func (sf *SpriteFetcher) fetchTileMSB(ppu *PPU, mem []Data8) {
	sf.TileMSB = mem[sf.TileLSBAddr+1]
}

func (sf *SpriteFetcher) pushFIFO(ppu *PPU) bool {
	obj := ppu.OAMBuffer.Buffer[sf.SpriteIDX]
	palette := ppu.ObjPalette(obj.Attributes)
	line := DecodeTileLine(sf.TileMSB, sf.TileLSB)
	if obj.Attributes&Bit5 != 0 {
		for i := range 4 {
			line[i], line[7-i] = line[7-i], line[i]
		}
	}
	for i := range 8 {
		line[i].Palette = palette
	}
	if !ppu.OBJEnable() {
		for i := range 8 {
			line[i].ColorIDXBGPriority = 0
		}
	}
	if obj.Attributes&Bit7 != 0 {
		for i := range 8 {
			line[i].SetBGPriority()
		}
	}
	offsetInSprite := ppu.Shifter.X + 8 - obj.X
	pixelsToPush := 8 - offsetInSprite
	pos := ppu.SpriteFIFO.ShiftPos
	for i := range int(pixelsToPush) {
		incLevel := i >= ppu.SpriteFIFO.Level
		pushPixel := incLevel || ppu.SpriteFIFO.Slots[pos].ColorIDX() == 0
		if pushPixel {
			pixel := line[int(offsetInSprite)+i]
			ppu.SpriteFIFO.Slots[pos] = pixel
		}
		if incLevel {
			ppu.SpriteFIFO.Level++
		}
		pos++
		if pos == 8 {
			pos = 0
		}
	}
	return true
}
