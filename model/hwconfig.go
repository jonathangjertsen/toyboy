package model

type HWConfig struct {
	Model       Model
	SystemClock ClockConfig
}

type Model int

const (
	DMG Model = iota
)

type ClockConfig struct {
	Frequency float64
}
