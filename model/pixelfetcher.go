package model

type PixelFetcher struct {
	Cycle       uint64
	PPU         *PPU
	State       PixelFetcherState
	X           uint8
	TileIndex   uint8
	TileLSBAddr uint16
	TileLSB     uint8
	TileMSB     uint8
	Suspended   bool
}

type BackgroundFetcher struct {
	PixelFetcher

	WindowYReached    bool
	WindowFetching    bool
	WindowLineCounter uint16
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
		}
	}
}

func (bgf *BackgroundFetcher) fetchTileNo() {
	// GBEDG: During the first step the fetcher fetches and stores the tile number of the tile which should be used.
	// Which Tilemap is used depends on whether the PPU is currently rendering Background or Window pixels
	// and on the bits 3 and 5 of the LCDC register.
	var addr uint16
	if bgf.WindowFetching {
		addr = bgf.PPU.WindowTilemapArea()
	} else {
		addr = bgf.PPU.BGTilemapArea()
	}

	// GBEDG: Additionally, the address which the tile number is read from is offset by the fetcher-internal X-Position-Counter, which is incremented each time the last step is completed.
	// GBEDG: The value of SCX / 8 is also added if the Fetcher is not fetching Window pixels. In order to make the wrap-around with SCX work, this offset is ANDed with 0x1f
	// GBEDG: Note: The sum of [...] the X-POS+SCX [...] is ANDed with 0x3ff in order to ensure that the address stays within the Tilemap memory regions.
	offsetX := (uint16(bgf.X) + (uint16(bgf.PPU.RegSCX/8) & 0x1f)) & 0x3ff
	addr += offsetX

	// GBEDG: An offset of 32 * (((LY + SCY) & 0xFF) / 8) is also added if background pixels are being fetched, otherwise, if window pixels are being fetched, this offset is determined by 32 * (WINDOW_LINE_COUNTER / 8)
	var offsetY uint16
	if bgf.WindowFetching {
		offsetY = 32 * (bgf.WindowLineCounter / 8)
	} else {
		offsetY = 32 * (uint16((bgf.PPU.RegLY+bgf.PPU.RegSCY)&0xff) / 8)
	}
	offsetY &= 0x3ff
	addr += offsetY

	bgf.TileIndex = bgf.PPU.Bus.VRAM.Read(addr)
}

func (bgf *BackgroundFetcher) fetchTileLSB() {
	var addr uint16
	if bgf.PPU.RegLCDC&0x10 == 0 {
		addr = 0x8000 + 16*uint16(bgf.TileIndex)
	} else {
		if bgf.TileIndex >= 128 {
			addr = 0x8800 + 16*uint16(bgf.TileIndex-128)
		} else {
			addr = 0x9000 + 16*uint16(bgf.TileIndex)
		}
	}
	if bgf.WindowFetching {
		addr += 2 * uint16(bgf.WindowLineCounter%8)
	} else {
		addr += 2 * uint16((bgf.PPU.RegLY+bgf.PPU.RegSCY)%8)
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
	bgf.PPU.BackgroundFIFO.Write8(decodeTileLine(bgf.TileLSB, bgf.TileMSB))
	return true
}

type SpriteFetcher struct {
	PixelFetcher
	Fetching  bool
	SpriteIDX int
	TileNo    int
}

func (sf *SpriteFetcher) fsm() {
	sfCycle := sf.Cycle
	sf.Cycle++

	// If the X-Position of any sprite in the sprite buffer is less than or equal to the current Pixel-X-Position + 8,
	// a sprite fetch is initiated
	sf.Suspended = true
	for idx := range sf.PPU.OAMBuffer.Level {
		obj := sf.PPU.OAMBuffer.Buffer[idx]
		doSpriteFetch := obj.X <= sf.X+8

		// presumably also...?
		if obj.X < sf.X {
			doSpriteFetch = false
		}
		if obj.Y < sf.PPU.RegLY {
			doSpriteFetch = false
		}
		if obj.Y > sf.PPU.RegLY+8 { // TODO: tall sprites
			doSpriteFetch = false
		}

		if doSpriteFetch {
			sf.Suspended = false
			sf.SpriteIDX = idx
			break
		}
	}

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
		}
	}
}

func (sf *SpriteFetcher) fetchTileNo() {
	sf.TileIndex = sf.PPU.OAMBuffer.Buffer[sf.SpriteIDX].TileIndex
}

func (sf *SpriteFetcher) fetchTileLSB() {
	addr := 0x8000 + 16*uint16(sf.TileIndex)
	// TODO: what about window?
	addr += 2 * uint16((sf.PPU.RegLY+sf.PPU.RegSCY)%8)
	sf.TileLSBAddr = addr
	sf.TileLSB = sf.PPU.Bus.VRAM.Read(addr)
}

func (sf *SpriteFetcher) fetchTileMSB() {
	sf.TileMSB = sf.PPU.Bus.VRAM.Read(sf.TileLSBAddr + 1)
}

func (sf *SpriteFetcher) pushFIFO() bool {
	line := decodeTileLine(sf.TileLSB, sf.TileMSB)
	obj := sf.PPU.OAMBuffer.Buffer[sf.SpriteIDX]
	offset := obj.X - sf.X + uint8(sf.PPU.SpriteFIFO.Level)
	pos := sf.PPU.SpriteFIFO.ShiftPos
	for range offset {
		pos++
		if pos == 8 {
			pos = 0
		}
	}
	for i := offset; i < 8; i++ {
		sf.PPU.SpriteFIFO.Slots[pos] = line[i]
		pos++
		if pos == 8 {
			pos = 0
		}
	}
	sf.PPU.SpriteFIFO.Level = 8
	return true
}
