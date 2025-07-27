package model

const (
	ReadCountdown  = 100
	WriteCountdown = 100
)

type MemoryRegion struct {
	Offset Addr
	Data   []Data8
}

func (mr *MemoryRegion) Read(addr Addr) Data8 {
	idx := addr - mr.Offset
	return mr.Data[idx]
}

func (mr *MemoryRegion) Write(addr Addr, v Data8) {
	idx := addr - mr.Offset
	mr.Data[idx] = v
}

func NewMemoryRegion(clock *ClockRT, start Addr, size Size16) MemoryRegion {
	mr := MemoryRegion{
		Data:   make([]Data8, size),
		Offset: start,
	}
	return mr
}
