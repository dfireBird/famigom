package palette

import _ "embed"

//go:embed palette.pal
var paletteData []byte

type Color struct {
	R byte
	G byte
	B byte
}

func GetColor(colorIdx byte) Color {
	dataIdx := colorIdx * 3
	R := paletteData[dataIdx]
	G := paletteData[dataIdx+1]
	B := paletteData[dataIdx+2]

	return Color{
		R: R,
		G: G,
		B: B,
	}
}
