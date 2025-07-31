package model

//go:generate go-enum --marshal --flag --values --nocomments

// ENUM(HBlank, VBlank, OAMScan, PixelDraw)
type PPUMode uint8

type PPU struct {
	// Registers and subsystems
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
	DMA     DMA

	// For other systems to hook in
	FrameCount uint

	// PPU overall state
	Mode PPUMode

	// OAM scan state
	OAMScanCycle uint64
	OAMBuffer    OAMBuffer

	// Pixel draw state
	PixelDrawCycle    uint64
	BackgroundFetcher BackgroundFetcher
	SpriteFetcher     SpriteFetcher
	Shifter           Shifter
	BackgroundFIFO    FIFO
	SpriteFIFO        FIFO
	BGPalette         Data8
	OBJPalette0       Data8
	OBJPalette1       Data8

	// HBlank/VBlank state
	HBlankRemainingCycles     uint64
	VBlankLineRemainingCycles uint64

	// Outputs
	FBViewport ViewPort
}

type FrameSync struct {
	ch chan func(*ViewPort)
}

func NewFrameSync() *FrameSync {
	return &FrameSync{
		ch: make(chan func(*ViewPort), 1),
	}
}

func NewPPU(rtClock *ClockRT, ints *Interrupts) *PPU {
	ppu := &PPU{}
	ppu.SpriteFetcher.Suspended = true
	ppu.SpriteFetcher.DoneX = 0xff

	ppu.beginFrame(ints)
	rtClock.ppu = ppu
	return ppu
}

func (ppu *PPU) GetDump() PPUDump {
	var dump PPUDump
	dump.BGFIFO = ppu.BackgroundFIFO.Dump()
	dump.SpriteFIFO = ppu.SpriteFIFO.Dump()
	dump.LastShifted = ppu.Shifter.LastShifted
	dump.OAMScanCycle = ppu.OAMScanCycle
	dump.PixelDrawCycle = ppu.PixelDrawCycle
	dump.HBlankRemainingCycles = ppu.HBlankRemainingCycles
	dump.VBlankLineRemainingCycles = ppu.VBlankLineRemainingCycles
	dump.PixelShifter = ppu.Shifter
	dump.BackgroundFetcher = ppu.BackgroundFetcher
	dump.SpriteFetcher = ppu.SpriteFetcher
	dump.OAMBuffer = ppu.OAMBuffer
	return dump
}

func (ppu *PPU) Sync(fs *FrameSync, f func(*ViewPort)) {
	done := make(chan struct{})
	fs.ch <- func(vp *ViewPort) {
		f(vp)
		done <- struct{}{}
	}
	<-done
}

func (ppu *PPU) Reset() {
	ppu.FrameCount = 0
}

func (ppu *PPU) WindowTilemapArea() Addr {
	if ppu.RegLCDC&Bit6 != 0 {
		return AddrTileMap1Begin
	}
	return AddrTileMap0Begin
}

func (ppu *PPU) WindowEnable() bool {
	if ppu.RegLCDC&Bit0 == 0 {
		return false // DMG only
	}
	return ppu.RegLCDC&Bit5 != 0
}

func (ppu *PPU) BGTilemapArea() Addr {
	if ppu.RegLCDC&Bit3 != 0 {
		return AddrTileMap1Begin
	}
	return AddrTileMap0Begin
}

func (ppu *PPU) ObjHeight() Data8 {
	if ppu.RegLCDC&Bit2 != 0 {
		return 16
	}
	return 8
}

func (ppu *PPU) OBJEnable() bool {
	return ppu.RegLCDC&Bit1 != 0
}

func (ppu *PPU) BGWindowEnable() bool {
	return ppu.RegLCDC&Bit0 != 0
}

func (ppu *PPU) SetLCDC(v Data8) {
	ppu.RegLCDC = v
}

func (ppu *PPU) SetSCY(v Data8) {
	ppu.RegSCY = v
}

func (ppu *PPU) SetSCX(v Data8) {
	ppu.RegSCX = v
}

func (ppu *PPU) SetWY(v Data8) {
	ppu.RegWY = v
}

func (ppu *PPU) SetWX(v Data8) {
	ppu.RegWX = v
}

func (ppu *PPU) SetLY(v Data8) {
	ppu.RegLY = v
}

func (ppu *PPU) SetLYC(ints *Interrupts, v Data8) {
	ppu.RegLYC = v
	ints.IRQCheck()
}

func (ppu *PPU) SetBGP(v Data8) {
	ppu.RegBGP = v
	ppu.BGPalette = v
}

func (ppu *PPU) SetOBP0(v Data8) {
	ppu.RegOBP0 = v
	ppu.OBJPalette0 = v
}

func (ppu *PPU) SetOBP1(v Data8) {
	ppu.RegOBP1 = v
	ppu.OBJPalette1 = v
}

func (ppu *PPU) fsm(ints *Interrupts, debug *Debug, clk *ClockRT, mem []Data8, fs *FrameSync) {
	if ppu.DMA.Source != 0 {
		ppu.DMA.fsm(mem)
	}
	switch ppu.Mode {
	case PPUModeVBlank:
		ppu.fsmVBlank(ints, debug, fs)
	case PPUModeHBlank:
		ppu.fsmHBlank(ints, debug)
	case PPUModePixelDraw:
		ppu.fsmPixelDraw(ints, debug, clk, mem)
	case PPUModeOAMScan:
		ppu.fsmOAMScan(ints, mem)
	}
}

func (ppu *PPU) setMode(ints *Interrupts, mode PPUMode) {
	ppu.Mode = mode
	ppu.Stat.SetMode(ints, mode)
}

func (ppu *PPU) beginFrame(ints *Interrupts) {
	ppu.beginOAMScan(ints)
	ppu.BackgroundFetcher.WindowYReached = false
}

func (ppu *PPU) beginOAMScan(ints *Interrupts) {
	ppu.setMode(ints, PPUModeOAMScan)
	ppu.OAMScanCycle = 0
	ppu.OAMBuffer.Level = 0
}

// start of scanline after 7OAM scan
func (ppu *PPU) beginPixelDraw(ints *Interrupts) {
	ppu.setMode(ints, PPUModePixelDraw)
	ppu.BackgroundFetcher.Cycle = 0
	ppu.BackgroundFetcher.State = FetcherStateFetchTileNo
	ppu.BackgroundFetcher.WindowFetching = false
	ppu.BackgroundFetcher.WindowPixelRenderedThisScanline = false
	ppu.BackgroundFetcher.X = 0
	ppu.SpriteFetcher.Cycle = 0
	ppu.SpriteFetcher.State = FetcherStateFetchTileNo
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
	ppu.Shifter.Discard = ppu.RegSCX % 8
}

func (ppu *PPU) beginHBlank(ints *Interrupts) {
	ppu.setMode(ints, PPUModeHBlank)
	if ppu.PixelDrawCycle > 376 {
		ppu.HBlankRemainingCycles = 1
	} else {
		ppu.HBlankRemainingCycles = 376 - ppu.PixelDrawCycle
	}
}

func (ppu *PPU) ObjPalette(attribs Data8) Data8 {
	var palette Data8
	if attribs&Bit4 != 0 {
		palette = ppu.OBJPalette1
	} else {
		palette = ppu.OBJPalette0
	}
	palette &= 0xfc
	return palette
}

func (ppu *PPU) beginVBlank(ints *Interrupts) {
	ppu.BackgroundFetcher.WindowLineCounter = 0

	ppu.setMode(ints, PPUModeVBlank)

	ppu.VBlankLineRemainingCycles = 456

	// TODO: do we ever clear the VBlank interrupt?
	ints.IRQSet(IntSourceVBlank)

	ppu.FrameCount++
}

func (ppu *PPU) fsmOAMScan(ints *Interrupts, mem []Data8) {
	cycle := ppu.OAMScanCycle
	ppu.OAMScanCycle++
	if ppu.OAMScanCycle == 80 {
		ppu.beginPixelDraw(ints)
	}

	// PPU checks the OAM entry every 2 cycles
	if cycle&1 == 0 {
		return
	}
	index := Addr((cycle - 1) / 2)

	// Read sprite out of OAM
	sprite := DecodeObject(mem[AddrOAMBegin+index*4 : AddrOAMBegin+(index+1)*4])

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

func (ppu *PPU) fsmPixelDraw(ints *Interrupts, debug *Debug, clk *ClockRT, mem []Data8) {
	if ppu.SpriteFetcher.Suspended && ppu.SpriteFetcher.DoneX != ppu.Shifter.X {
		for idx := range ppu.OAMBuffer.Level {
			obj := ppu.OAMBuffer.Buffer[idx]
			if obj.X <= ppu.Shifter.X+8 && obj.X > ppu.Shifter.X {
				// Initiate sprite fetch
				ppu.SpriteFetcher.State = FetcherStateFetchTileNo
				ppu.SpriteFetcher.SpriteIDX = idx
				ppu.Shifter.Suspended = true
				ppu.SpriteFetcher.Suspended = false
				ppu.SpriteFetcher.DoneX = 0xff
				break
			}
		}
	} else {
		if ppu.SpriteFetcher.DoneX != 0xff {
			ppu.BackgroundFetcher.Suspended = false
			ppu.Shifter.Suspended = false
		}
	}

	ppu.SpriteFetcher.fsm(ppu, mem)
	ppu.BackgroundFetcher.fsm(ppu, mem)
	ppu.Shifter.fsm(ppu, debug, clk)

	// GBEDG: After each pixel shifted out, the PPU checks if it has reached the window
	if !ppu.BackgroundFetcher.WindowFetching && ppu.BackgroundFetcher.windowReached(ppu) {
		ppu.BackgroundFetcher.WindowFetching = true
		ppu.BackgroundFIFO.Clear()
		ppu.BackgroundFetcher.X = 0
	}

	if ppu.Shifter.X >= 160 {
		ppu.beginHBlank(ints)
		ppu.Shifter.X = 0
	}

	ppu.PixelDrawCycle++
}

func (ppu *PPU) fsmVBlank(ints *Interrupts, debug *Debug, fs *FrameSync) {
	if ppu.VBlankLineRemainingCycles > 0 {
		ppu.VBlankLineRemainingCycles--
		return
	}

	nSyncers := len(fs.ch)
	for range nSyncers {
		f := <-fs.ch
		f(&ppu.FBViewport)
	}
	ppu.IncRegLY(ints, debug)

	if ppu.RegLY == 0 {
		ppu.beginFrame(ints)
	}
}

func (ppu *PPU) fsmHBlank(ints *Interrupts, debug *Debug) {
	if ppu.HBlankRemainingCycles > 0 {
		ppu.HBlankRemainingCycles--
		return
	}

	ppu.IncRegLY(ints, debug)
	if ppu.BackgroundFetcher.WindowPixelRenderedThisScanline {
		ppu.BackgroundFetcher.WindowLineCounter++
	}

	if ppu.RegLY < 144 {
		ppu.beginOAMScan(ints)
	} else if ppu.RegLY == 144 {
		ppu.beginVBlank(ints)
	} else {
		panicv(ppu.RegLY)
	}
}

func (ppu *PPU) IncRegLY(ints *Interrupts, debug *Debug) {
	ppu.RegLY++
	if ppu.RegLY >= 153 {
		ppu.RegLY = 0
	}
	ppu.Stat.SetLYCEqLY(ints, ppu.RegLY == ppu.RegLYC)
	debug.SetY(ppu.RegLY)
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
	case AddrDMA:
		return ppu.DMA.Reg
	}
	return 0
}

func (ppu *PPU) Write(addr Addr, v Data8, ints *Interrupts) {
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
		ppu.SetLYC(ints, v)
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
	case AddrDMA:
		ppu.DMA.Write(v)
	default:
		panicf("Write to unknown LCD register %#v", addr)
	}
}
