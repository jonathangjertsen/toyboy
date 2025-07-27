package plugin

import (
	"time"
)

type ClockMeasurement struct {
	initCount uint64
	counter   *uint64
	t0        time.Time
}

func NewClockMeasurement(counter *uint64) *ClockMeasurement {
	return &ClockMeasurement{
		counter: counter,
	}
}

func (cm *ClockMeasurement) Start() {
	cm.initCount = *cm.counter
	cm.t0 = time.Now()
}

func (cm *ClockMeasurement) Stop() (uint64, time.Duration) {
	return *cm.counter - cm.initCount, time.Since(cm.t0)
}
