package ram

import (
	. "github.com/dfirebird/famigom/types"
)

const (
	maxRam = 1 << 11

	lowAddr=  0x0000
		highAddr= 0x1FFF
)

type RAM [maxRam]byte

func CreateRAM() *RAM {
	ram := RAM{}
	return &ram
}

func (r *RAM) ReadMemory(addr Word) (bool, byte) {
	if (lowAddr <= addr && addr <= highAddr) {
		return true, r[addr % maxRam]
	}

	return false, 0x00
}

func (r *RAM) WriteMemory(addr Word, value byte) {
	if (lowAddr <= addr && addr <= highAddr) {
		r[addr % maxRam] = value
	}
}
