package nametable

import "github.com/dfirebird/famigom/program"

type NametableMirroring int

const (
	Horizontal NametableMirroring = iota
	Vertical
	SingleScreenLo
	SingleScreenHi
)

func FromNametableArrangement(v program.NametableArrangement) NametableMirroring {
	if v == program.Vertical {
		return Horizontal
	} else {
		return Vertical
	}
}
