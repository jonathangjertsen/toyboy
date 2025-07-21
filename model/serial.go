package model

type Serial struct {
	Mem MemoryRegion
}

func NewSerial(clk *ClockRT) *Serial {
	return &Serial{
		Mem: NewMemoryRegion(clk, AddrSB, 2),
	}
}

func (ser *Serial) Read(addr Addr) Data8 {
	return ser.Mem.Read(addr)
}

func (ser *Serial) Write(addr Addr, v Data8) {
	ser.Mem.Write(addr, v)
}

func (ser *Serial) GetCounters(addr Addr) (uint64, uint64) {
	return ser.Mem.GetCounters(addr)
}
