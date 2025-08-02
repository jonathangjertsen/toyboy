package model

import (
	"time"
)

type AudioSample = int16

type Audio interface {
	Clock(*APU)
	SetMPeriod(time.Duration)
}
