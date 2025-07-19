package model

type Joypad struct {
	clk        *ClockRT
	Interrupts *Interrupts
	Written    MemoryRegion
	Action     uint8
	Direction  uint8
}

type JoypadState struct {
	Up     bool
	Left   bool
	Right  bool
	Down   bool
	A      bool
	B      bool
	Start  bool
	Select bool
}

func NewJoypad(clock *ClockRT, ints *Interrupts) *Joypad {
	jp := &Joypad{
		clk:        clock,
		Interrupts: ints,
		Written:    NewMemoryRegion(clock, 0xff00, 0x0001),
		Action:     0xf,
		Direction:  0xf,
	}
	jp.Written.Data[0] = 0x1f
	return jp
}

func (jp *Joypad) Write(addr uint16, v uint8) {
	// TODO: can this trigger an interrupt?
	jp.Written.Write(addr, v)
}

func (jp *Joypad) Read(addr uint16) uint8 {
	written := jp.Written.Read(addr)
	out := uint8(0x0f)
	if written&0x20 == 0 {
		out &= jp.Action
	}
	if written&0x10 == 0 {
		out &= jp.Direction
	}
	out |= (written & 0xf0)
	return out
}

func (jp *Joypad) SetState(jps JoypadState) {
	// TODO:
	jp.clk.Sync(func() {
		actionMask := uint8(0b0000)
		directionMask := uint8(0b0000)
		if jps.A {
			actionMask |= 0b0001
		}
		if jps.B {
			actionMask |= 0b0010
		}
		if jps.Select {
			actionMask |= 0b0100
		}
		if jps.Start {
			actionMask |= 0b1000
		}
		if jps.Right {
			directionMask |= 0b0001
		}
		if jps.Left {
			directionMask |= 0b0010
		}
		if jps.Up {
			directionMask |= 0b0100
		}
		if jps.Down {
			directionMask |= 0b1000
		}

		newAction := 0xf ^ actionMask
		newDirection := 0xf ^ directionMask

		doJoypadInterrupt := false
		if jp.Written.Data[0]&0x20 == 0 {
			doJoypadInterrupt = (jp.Action & ^newAction) != 0
		} else {
			doJoypadInterrupt = (jp.Direction & ^newDirection) != 0
		}

		jp.Action = newAction
		jp.Direction = newDirection

		if doJoypadInterrupt {
			jp.Interrupts.IRQSet(0x10)
		}
	})
}
