package model

type LengthTimer struct {
	lengthEnable     bool
	lengthTimer      Data16
	lengthTimerReset Data8
}

func (lt *LengthTimer) SetResetValue(v Data8) {
	lt.lengthTimerReset = v
}

func (lt *LengthTimer) clock(expireValue Data16) bool {
	if !lt.lengthEnable {
		return false
	}
	if lt.lengthTimer < expireValue {
		lt.lengthTimer++
		return false
	}
	return lt.lengthTimer == expireValue
}
