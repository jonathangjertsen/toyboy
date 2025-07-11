package model

type HWConfig struct {
	SystemClock ClockConfig
}

type ClockConfig struct {
	Frequency float64
}
