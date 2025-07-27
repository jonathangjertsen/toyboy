package plugin

import (
	"time"
)

type ClockMeasurement struct {
	initCount uint
	counter   *uint
	t0        time.Time
}

func NewClockMeasurement(counter *uint) *ClockMeasurement {
	return &ClockMeasurement{
		counter: counter,
	}
}

func (cm *ClockMeasurement) Start() {
	cm.initCount = *cm.counter
	cm.t0 = time.Now()
}

func (cm *ClockMeasurement) Stop() (uint, time.Duration) {
	return *cm.counter - cm.initCount, time.Since(cm.t0)
}
