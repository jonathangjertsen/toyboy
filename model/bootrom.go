package model

type BootROMLock struct {
	BootOff bool
	OnLock  func()
}

func NewBootROMLock(clock *ClockRT) *BootROMLock {
	return &BootROMLock{}
}

func (brl *BootROMLock) Write(addr Addr, v Data8) {
	if brl.BootOff {
		return
	}
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
