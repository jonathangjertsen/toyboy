package model

type Joypad struct {
	Action    Data8
	Direction Data8
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

func (jp *Joypad) Write(addr Addr, v Data8) {
	// TODO: this can trigger an interrupt
}

func (jp *Joypad) Read(written Data8, addr Addr) Data8 {
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

func (jp *Joypad) SetState(clk *ClockRT, gb *Gameboy, jps JoypadState) {
	clk.Sync(func() {
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
		if gb.Mem[AddrP1]&0x20 == 0 {
			doJoypadInterrupt = (jp.Action & ^newAction) != 0
		} else {
			doJoypadInterrupt = (jp.Direction & ^newDirection) != 0
		}

		jp.Action = newAction
		jp.Direction = newDirection

		if doJoypadInterrupt {
			gb.Interrupts.IRQSet(gb, IntSourceJoypad)
		}
	})
}
