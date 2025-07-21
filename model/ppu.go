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
type PixelFetcherState Data8

type Pixel struct {
	Color              Color
	Palette            Data8
	SpritePriority     bool // CGB only
	BackgroundPriority bool
}

type Sprite struct {
	X          Data8
	Y          Data8
	TileIndex  Data8
	Attributes Data8
}

func DecodeSprite(data []Data8) Sprite {
	return Sprite{
		Y:          data[0],
		X:          data[1],
		TileIndex:  data[2],
		Attributes: data[3],
	}
}

type PPU struct {
	MemoryRegion

	Interrupts *Interrupts
	Debugger   *Debugger

	FrameClock *Clock
	FrameCount uint64

	RegLCDC Data8
	Stat    Stat
	RegSCY  Data8
	RegSCX  Data8
	RegWY   Data8
	RegWX   Data8
	RegLY   Data8
	RegLYC  Data8
	RegBGP  Data8
	RegOBP0 Data8
	RegOBP1 Data8

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

func (ppu *PPU) Reset() {
	ppu.FrameCount = 0
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
	return ppu.RegLCDC.Bit(7)
}

func (ppu *PPU) WindowTilemapArea() Addr {
	if ppu.RegLCDC.Bit(6) {
		return AddrTileMap1Begin
	}
	return AddrTileMap0Begin
}

func (ppu *PPU) WindowEnable() bool {
	if ppu.RegLCDC.Bit(0) {
		return false // DMG only
	}
	return ppu.RegLCDC.Bit(5)
}

func (ppu *PPU) BGTilemapArea() Addr {
	if ppu.RegLCDC.Bit(3) {
		return AddrTileMap1Begin
	}
	return AddrTileMap0Begin
}

func (ppu *PPU) ObjHeight() Data8 {
	if ppu.RegLCDC.Bit(2) {
		return 16
	}
	return 8
}

func (ppu *PPU) OBJEnable() uint8 {
	bitSet := ppu.RegLCDC.Bit(1)
	_ = bitSet
	panic("not implemented")
}

func (ppu *PPU) BGWindowEnablePriority() uint8 {
	bitSet := ppu.RegLCDC.Bit(0)
	_ = bitSet
	panic("not implemented")
}

func (ppu *PPU) SetLCDC(v Data8) {
	ppu.Debug("SetLCDC", "%s", v.Hex())
	ppu.RegLCDC = v
}

func (ppu *PPU) SetSCY(v Data8) {
	ppu.Debug("SetSCY", "%s", v.Hex())
	ppu.RegSCY = v
}

func (ppu *PPU) SetSCX(v Data8) {
	ppu.Debug("SetSCX", "%s", v.Hex())
	ppu.RegSCX = v
}

func (ppu *PPU) SetWY(v Data8) {
	ppu.Debug("SetWY", "%s", v.Hex())
	ppu.RegWY = v
	panic("not implemented: SetWY")
}

func (ppu *PPU) SetWX(v Data8) {
	ppu.Debug("SetWX", "%s", v.Hex())
	ppu.RegWX = v
	panic("not implemented: SetWX")
}

func (ppu *PPU) SetLY(v Data8) {
	ppu.Debug("SetLY", "%s", v.Hex())
	ppu.RegLY = v
	panic("not implemented: SetLY")
}

func (ppu *PPU) SetLYC(v Data8) {
	ppu.Debug("SetLYC", "%s", v.Hex())
	ppu.RegLYC = v
	panic("not implemented: SetLYC")
}

func (ppu *PPU) SetBGP(v Data8) {
	ppu.Debug("SetBGP", "%s", v.Hex())
	ppu.RegBGP = v

	ppu.BGPalette[0] = Color((v >> 0) & 0x3)
	ppu.BGPalette[1] = Color((v >> 2) & 0x3)
	ppu.BGPalette[2] = Color((v >> 4) & 0x3)
	ppu.BGPalette[3] = Color((v >> 6) & 0x3)
}

func (ppu *PPU) SetOBP0(v Data8) {
	ppu.Debug("SetOBP0", "%s", v.Hex())
	ppu.RegOBP0 = v

	ppu.OBJPalette0[0] = Color((v >> 0) & 0x3)
	ppu.OBJPalette0[1] = Color((v >> 2) & 0x3)
	ppu.OBJPalette0[2] = Color((v >> 4) & 0x3)
	ppu.OBJPalette0[3] = Color((v >> 6) & 0x3)
}

func (ppu *PPU) SetOBP1(v Data8) {
	ppu.Debug("SetOBP1", "%s", v.Hex())
	ppu.RegOBP1 = v

	ppu.OBJPalette1[0] = Color((v >> 0) & 0x3)
	ppu.OBJPalette1[1] = Color((v >> 2) & 0x3)
	ppu.OBJPalette1[2] = Color((v >> 4) & 0x3)
	ppu.OBJPalette1[3] = Color((v >> 6) & 0x3)
}

func NewPPU(rtClock *ClockRT, clock *Clock, interrupts *Interrupts, bus *Bus, debugger *Debugger) *PPU {
	ppu := &PPU{
		MemoryRegion: NewMemoryRegion(rtClock, AddrPPUBegin, SizePPU),
		Bus:          bus,
		Debugger:     debugger,
		Interrupts:   interrupts,
		Stat:         Stat{Interrupts: interrupts},
		FrameClock:   NewClock(),
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
	ppu.Stat.SetMode(mode)
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
	ppu.SpriteFetcher.State = PixelFetcherStateFetchTileNo
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

	// TODO: do we ever clear the VBlank interrupt?
	ppu.Interrupts.IRQSet(0x01)

	ppu.FrameClock.Cycle(ppu.FrameCount)
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
	index := Addr((cycle - 1) / 2)

	// Read sprite out of OAM
	spriteData := make([]Data8, 4)
	for offs := Addr(0); offs < 4; offs++ {
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
				ppu.SpriteFetcher.State = PixelFetcherStateFetchTileNo
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
	ppu.Stat.SetLYCEqLY(ppu.RegLY == ppu.RegLYC)
	ppu.Debugger.SetY(ppu.RegLY)
}

func (ppu *PPU) Read(addr Addr) Data8 {
	switch Addr(addr) {
	case AddrLCDC:
		return ppu.RegLCDC
	case AddrSTAT:
		return ppu.Stat.Reg
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

func (ppu *PPU) Write(addr Addr, v Data8) {
	switch Addr(addr) {
	case AddrLCDC:
		ppu.SetLCDC(v)
	case AddrSTAT:
		ppu.Stat.Write(v)
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
