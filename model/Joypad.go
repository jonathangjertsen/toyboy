package model

type Joypad struct {
	Written   MemoryRegion
	Action    uint8
	Direction uint8
	Select    uint8
}

func NewJoypad(clock *ClockRT) *Joypad {
	return &Joypad{
		Written: NewMemoryRegion(clock, 0xff00, 0x0001),
	}
}

func (jp *Joypad) Write(addr uint16, v uint8) {
	jp.Written.Write(addr, v)
}

func (jp *Joypad) Read(addr uint16) uint8 {
	written := jp.Written.Read(addr)
	out := uint8(0x0f)
	out |= (written & 0xf0)
	if written&0x20 == 0 {
		out &= jp.Action
	}
	if written&0x10 == 0 {
		out &= jp.Direction
	}
	return out
}
