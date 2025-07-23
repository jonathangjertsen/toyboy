package model

type OAMBuffer struct {
	Buffer [10]Object
	Level  int
}

func (oamBuffer *OAMBuffer) Add(sprite Object) {
	oamBuffer.Buffer[oamBuffer.Level] = sprite
	oamBuffer.Level++
}

func (oamBuffer *OAMBuffer) Full() bool {
	return oamBuffer.Level >= 10
}
