package model

type DMA struct {
	Reg    Data8
	Source Addr
	Dest   Addr

	mem []Data8
}

func (d *DMA) Write(v Data8) {
	d.Reg = v
	d.Source = Addr(join16(v, 0x00))
	d.Dest = AddrOAMBegin
}

func (d *DMA) fsm() {
	if d.Source == 0 {
		return
	}

	// Write next data
	// TODO: presumably this is not actually how it works
	d.mem[d.Dest] = d.mem[d.Source]

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
