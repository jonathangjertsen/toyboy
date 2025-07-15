package model

type MemoryRegion struct {
	Offset uint16
	Data   []uint8
}

func (mr *MemoryRegion) Read(addr uint16) uint8 {
	return mr.Data[addr-mr.Offset]
}

func (mr *MemoryRegion) Write(addr uint16, v uint8) {
	mr.Data[addr-mr.Offset] = v
}

func NewMemoryRegion(start, size uint16) MemoryRegion {
	return MemoryRegion{
		Data:   make([]uint8, size),
		Offset: start,
	}
}
