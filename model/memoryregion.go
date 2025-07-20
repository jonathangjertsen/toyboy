package model

const (
	ReadCountdown  = 100
	WriteCountdown = 100
)

type MemoryRegion struct {
	Offset Addr
	Data   []Data8

	WriteCountdowns  []uint64
	ReadCountdowns   []uint64
	CountdownDisable bool
}

func (mr *MemoryRegion) GetCounters(addr Addr) (uint64, uint64) {
	return mr.ReadCountdowns[addr-mr.Offset], mr.WriteCountdowns[addr-mr.Offset]
}

func (mr *MemoryRegion) Read(addr Addr) Data8 {
	idx := addr - mr.Offset
	if !mr.CountdownDisable { // todo move out of fast path
		mr.ReadCountdowns[idx] = ReadCountdown
	}
	return mr.Data[idx]
}

func (mr *MemoryRegion) Write(addr Addr, v Data8) {
	idx := addr - mr.Offset
	if !mr.CountdownDisable { // todo move out of fast path
		mr.WriteCountdowns[idx] = WriteCountdown
	}
	mr.Data[idx] = v
}

func (mr *MemoryRegion) DecrementCounters() {
	for i, c := range mr.WriteCountdowns {
		if c > 0 {
			mr.WriteCountdowns[i]--
		}
	}
	for i, c := range mr.ReadCountdowns {
		if c > 0 {
			mr.ReadCountdowns[i]--
		}
	}
}

func NewMemoryRegion(clock *ClockRT, start Addr, size Size16) MemoryRegion {
	mr := MemoryRegion{
		Data:            make([]Data8, size),
		WriteCountdowns: make([]uint64, size),
		ReadCountdowns:  make([]uint64, size),
		Offset:          start,
	}
	clock.AttachUIDevice(mr.DecrementCounters)
	return mr
}
