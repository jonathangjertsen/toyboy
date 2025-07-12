package model

type MemoryRegion struct {
	name   string
	offset uint16
	data   []uint8
}

func (mr *MemoryRegion) Name() string {
	return mr.name
}

func (mr *MemoryRegion) Range() (uint16, uint16) {
	return mr.offset, uint16(len(mr.data))
}

func (mr *MemoryRegion) Read(addr uint16) uint8 {
	// fmt.Printf("read %s[%x]\n", mr.name, addr-mr.offset)
	return mr.data[addr-mr.offset]
}

func (mr *MemoryRegion) Write(addr uint16, v uint8) {
	// fmt.Printf("wrote %s[%x]=%x\n", mr.name, addr-mr.offset, v)
	mr.data[addr-mr.offset] = v
}

func NewMemoryRegion(name string, start, size uint16) *MemoryRegion {
	return &MemoryRegion{
		name:   name,
		data:   make([]uint8, size),
		offset: start,
	}
}
