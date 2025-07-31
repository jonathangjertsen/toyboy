package model

type Joypad struct {
	clk        *ClockRT
	Interrupts *Interrupts
	Action     Data8
	Direction  Data8
	mem        []Data8
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

func NewJoypad(clock *ClockRT, ints *Interrupts, mem []Data8) *Joypad {
	jp := &Joypad{
		clk:        clock,
		Interrupts: ints,
		Action:     0xf,
		Direction:  0xf,
		mem:        mem,
	}
	mem[AddrP1] = 0x1f
	return jp
}

func (jp *Joypad) Write(addr Addr, v Data8) {
	// TODO: this can trigger an interrupt
}

func (jp *Joypad) Read(addr Addr) Data8 {
	written := jp.mem[AddrP1]
	out := Data8(0x0f)
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
		actionMask := Data8(0b0000)
		directionMask := Data8(0b0000)
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
		if jp.mem[AddrP1]&0x20 == 0 {
			doJoypadInterrupt = (jp.Action & ^newAction) != 0
		} else {
			doJoypadInterrupt = (jp.Direction & ^newDirection) != 0
		}

		jp.Action = newAction
		jp.Direction = newDirection

		if doJoypadInterrupt {
			jp.Interrupts.IRQSet(IntSourceJoypad)
		}
	})
}
