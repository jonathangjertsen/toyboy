package model

type BootROMLock struct {
	MemoryRegion

	BootOff bool
	OnLock  func()
}

func NewBootROMLock(clock *ClockRT) *BootROMLock {
	return &BootROMLock{
		MemoryRegion: NewMemoryRegion(clock, AddrBootROMLock, 0x0001),
	}
}

func (brl *BootROMLock) Read(addr Addr) Data8 {
	return brl.MemoryRegion.Read(addr)
}

func (brl *BootROMLock) Write(addr Addr, v Data8) {
	if brl.BootOff {
		return
	}
	brl.MemoryRegion.Write(addr, v)
	if v&1 == 1 {
		brl.Lock()
	}
}

func (brl *BootROMLock) Lock() {
	brl.BootOff = true
	if brl.OnLock != nil {
		brl.OnLock()
	}
}
