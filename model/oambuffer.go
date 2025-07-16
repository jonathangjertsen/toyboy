package model

type OAMBuffer struct {
	Buffer [10]Sprite
	Level  int
}

func (oamBuffer *OAMBuffer) Add(sprite Sprite) {
	oamBuffer.Buffer[oamBuffer.Level] = sprite
	oamBuffer.Level++
}

func (oamBuffer *OAMBuffer) Full() bool {
	return oamBuffer.Level >= 10
}
