package model

type LengthTimer struct {
	Enable  bool
	Counter Data16
	Reset   Data8
}

func (lt *LengthTimer) SetResetValue(v Data8) {
	lt.Reset = v
}

func (lt *LengthTimer) clock(expireValue Data16) bool {
	if !lt.Enable {
		return false
	}
	if lt.Counter < expireValue {
		lt.Counter++
		return false
	}
	return lt.Counter == expireValue
}
