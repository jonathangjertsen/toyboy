package model

// Also known as a "Sprite".
type Object struct {
	X          Data8
	Y          Data8
	TileIndex  Data8
	Attributes Data8
}

func DecodeObject(data []Data8) Object {
	return Object{
		Y:          data[0],
		X:          data[1],
		TileIndex:  data[2],
		Attributes: data[3],
	}
}
