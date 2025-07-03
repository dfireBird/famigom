package ram

import (
	. "github.com/dfirebird/famigom/types"
)

const maxRam = 1 << 11

type RAM [maxRam]byte

func CreateRAM() (*RAM, AddrRange) {
	ram := RAM{}
	addrRange := AddrRange {
		LowAddr:  0x0000,
		HighAddr: 0x1FFF,
	}

	return &ram, addrRange
}

func (r *RAM) ReadMemory(addr Word) byte {
	return r[addr % maxRam]
}

func (r *RAM) WriteMemory(addr Word, value byte) {
    r[addr % maxRam] = value
}
