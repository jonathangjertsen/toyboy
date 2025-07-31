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

	ppu *PPU
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

func (bgf *BackgroundFetcher) fsm() {
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
		bgf.fetchTileNo()
		bgf.State = FetcherStateFetchTileLSB
	case FetcherStateFetchTileLSB:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileLSB()
		bgf.State = FetcherStateFetchTileMSB
	case FetcherStateFetchTileMSB:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileMSB()
		bgf.State = FetcherStatePushFIFO
	case FetcherStatePushFIFO:
		if bgf.pushFIFO() {
			bgf.State = FetcherStateFetchTileNo
			bgf.X++
		}
	}
}

func (bgf *BackgroundFetcher) fetchTileNo() {
	// GBEDG: During the first step the fetcher fetches and stores the tile number of the tile which should be used.
	// Which Tilemap is used depends on whether the PPU is currently rendering Background or Window pixels
	// and on the bits 3 and 5 of the LCDC register.
	var addr Addr
	if bgf.WindowFetching {
		addr = bgf.ppu.WindowTilemapArea()
		bgf.WindowPixelRenderedThisScanline = true
	} else {
		addr = bgf.ppu.BGTilemapArea()
	}

	// GBEDG: Additionally, the address which the tile number is read from is offset by the fetcher-internal X-Position-Counter,
	//        which is incremented each time the last step is completed.
	offsetX := Addr(bgf.X)
	if !bgf.WindowFetching {
		// GBEDG: The value of SCX / 8 is also added if the Fetcher is not fetching Window pixels.
		//        In order to make the wrap-around with SCX work, this offset is ANDed with 0x1f
		offsetX += Addr((bgf.ppu.RegSCX / 8))
		offsetX &= 0x1f
	}

	var offsetY Addr
	if !bgf.WindowFetching {
		// GBEDG: An offset of 32 * (((LY + SCY) & 0xFF) / 8) is also added if background pixels are being fetched,
		offs := (bgf.ppu.RegLY + bgf.ppu.RegSCY) & 0xff
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
	bgf.TileIndex = bgf.ppu.Bus.ProbeAddress(addr)
}

func (bgf *BackgroundFetcher) fetchTileLSB() {
	idx := bgf.TileIndex
	signedAddressing := bgf.ppu.RegLCDC&Bit4 == 0
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
		addr += 2 * Addr((bgf.ppu.RegLY+bgf.ppu.RegSCY)%8)
	}
	bgf.TileLSBAddr = addr
	bgf.TileLSB = bgf.ppu.Bus.ProbeAddress(addr)
}

func (bgf *BackgroundFetcher) fetchTileMSB() {
	bgf.TileMSB = bgf.ppu.Bus.ProbeAddress(bgf.TileLSBAddr + 1)
}

func (bgf *BackgroundFetcher) windowReached() bool {
	if !bgf.ppu.WindowEnable() {
		return false
	}
	if bgf.ppu.RegWY == bgf.ppu.RegLY {
		bgf.WindowYReached = true
	}
	if !bgf.WindowYReached {
		return false
	}
	if bgf.ppu.Shifter.X < bgf.ppu.RegWX-7 {
		return false
	}
	return true
}

func (bgf *BackgroundFetcher) pushFIFO() bool {
	if bgf.ppu.BackgroundFIFO.Level > 0 {
		return false
	}
	if bgf.ppu.BGWindowEnable() {
		line := DecodeTileLine(bgf.TileMSB, bgf.TileLSB)
		for i := range 8 {
			line[i].Palette = bgf.ppu.BGPalette
		}
		bgf.ppu.BackgroundFIFO.Write8(line)
	} else {
		bgf.ppu.BackgroundFIFO.Write8(TileLine{})
	}
	return true
}

type SpriteFetcher struct {
	Fetcher
	SpriteIDX int
	DoneX     Data8
}

func (sf *SpriteFetcher) fsm() {
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
		sf.fetchTileNo()
		sf.State = FetcherStateFetchTileLSB
	case FetcherStateFetchTileLSB:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileLSB()
		sf.State = FetcherStateFetchTileMSB
	case FetcherStateFetchTileMSB:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileMSB()
		sf.State = FetcherStatePushFIFO
	case FetcherStatePushFIFO:
		if sf.pushFIFO() {
			sf.State = FetcherStateFetchTileNo
			sf.Suspended = true
			sf.DoneX = sf.ppu.Shifter.X
		}
	}
}

func (sf *SpriteFetcher) fetchTileNo() {
	sf.TileIndex = sf.ppu.OAMBuffer.Buffer[sf.SpriteIDX].TileIndex
}

func (sf *SpriteFetcher) fetchTileLSB() {
	obj := sf.ppu.OAMBuffer.Buffer[sf.SpriteIDX]

	screenY := sf.ppu.RegLY + sf.ppu.RegSCY
	offsetInObj := screenY - obj.Y

	addr := Addr(0x8000)
	if sf.ppu.ObjHeight() == 8 {
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
	sf.TileLSB = sf.ppu.Bus.ProbeAddress(addr)
}

func (sf *SpriteFetcher) fetchTileMSB() {
	sf.TileMSB = sf.ppu.Bus.ProbeAddress(sf.TileLSBAddr + 1)
}

func (sf *SpriteFetcher) pushFIFO() bool {
	obj := sf.ppu.OAMBuffer.Buffer[sf.SpriteIDX]
	palette := sf.ppu.ObjPalette(obj.Attributes)
	line := DecodeTileLine(sf.TileMSB, sf.TileLSB)
	if obj.Attributes&Bit5 != 0 {
		for i := range 4 {
			line[i], line[7-i] = line[7-i], line[i]
		}
	}
	for i := range 8 {
		line[i].Palette = palette
	}
	if !sf.ppu.OBJEnable() {
		for i := range 8 {
			line[i].ColorIDXBGPriority = 0
		}
	}
	if obj.Attributes&Bit7 != 0 {
		for i := range 8 {
			line[i].SetBGPriority()
		}
	}
	offsetInSprite := sf.ppu.Shifter.X + 8 - obj.X
	pixelsToPush := 8 - offsetInSprite
	pos := sf.ppu.SpriteFIFO.ShiftPos
	for i := range int(pixelsToPush) {
		incLevel := i >= sf.ppu.SpriteFIFO.Level
		pushPixel := incLevel || sf.ppu.SpriteFIFO.Slots[pos].ColorIDX() == 0
		if pushPixel {
			pixel := line[int(offsetInSprite)+i]
			sf.ppu.SpriteFIFO.Slots[pos] = pixel
		}
		if incLevel {
			sf.ppu.SpriteFIFO.Level++
		}
		pos++
		if pos == 8 {
			pos = 0
		}
	}
	return true
}
