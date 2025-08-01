package model

type DMA struct {
	Reg    Data8
	Source Addr
	Dest   Addr
}

func (d *DMA) Write(v Data8) {
	d.Reg = v
	d.Source = Addr(join16(v, 0x00))
	d.Dest = AddrOAMBegin
}

func (d *DMA) fsm(gb *Gameboy) {
	if d.Source == 0 {
		return
	}

	// Write next data
	// TODO: presumably this is not actually how it works
	gb.Mem[d.Dest] = gb.Mem[d.Source]

	if d.Dest == AddrOAMEnd {
		// Done
		d.Source = 0
		d.Dest = AddrOAMBegin
	} else {
		// Set next source and dest
		d.Source++
		d.Dest++
	}
}
