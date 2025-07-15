package plugin

import (
	"sync/atomic"
	"time"
)

type ClockMeasurement struct {
	counter atomic.Uint64
	t0      time.Time
}

func NewClockMeasurement() *ClockMeasurement {
	return &ClockMeasurement{}
}

func (cm *ClockMeasurement) Start() {
	cm.counter.Store(0)
	cm.t0 = time.Now()
}

func (cm *ClockMeasurement) Stop() (uint64, time.Duration) {
	return cm.counter.Load(), time.Since(cm.t0)
}

func (cm *ClockMeasurement) Clocked() {
	cm.counter.Add(1)
}
