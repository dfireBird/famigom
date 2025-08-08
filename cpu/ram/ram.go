package ram

import (
	"github.com/dfirebird/famigom/types"
)

const (
	maxRAM = 1 << 11

	lowAddr  = 0x0000
	highAddr = 0x1FFF
)

type RAM [maxRAM]byte

func CreateRAM() *RAM {
	ram := RAM{}
	return &ram
}

func (r *RAM) ReadMemory(addr types.Word) (bool, byte) {
	if lowAddr <= addr && addr <= highAddr {
		return true, r[addr%maxRAM]
	}

	return false, 0x00
}

func (r *RAM) WriteMemory(addr types.Word, value byte) {
	if lowAddr <= addr && addr <= highAddr {
		r[addr%maxRAM] = value
	}
}

type DMA struct {
	dmaCallback *func(byte)
}

func CreateDMADevice(dmaCallback *func(byte)) *DMA {
	return &DMA{
		dmaCallback: dmaCallback,
	}
}

func (d *DMA) ReadMemory(addr types.Word) (bool, byte) {
	return false, 0
}

func (d *DMA) WriteMemory(addr types.Word, value byte) {
	if addr == 0x4014 {
		(*d.dmaCallback)(value)
	}
}
