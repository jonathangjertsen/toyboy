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
	Ch chan func(*ViewPort)
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
	fs.Ch <- func(vp *ViewPort) {
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

func (ppu *PPU) SetLYC(gb *Gameboy, v Data8) {
	ppu.RegLYC = v
	gb.IRQCheck()
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

func (ppu *PPU) fsm(gb *Gameboy, clk *ClockRT, fs *FrameSync) {
	if ppu.DMA.Source != 0 {
		gb.fsmPPUDMA()
	}
	switch ppu.Mode {
	case PPUModeVBlank:
		ppu.fsmVBlank(gb, fs)
	case PPUModeHBlank:
		ppu.fsmHBlank(gb)
	case PPUModePixelDraw:
		gb.fsmPixelDraw(clk)
	case PPUModeOAMScan:
		ppu.fsmOAMScan(gb)
	}
}

func (ppu *PPU) setMode(gb *Gameboy, mode PPUMode) {
	ppu.Mode = mode
	ppu.Stat.SetMode(gb, mode)
}

func (gb *Gameboy) beginFrame() {
	gb.beginOAMScan()
	gb.PPU.BackgroundFetcher.WindowYReached = false
}

func (gb *Gameboy) beginOAMScan() {
	gb.PPU.setMode(gb, PPUModeOAMScan)
	gb.PPU.OAMScanCycle = 0
	gb.PPU.OAMBuffer.Level = 0
}

// start of scanline after 7OAM scan
func (ppu *PPU) beginPixelDraw(gb *Gameboy) {
	ppu.setMode(gb, PPUModePixelDraw)
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

func (ppu *PPU) beginHBlank(gb *Gameboy) {
	ppu.setMode(gb, PPUModeHBlank)
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

func (ppu *PPU) beginVBlank(gb *Gameboy) {
	ppu.BackgroundFetcher.WindowLineCounter = 0

	ppu.setMode(gb, PPUModeVBlank)

	ppu.VBlankLineRemainingCycles = 456

	// TODO: do we ever clear the VBlank interrupt?
	gb.IRQSet(IntSourceVBlank)

	ppu.FrameCount++
}

func (ppu *PPU) fsmOAMScan(gb *Gameboy) {
	cycle := ppu.OAMScanCycle
	ppu.OAMScanCycle++
	if ppu.OAMScanCycle == 80 {
		ppu.beginPixelDraw(gb)
	}

	// PPU checks the OAM entry every 2 cycles
	if cycle&1 == 0 {
		return
	}
	index := Addr((cycle - 1) / 2)

	// Read sprite out of OAM
	sprite := DecodeObject(gb.Mem[AddrOAMBegin+index*4 : AddrOAMBegin+(index+1)*4])

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

func (gb *Gameboy) fsmPixelDraw(clk *ClockRT) {
	gb.checkSpriteHit()
	gb.spriteFetcherFSM()
	gb.backgroundFetcherFSM()
	gb.shifterFSM(clk)
	gb.checkWindowReached()
	gb.checkHBlankReached()
	gb.PPU.PixelDrawCycle++
}

func (gb *Gameboy) checkSpriteHit() {
	if gb.PPU.SpriteFetcher.Suspended && gb.PPU.SpriteFetcher.DoneX != gb.PPU.Shifter.X {
		for idx := range gb.PPU.OAMBuffer.Level {
			obj := gb.PPU.OAMBuffer.Buffer[idx]
			if obj.X <= gb.PPU.Shifter.X+8 && obj.X > gb.PPU.Shifter.X {
				// Initiate sprite fetch
				gb.PPU.SpriteFetcher.State = FetcherStateFetchTileNo
				gb.PPU.SpriteFetcher.SpriteIDX = idx
				gb.PPU.Shifter.Suspended = true
				gb.PPU.SpriteFetcher.Suspended = false
				gb.PPU.SpriteFetcher.DoneX = 0xff
				break
			}
		}
	} else {
		if gb.PPU.SpriteFetcher.DoneX != 0xff {
			gb.PPU.BackgroundFetcher.Suspended = false
			gb.PPU.Shifter.Suspended = false
		}
	}
}

func (gb *Gameboy) checkWindowReached() {
	if !gb.PPU.BackgroundFetcher.WindowFetching && gb.PPU.BackgroundFetcher.windowReached(gb) {
		gb.PPU.BackgroundFetcher.WindowFetching = true
		gb.PPU.BackgroundFIFO.Clear()
		gb.PPU.BackgroundFetcher.X = 0
	}
}

func (gb *Gameboy) checkHBlankReached() {
	if gb.PPU.Shifter.X >= 160 {
		gb.PPU.beginHBlank(gb)
		gb.PPU.Shifter.X = 0
	}
}

func (ppu *PPU) fsmVBlank(gb *Gameboy, fs *FrameSync) {
	if ppu.VBlankLineRemainingCycles > 0 {
		ppu.VBlankLineRemainingCycles--
		return
	}

	nSyncers := len(fs.Ch)
	for range nSyncers {
		f := <-fs.Ch
		f(&ppu.FBViewport)
	}
	ppu.IncRegLY(gb)

	if ppu.RegLY == 0 {
		gb.beginFrame()
	}
}

func (ppu *PPU) fsmHBlank(gb *Gameboy) {
	if ppu.HBlankRemainingCycles > 0 {
		ppu.HBlankRemainingCycles--
		return
	}

	ppu.IncRegLY(gb)
	if ppu.BackgroundFetcher.WindowPixelRenderedThisScanline {
		ppu.BackgroundFetcher.WindowLineCounter++
	}

	if ppu.RegLY < 144 {
		gb.beginOAMScan()
	} else if ppu.RegLY == 144 {
		ppu.beginVBlank(gb)
	} else {
		panicv(ppu.RegLY)
	}
}

func (ppu *PPU) IncRegLY(gb *Gameboy) {
	ppu.RegLY++
	if ppu.RegLY >= 153 {
		ppu.RegLY = 0
	}
	ppu.Stat.SetLYCEqLY(gb, ppu.RegLY == ppu.RegLYC)
	gb.Debug.SetY(ppu.RegLY)
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

func (gb *Gameboy) WritePPU(addr Addr, v Data8) {
	ppu := &gb.PPU

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
		ppu.SetLYC(gb, v)
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
