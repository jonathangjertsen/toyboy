package model

import (
	"fmt"
	"slices"
)

//go:generate go-enum --marshal --flag --values --nocomments

var ppuDebugEvents = []string{
	"SetLCDC",
	"SetSTAT",
	//"SetSCY",
	"SetSCX",
	"SetWY",
	"SetWX",
	"SetLY",
	"SetLYC",
	"SetBGP",
	"SetOBP0",
	"SetOBP1",

	// "BeginHBlank",
	// "BeginVBlank",
}

// ENUM(HBlank, VBlank, OAMScan, PixelDraw)
type PPUMode uint8

// ENUM(WhiteOrTransparent, LightGray, DarkGray, Black)
type Color uint8

// ENUM(FetchTileNo, FetchTileLSB, FetchTileMSB, PushFIFO)
type PixelFetcherState uint8

type Pixel struct {
	Color              Color
	Palette            uint8
	SpritePriority     bool // CGB only
	BackgroundPriority bool
}

type Sprite struct {
	X          uint8
	Y          uint8
	TileIndex  uint8
	Attributes uint8
}

type PPUHooks interface {
	FrameCompleted(ViewPort)
}

type PPU struct {
	MemoryRegion

	RegLCDC uint8
	RegSTAT uint8
	RegSCY  uint8
	RegSCX  uint8
	RegWY   uint8
	RegWX   uint8
	RegLY   uint8
	RegLYC  uint8
	RegBGP  uint8
	RegOBP0 uint8
	RegOBP1 uint8

	Hooks PPUHooks

	Bus  *Bus
	Mode PPUMode

	OAMBuffer      [10]Sprite
	OAMBufferLevel int

	// OAM scan state
	OAMScanCycle uint64

	// Pixel draw state
	PixelDrawCycle           uint64
	BackgroundFetcher        BackgroundFetcher
	RemainingPixelsToDiscard uint8
	SpriteFetcher            SpriteFetcher
	PixelShifter             PixelShifter

	HBlankRemainingCycles     uint64
	VBlankLineRemainingCycles uint64

	BackgroundFIFO PixelFIFO
	SpriteFIFO     PixelFIFO

	Palette [4]Color

	FBBackground FrameBuffer
	FBWindow     FrameBuffer
	FBViewport   ViewPort
}

type FrameBuffer [256][256]Color

type ViewPort [144][160]Color

type PixelFIFO struct {
	Slots    [8]Pixel
	Level    int
	ShiftPos int
	PushPos  int
}

func (fifo *PixelFIFO) Clear() {
	clear(fifo.Slots[:])
	fifo.Level = 0
	fifo.ShiftPos = 0
	fifo.PushPos = 0
}

func (fifo *PixelFIFO) ShiftOut() (Pixel, bool) {
	var p Pixel
	if fifo.Level == 0 {
		return p, false
	}
	fifo.Level--
	p = fifo.Slots[fifo.ShiftPos]
	fifo.ShiftPos++
	if fifo.ShiftPos == 8 {
		fifo.ShiftPos = 0
	}
	return p, true
}

func (fifo *PixelFIFO) Write8(pixels [8]Pixel) {
	pos := fifo.ShiftPos
	for i := range 8 {
		fifo.Slots[pos] = pixels[i]
		pos++
		if pos == 8 {
			pos = 0
		}
	}
	fifo.Level = 8
}

func (ppu *PPU) Debug(event string, f string, v ...any) {
	if !slices.Contains(ppuDebugEvents, event) {
		return
	}
	fmt.Printf("PPU | %s | ", event)
	fmt.Printf(f, v...)
	fmt.Printf("\n")
}

func (ppu *PPU) Enabled() bool {
	return ppu.RegLCDC&0x80 != 0
}

func (ppu *PPU) WindowTilemapArea() uint16 {
	if ppu.RegLCDC&0x40 != 0 {
		return 0x9800
	}
	return 0x9c00
}

func (ppu *PPU) WindowEnable() bool {
	if ppu.RegLCDC&0x01 == 0 {
		return false // DMG only
	}
	return ppu.RegLCDC&0x20 != 0
}

func (ppu *PPU) BGTilemapArea() uint16 {
	if ppu.RegLCDC&0x08 != 0 {
		return 0x9800
	}
	return 0x9c00
}

func (ppu *PPU) ObjHeight() uint8 {
	if ppu.RegLCDC&0x04 != 0 {
		return 16
	}
	return 8
}

func (ppu *PPU) OBJEnable() uint8 {
	bitSet := ppu.RegLCDC&0x02 != 0
	_ = bitSet
	panic("not implemented")
}

func (ppu *PPU) BGWindowEnablePriority() uint8 {
	bitSet := ppu.RegLCDC&0x01 != 0
	_ = bitSet
	panic("not implemented")
}

func (ppu *PPU) SetLCDC(v uint8) {
	ppu.Debug("SetLCDC", "%02x", v)
	ppu.RegLCDC = v
}

func (ppu *PPU) SetSTAT(v uint8) {
	ppu.Debug("SetSTAT", "%02x", v)
	ppu.RegSTAT = v
	panic("not implemented: SetSTAT")
}

func (ppu *PPU) SetSCY(v uint8) {
	ppu.Debug("SetSCY", "%02x", v)
	ppu.RegSCY = v
}

func (ppu *PPU) SetSCX(v uint8) {
	ppu.Debug("SetSCX", "%02x", v)
	ppu.RegSCX = v
	panic("not implemented: SetSCX")
}

func (ppu *PPU) SetWY(v uint8) {
	ppu.Debug("SetWY", "%02x", v)
	ppu.RegWY = v
	panic("not implemented: SetWY")
}

func (ppu *PPU) SetWX(v uint8) {
	ppu.Debug("SetWX", "%02x", v)
	ppu.RegWX = v
	panic("not implemented: SetWX")
}

func (ppu *PPU) SetLY(v uint8) {
	ppu.Debug("SetLY", "%02x", v)
	ppu.RegLY = v
	panic("not implemented: SetLY")
}

func (ppu *PPU) SetLYC(v uint8) {
	ppu.Debug("SetLYC", "%02x", v)
	ppu.RegLYC = v
	panic("not implemented: SetLYC")
}

func (ppu *PPU) SetBGP(v uint8) {
	ppu.Debug("SetBGP", "%02x", v)
	ppu.RegBGP = v

	ppu.Palette[0] = Color((v >> 0) & 0x3)
	ppu.Palette[1] = Color((v >> 2) & 0x3)
	ppu.Palette[2] = Color((v >> 4) & 0x3)
	ppu.Palette[3] = Color((v >> 6) & 0x3)
}

func (ppu *PPU) SetOBP0(v uint8) {
	ppu.Debug("SetOBP0", "v=%02x", v)
	ppu.RegOBP0 = v
	panic("not implemented: SetOBP0")
}

func (ppu *PPU) SetOBP1(v uint8) {
	ppu.Debug("SetOBP1", "v=%02x", v)
	ppu.RegOBP1 = v
	panic("not implemented: SetOBP1")
}

func NewPPU(rtClock *ClockRT, clock *Clock, bus *Bus, hooks PPUHooks) *PPU {
	ppu := &PPU{
		MemoryRegion: NewMemoryRegion(rtClock, AddrPPUBegin, AddrPPUEnd),
		Bus:          bus,
		Hooks:        hooks,
	}
	ppu.BackgroundFetcher.PPU = ppu
	ppu.SpriteFetcher.PPU = ppu
	ppu.PixelShifter.PPU = ppu
	ppu.beginFrame()
	clock.AttachDevice(ppu.fsm)
	return ppu
}

func (ppu *PPU) Dump() {
	fmt.Printf("\n--------\nPPU dump:\n")
	fmt.Printf("Mode: %v\n", ppu.Mode)
	fmt.Printf("OAMScanCycle: %d\n", ppu.OAMScanCycle)
	fmt.Printf("OAMBuffer: %v\n", ppu.OAMBuffer)
	fmt.Printf("\n--------\n")
}

func (ppu *PPU) fsm(c Cycle) {
	if !ppu.Enabled() {
		return
	}
	switch ppu.Mode {
	case PPUModeVBlank:
		ppu.fsmVBlank()
	case PPUModeHBlank:
		ppu.fsmHBlank()
	case PPUModePixelDraw:
		ppu.fsmPixelDraw()
	case PPUModeOAMScan:
		ppu.fsmOAMScan()
	default:
		panicf("not implemented mode: %v", ppu.Mode)
	}
}

func (ppu *PPU) beginFrame() {
	ppu.RegLY = 0
	ppu.BackgroundFetcher.Cycle = 0
	ppu.BackgroundFetcher.State = PixelFetcherStateFetchTileNo
	ppu.BackgroundFetcher.WindowFetching = false
	ppu.BackgroundFetcher.WindowYReached = false
	ppu.beginOAMScan()
}

func (ppu *PPU) beginOAMScan() {
	ppu.Mode = PPUModeOAMScan
	ppu.OAMScanCycle = 0
	ppu.OAMBufferLevel = 0
}

// start of scanline after OAM scan
func (ppu *PPU) beginPixelDraw() {
	ppu.Mode = PPUModePixelDraw
	ppu.BackgroundFetcher.Cycle = 0
	ppu.BackgroundFetcher.State = PixelFetcherStateFetchTileNo
	ppu.BackgroundFetcher.WindowFetching = false
	ppu.BackgroundFIFO.Clear()
	ppu.SpriteFIFO.Clear()
	ppu.PixelDrawCycle = 0

	// GBEDG: The SCX register makes it possible to scroll the background on a per-pixel basis rather than a per-tile one.
	// While the per-tile-part of horizontal scrolling is handled within the fetching process,
	// the remaining scrolling is actually done at the start of a scanline while shifting pixels out of the background FIFO.
	// SCX mod 8 pixels are discarded at the start of each scanline rather than being pushed to the LCD,
	// which is also the cause of PPU Mode 3 being extended by SCX mod 8 cycles.
	ppu.RemainingPixelsToDiscard = ppu.RegSCX % 8
}

func (ppu *PPU) beginHBlank() {
	ppu.Debug("BeginHBlank", "pixelCycle=%v", ppu.PixelDrawCycle)

	ppu.Mode = PPUModeHBlank

	if ppu.PixelDrawCycle > 376 {
		panicv(ppu.PixelDrawCycle)
		ppu.HBlankRemainingCycles = 0
	}
	ppu.HBlankRemainingCycles = 376 - ppu.PixelDrawCycle
}

func (ppu *PPU) beginVBlank() {
	ppu.Debug("BeginVBlank", "")

	ppu.Mode = PPUModeVBlank

	ppu.VBlankLineRemainingCycles = 456
}

func (ppu *PPU) fsmOAMScan() {
	cycle := ppu.OAMScanCycle
	ppu.OAMScanCycle++
	if ppu.OAMScanCycle == 80 {
		ppu.beginPixelDraw()
	}

	// PPU checks the OAM entry every 2 cycles
	if cycle&1 == 0 {
		return
	}
	index := uint16((cycle - 1) / 2)

	// Read sprite out of OAM
	var sprite Sprite
	sprite.Y = ppu.Bus.OAM.Read(0xfe00 + index*4 + 0)
	sprite.X = ppu.Bus.OAM.Read(0xfe00 + index*4 + 1)
	sprite.TileIndex = ppu.Bus.OAM.Read(0xfe00 + index*4 + 2)
	sprite.Attributes = ppu.Bus.OAM.Read(0xfe00 + index*4 + 3)

	// Check if sprite should be added to buffer
	if len(ppu.OAMBuffer) >= 10 {
		return
	}
	if sprite.X == 0 {
		return
	}
	if ppu.RegLY+16 < sprite.Y {
		return
	}
	if ppu.RegLY+16 <= sprite.Y+ppu.ObjHeight() {
		return
	}

	// Add
	ppu.OAMBuffer[ppu.OAMBufferLevel] = sprite
	ppu.OAMBufferLevel++
}

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
	bgf.PPU.BackgroundFIFO.Write8(bgf.PPU.decodeTileLine(bgf.TileLSB, bgf.TileMSB))
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
	for idx := range sf.PPU.OAMBufferLevel {
		obj := sf.PPU.OAMBuffer[idx]
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
	sf.TileIndex = sf.PPU.OAMBuffer[sf.SpriteIDX].TileIndex
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
	line := sf.PPU.decodeTileLine(sf.TileLSB, sf.TileMSB)
	obj := sf.PPU.OAMBuffer[sf.SpriteIDX]
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

func (ppu *PPU) fsmPixelDraw() {
	ppu.SpriteFetcher.fsm()
	ppu.BackgroundFetcher.Suspended = !ppu.SpriteFetcher.Suspended
	ppu.PixelShifter.Suspended = !ppu.SpriteFetcher.Suspended

	ppu.BackgroundFetcher.fsm()

	ppu.PixelShifter.fsm()

	// GBEDG: After each pixel shifted out, the PPU checks if it has reached the window
	if !ppu.BackgroundFetcher.WindowFetching && ppu.BackgroundFetcher.windowReached() {
		ppu.BackgroundFetcher.WindowFetching = true
		ppu.BackgroundFIFO.Clear()
		ppu.BackgroundFetcher.X = 0
	}

	if ppu.PixelShifter.X >= 160 {
		ppu.beginHBlank()
		ppu.PixelShifter.X = 0
	}

	ppu.PixelDrawCycle++
}

type PixelShifter struct {
	RemainingPixelsToDiscard uint8
	Suspended                bool
	X                        uint8

	PPU *PPU
}

func (ps *PixelShifter) fsm() {
	if ps.Suspended {
		return
	}

	if ps.RemainingPixelsToDiscard > 0 {
		if _, shifted := ps.PPU.BackgroundFIFO.ShiftOut(); shifted {
			ps.RemainingPixelsToDiscard--
		}
		return
	}

	pixel, havePixel := ps.getPixel()
	if !havePixel {
		return
	}

	// Write pixel to LCD
	ps.PPU.FBViewport[ps.PPU.RegLY][ps.X] = pixel.Color
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

func (ppu *PPU) decodeTileLine(lsb, msb uint8) [8]Pixel {
	var pixels [8]Pixel
	for i := 0; i < 8; i++ {
		lsbMask := lsb & (1 << (8 - i))
		lsbMask >>= (8 - i)
		msbMask := msb & (1 << (8 - i))
		msbMask >>= (8 - i)
		msbMask <<= 1
		pixels[i].Color = Color(lsbMask | msbMask)
	}
	return pixels
}

func (ppu *PPU) fsmVBlank() {
	if ppu.VBlankLineRemainingCycles > 0 {
		ppu.VBlankLineRemainingCycles--
		return
	}

	ppu.Hooks.FrameCompleted(ppu.FBViewport)

	if ppu.RegLY == 153 {
		ppu.RegLY = 0
		ppu.beginFrame()
	} else {
		ppu.RegLY++
	}
}

func (ppu *PPU) fsmHBlank() {
	if ppu.HBlankRemainingCycles > 0 {
		ppu.HBlankRemainingCycles--
		return
	}

	ppu.RegLY++
	if ppu.BackgroundFetcher.WindowFetching {
		ppu.BackgroundFetcher.WindowLineCounter++
	}

	if ppu.RegLY < 144 {
		ppu.beginOAMScan()
	} else if ppu.RegLY == 144 {
		ppu.beginVBlank()
	} else {
		panicv(ppu.RegLY)
	}
}

func (ppu *PPU) Read(addr uint16) uint8 {
	switch Addr(addr) {
	case AddrLCDC:
		return ppu.RegLCDC
	case AddrSTAT:
		return ppu.RegSTAT
	case AddrSCY:
		return ppu.RegSCY
	case AddrSCX:
		return ppu.RegSCX
	case AddrLY:
		return ppu.RegLY
	case AddrLYC:
		return ppu.RegLYC
	case AddrBGP:
		return ppu.RegBGP
	case AddrOBP0:
		return ppu.RegOBP0
	case AddrOBP1:
		return ppu.RegOBP1
	case AddrWY:
		return ppu.RegWY
	case AddrWX:
		return ppu.RegWX
	}
	return 0
}

func (ppu *PPU) Write(addr uint16, v uint8) {
	switch Addr(addr) {
	case AddrLCDC:
		ppu.SetLCDC(v)
	case AddrSTAT:
		ppu.SetSTAT(v)
	case AddrSCY:
		ppu.SetSCY(v)
	case AddrSCX:
		ppu.SetSCX(v)
	case AddrLY:
		ppu.SetLY(v)
	case AddrLYC:
		ppu.SetLYC(v)
	case AddrBGP:
		ppu.SetBGP(v)
	case AddrOBP0:
		ppu.SetOBP0(v)
	case AddrOBP1:
		ppu.SetOBP1(v)
	case AddrWY:
		ppu.SetWY(v)
	case AddrWX:
		ppu.SetWX(v)
	default:
		panicf("Write to unknown LCD register %#v", addr)
	}
}
