package model

type Rewind struct {
	Buffer []ExecLogEntry
	Idx    int
	Full   bool
}

func NewRewind(size int) *Rewind {
	return &Rewind{
		Buffer: make([]ExecLogEntry, size),
	}
}

func (rb *Rewind) Reset() {
	clear(rb.Buffer)
	rb.Idx = 0
	rb.Full = false
}

func (rb *Rewind) Start() int {
	if rb.Full {
		return rb.Next(rb.Idx)
	}
	return 0
}

func (rb *Rewind) End() int {
	return rb.Idx
}

func (rb *Rewind) Next(i int) int {
	i++
	if i == len(rb.Buffer) {
		i = 0
	}
	return i
}

func (rb *Rewind) At(i int) *ExecLogEntry {
	return &rb.Buffer[i]
}

func (rb *Rewind) Curr() *ExecLogEntry {
	idx := rb.Idx
	if idx > 0 {
		idx--
	} else {
		idx = len(rb.Buffer) - 1
	}
	return &rb.Buffer[idx]
}

func (rb *Rewind) Push() *ExecLogEntry {
	idx := rb.Idx
	rb.Idx++
	if rb.Idx >= len(rb.Buffer) {
		rb.Idx = 0
		rb.Full = true
	}
	return &rb.Buffer[idx]
}
