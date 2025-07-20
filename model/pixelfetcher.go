package model

type PixelFetcher struct {
	Cycle       uint64
	State       PixelFetcherState
	X           Data8
	TileIndex   Data8
	TileLSBAddr Addr
	TileLSB     Data8
	TileMSB     Data8
	Suspended   bool

	PPU *PPU
}

type BackgroundFetcher struct {
	PixelFetcher
	TileIndexAddr Addr

	TileOffsetX Addr
	TileOffsetY Addr

	WindowYReached    bool
	WindowFetching    bool
	WindowLineCounter Data8
}

func (bgf *BackgroundFetcher) fsm() {
	if bgf.Suspended {
		return
	}

	bgfCycle := bgf.Cycle
	bgf.Cycle++

	switch bgf.State {
	case PixelFetcherStateFetchTileNo:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileNo()
		bgf.State = PixelFetcherStateFetchTileLSB
	case PixelFetcherStateFetchTileLSB:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileLSB()
		bgf.State = PixelFetcherStateFetchTileMSB
	case PixelFetcherStateFetchTileMSB:
		// Takes 2 cycles
		if bgfCycle&1 == 0 {
			return
		}
		bgf.fetchTileMSB()
		bgf.State = PixelFetcherStatePushFIFO
	case PixelFetcherStatePushFIFO:
		if bgf.pushFIFO() {
			bgf.State = PixelFetcherStateFetchTileNo
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
		addr = bgf.PPU.WindowTilemapArea()
	} else {
		addr = bgf.PPU.BGTilemapArea()
	}

	// GBEDG: Additionally, the address which the tile number is read from is offset by the fetcher-internal X-Position-Counter,
	//        which is incremented each time the last step is completed.
	offsetX := Addr(bgf.X)
	if !bgf.WindowFetching {
		// GBEDG: The value of SCX / 8 is also added if the Fetcher is not fetching Window pixels.
		//        In order to make the wrap-around with SCX work, this offset is ANDed with 0x1f
		offsetX += Addr((bgf.PPU.RegSCX / 8))
		offsetX &= 0x1f
	}

	var offsetY Addr
	if !bgf.WindowFetching {
		// GBEDG: An offset of 32 * (((LY + SCY) & 0xFF) / 8) is also added if background pixels are being fetched,
		offs := (bgf.PPU.RegLY + bgf.PPU.RegSCY) & 0xff
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
	bgf.TileIndex = bgf.PPU.Bus.VRAM.Read(addr)
}

func (bgf *BackgroundFetcher) fetchTileLSB() {
	idx := bgf.TileIndex
	signedAddressing := !bgf.PPU.RegLCDC.Bit(4)
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
		addr += 2 * Addr((bgf.PPU.RegLY+bgf.PPU.RegSCY)%8)
	}
	bgf.TileLSBAddr = addr
	bgf.TileLSB = bgf.PPU.Bus.VRAM.Read(addr)
}

func (bgf *BackgroundFetcher) fetchTileMSB() {
	bgf.TileMSB = bgf.PPU.Bus.VRAM.Read(bgf.TileLSBAddr + 1)
}

func (bgf *BackgroundFetcher) windowReached() bool {
	if bgf.PPU.RegWY == bgf.PPU.RegLY {
		bgf.WindowYReached = true
	}
	if !bgf.WindowYReached {
		return false
	}
	if !bgf.PPU.WindowEnable() {
		return false
	}
	if bgf.X < bgf.PPU.RegWX-7 {
		return false
	}
	return true
}

func (bgf *BackgroundFetcher) pushFIFO() bool {
	if bgf.PPU.BackgroundFIFO.Level > 0 {
		return false
	}
	bgf.PPU.BackgroundFIFO.Write8(DecodeTileLine(bgf.TileMSB, bgf.TileLSB))
	return true
}

type SpriteFetcher struct {
	PixelFetcher
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
	case PixelFetcherStateFetchTileNo:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileNo()
		sf.State = PixelFetcherStateFetchTileLSB
	case PixelFetcherStateFetchTileLSB:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileLSB()
		sf.State = PixelFetcherStateFetchTileMSB
	case PixelFetcherStateFetchTileMSB:
		// Takes 2 cycles
		if sfCycle&1 == 0 {
			return
		}
		sf.fetchTileMSB()
		sf.State = PixelFetcherStatePushFIFO
	case PixelFetcherStatePushFIFO:
		if sf.pushFIFO() {
			sf.State = PixelFetcherStateFetchTileNo
			sf.Suspended = true
			sf.DoneX = sf.PPU.PixelShifter.X
		}
	}
}

func (sf *SpriteFetcher) fetchTileNo() {
	sf.TileIndex = sf.PPU.OAMBuffer.Buffer[sf.SpriteIDX].TileIndex
}

func (sf *SpriteFetcher) fetchTileLSB() {
	obj := sf.PPU.OAMBuffer.Buffer[sf.SpriteIDX]
	addr := 0x8000 + 16*Addr(sf.TileIndex)
	addr += 2 * Addr((sf.PPU.RegLY+sf.PPU.RegSCY-obj.Y)%8)
	sf.TileLSBAddr = addr
	sf.TileLSB = sf.PPU.Bus.VRAM.Read(addr)
}

func (sf *SpriteFetcher) fetchTileMSB() {
	sf.TileMSB = sf.PPU.Bus.VRAM.Read(sf.TileLSBAddr + 1)
}

func (sf *SpriteFetcher) pushFIFO() bool {
	line := DecodeTileLine(sf.TileMSB, sf.TileLSB)
	obj := sf.PPU.OAMBuffer.Buffer[sf.SpriteIDX]
	offsetInSprite := sf.PPU.PixelShifter.X + 8 - obj.X
	pixelsToPush := 8 - offsetInSprite
	pos := sf.PPU.SpriteFIFO.ShiftPos
	for i := range int(pixelsToPush) {
		incLevel := i >= sf.PPU.SpriteFIFO.Level
		pushPixel := incLevel || sf.PPU.SpriteFIFO.Slots[pos].Color == ColorWhiteOrTransparent
		if pushPixel {
			pixel := line[int(offsetInSprite)+i]
			sf.PPU.SpriteFIFO.Slots[pos] = pixel
		}
		if incLevel {
			sf.PPU.SpriteFIFO.Level++
		}
		pos++
		if pos == 8 {
			pos = 0
		}
	}
	return true
}
