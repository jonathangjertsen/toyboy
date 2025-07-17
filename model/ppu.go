package model

import (
	"fmt"
	"image/color"
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

func (c Color) Grayscale() uint8 {
	switch c {
	case ColorWhiteOrTransparent:
		return 0xf0
	case ColorLightGray:
		return 0xa0
	case ColorDarkGray:
		return 0x70
	case ColorBlack:
		return 0x30
	}
	return 0x00
}

func (c Color) RGBA() color.RGBA {
	switch c {
	case ColorWhiteOrTransparent:
		return color.RGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff}
	case ColorLightGray:
		return color.RGBA{R: 0xa0, G: 0xa0, B: 0xa0, A: 0xff}
	case ColorDarkGray:
		return color.RGBA{R: 0x70, G: 0x70, B: 0x70, A: 0xff}
	case ColorBlack:
		return color.RGBA{R: 0x30, G: 0x30, B: 0x30, A: 0xff}
	}
	return color.RGBA{}
}

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

func DecodeSprite(data []uint8) Sprite {
	return Sprite{
		Y:          data[0],
		X:          data[1],
		TileIndex:  data[2],
		Attributes: data[3],
	}
}

type PPU struct {
	MemoryRegion

	Debugger *Debugger

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

	Bus  *Bus
	Mode PPUMode

	OAMBuffer OAMBuffer

	// OAM scan state
	OAMScanCycle uint64

	// Pixel draw state
	PixelDrawCycle    uint64
	BackgroundFetcher BackgroundFetcher
	SpriteFetcher     SpriteFetcher
	PixelShifter      PixelShifter

	HBlankRemainingCycles     uint64
	VBlankLineRemainingCycles uint64

	BackgroundFIFO PixelFIFO
	SpriteFIFO     PixelFIFO

	BGPalette   [4]Color
	OBJPalette0 [4]Color
	OBJPalette1 [4]Color

	FBBackground FrameBuffer
	FBWindow     FrameBuffer
	FBViewport   ViewPort
	LastFrame    ViewPort
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
		return 0x9c00
	}
	return 0x9800
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
	ppu.RegSTAT = maskedWrite(ppu.RegSTAT, v, 0xf8)
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

	ppu.BGPalette[0] = Color((v >> 0) & 0x3)
	ppu.BGPalette[1] = Color((v >> 2) & 0x3)
	ppu.BGPalette[2] = Color((v >> 4) & 0x3)
	ppu.BGPalette[3] = Color((v >> 6) & 0x3)
}

func (ppu *PPU) SetOBP0(v uint8) {
	ppu.Debug("SetOBP0", "v=%02x", v)
	ppu.RegOBP0 = v

	ppu.OBJPalette0[0] = Color((v >> 0) & 0x3)
	ppu.OBJPalette0[1] = Color((v >> 2) & 0x3)
	ppu.OBJPalette0[2] = Color((v >> 4) & 0x3)
	ppu.OBJPalette0[3] = Color((v >> 6) & 0x3)
}

func (ppu *PPU) SetOBP1(v uint8) {
	ppu.Debug("SetOBP1", "v=%02x", v)
	ppu.RegOBP1 = v

	ppu.OBJPalette1[0] = Color((v >> 0) & 0x3)
	ppu.OBJPalette1[1] = Color((v >> 2) & 0x3)
	ppu.OBJPalette1[2] = Color((v >> 4) & 0x3)
	ppu.OBJPalette1[3] = Color((v >> 6) & 0x3)
}

func NewPPU(rtClock *ClockRT, clock *Clock, bus *Bus, debugger *Debugger) *PPU {
	ppu := &PPU{
		MemoryRegion: NewMemoryRegion(rtClock, AddrPPUBegin, AddrPPUEnd),
		Bus:          bus,
		Debugger:     debugger,
	}
	ppu.BackgroundFetcher.PPU = ppu
	ppu.SpriteFetcher.PPU = ppu
	ppu.SpriteFetcher.Suspended = true
	ppu.SpriteFetcher.DoneX = 0xff
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

func (ppu *PPU) setMode(mode PPUMode) {
	ppu.Mode = mode
	ppu.RegSTAT = maskedWrite(ppu.RegSTAT, uint8(mode), 0x7)
}

func (ppu *PPU) beginFrame() {
	ppu.beginOAMScan()
}

func (ppu *PPU) beginOAMScan() {
	ppu.setMode(PPUModeOAMScan)
	ppu.OAMScanCycle = 0
	ppu.OAMBuffer.Level = 0
}

// start of scanline after OAM scan
func (ppu *PPU) beginPixelDraw() {
	ppu.setMode(PPUModePixelDraw)
	ppu.BackgroundFetcher.Cycle = 0
	ppu.BackgroundFetcher.State = PixelFetcherStateFetchTileNo
	ppu.BackgroundFetcher.WindowFetching = false
	ppu.BackgroundFetcher.X = 0
	ppu.SpriteFetcher.Cycle = 0
	ppu.SpriteFetcher.DoneX = 0xff
	ppu.SpriteFetcher.X = 0
	ppu.BackgroundFIFO.Clear()
	ppu.SpriteFIFO.Clear()
	ppu.PixelDrawCycle = 0

	// GBEDG: The SCX register makes it possible to scroll the background on a per-pixel basis rather than a per-tile one.
	// While the per-tile-part of horizontal scrolling is handled within the fetching process,
	// the remaining scrolling is actually done at the start of a scanline while shifting pixels out of the background FIFO.
	// SCX mod 8 pixels are discarded at the start of each scanline rather than being pushed to the LCD,
	// which is also the cause of PPU Mode 3 being extended by SCX mod 8 cycles.
	ppu.PixelShifter.Discard = ppu.RegSCX % 8
}

func (ppu *PPU) beginHBlank() {
	ppu.Debug("BeginHBlank", "pixelCycle=%v", ppu.PixelDrawCycle)

	ppu.setMode(PPUModeHBlank)

	if ppu.PixelDrawCycle > 376 {
		panicv(ppu.PixelDrawCycle)
		ppu.HBlankRemainingCycles = 0
	}
	ppu.HBlankRemainingCycles = 376 - ppu.PixelDrawCycle
}

func (ppu *PPU) beginVBlank() {
	ppu.Debug("BeginVBlank", "")

	ppu.setMode(PPUModeVBlank)

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
	spriteData := make([]uint8, 4)
	for offs := uint16(0); offs < 4; offs++ {
		spriteData[offs] = ppu.Bus.OAM.Read(AddrOAMBegin + index*4 + offs)
	}
	sprite := DecodeSprite(spriteData)

	// Check if sprite should be added to buffer
	if !(sprite.X > 0) {
		return
	}
	if !(ppu.RegLY+16 >= sprite.Y) {
		return
	}
	if !(ppu.RegLY+16 < sprite.Y+ppu.ObjHeight()) {
		return
	}
	if ppu.OAMBuffer.Full() {
		return
	}

	ppu.OAMBuffer.Add(sprite)
}

func (ppu *PPU) fsmPixelDraw() {
	if ppu.SpriteFetcher.Suspended && ppu.SpriteFetcher.DoneX != ppu.PixelShifter.X {
		for idx := range ppu.OAMBuffer.Level {
			obj := ppu.OAMBuffer.Buffer[idx]
			if obj.X <= ppu.PixelShifter.X+8 && obj.X > ppu.PixelShifter.X {
				// Initiate sprite fetch
				ppu.SpriteFetcher.State = 1
				ppu.SpriteFetcher.SpriteIDX = idx
				ppu.PixelShifter.Suspended = true
				ppu.SpriteFetcher.Suspended = false
				ppu.SpriteFetcher.DoneX = 0xff
				break
			}
		}
	} else {
		if ppu.SpriteFetcher.DoneX != 0xff {
			ppu.BackgroundFetcher.Suspended = false
			ppu.PixelShifter.Suspended = false
		}
	}

	ppu.SpriteFetcher.fsm()
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

func (ppu *PPU) fsmVBlank() {
	if ppu.VBlankLineRemainingCycles > 0 {
		ppu.VBlankLineRemainingCycles--
		return
	}

	ppu.LastFrame = ppu.FBViewport
	ppu.IncRegLY()
	if ppu.RegLY == 0 {
		ppu.beginFrame()
	}
}

func (ppu *PPU) fsmHBlank() {
	if ppu.HBlankRemainingCycles > 0 {
		ppu.HBlankRemainingCycles--
		return
	}

	ppu.IncRegLY()
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

func (ppu *PPU) IncRegLY() {
	ppu.RegLY++
	if ppu.RegLY >= 153 {
		ppu.RegLY = 0
	}
	if ppu.RegLY == ppu.RegLYC {
		ppu.RegSTAT |= uint8(1 << 2)
	} else {
		ppu.RegSTAT &= ^uint8(1 << 2)
	}
	ppu.Debugger.SetY(ppu.RegLY)
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
