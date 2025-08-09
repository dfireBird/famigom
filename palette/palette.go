package palette

import _ "embed"

//go:embed palette.pal
var palleteData []byte

type Color struct {
	R byte
	G byte
	B byte
}

func GetColor(colorIdx byte) Color {
	dataIdx := colorIdx * 3
	R := palleteData[dataIdx]
	G := palleteData[dataIdx+1]
	B := palleteData[dataIdx+2]

	return Color{
		R: R,
		G: G,
		B: B,
	}
}
