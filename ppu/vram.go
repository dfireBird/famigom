package ppu

import (
	"github.com/dfirebird/famigom/program"
	"github.com/dfirebird/famigom/types"
)

const (
	maxVRAMSize       = 2048
	maxPalleteRAMSize = 32

	vramLoAddr = 0x2000
	vramHiAddr = 0x2FFF

	palleteRAMLoAddr = 0x3F00
	palleteRAMHiAddr = 0x3FFF

	palleteRAMIndexMask = 0x3F1F
)

type VRAM struct {
	mirroring program.NametableArrangement
	data      [maxVRAMSize]byte
}

func (v *VRAM) ReadPRGMemory(addr types.Word) (bool, byte) {
	if vramLoAddr <= addr && addr <= vramHiAddr {
		idx := v.getNTAddrWithMirroring(addr)
		return true, v.data[idx]
	}

	return false, 0
}

func (v *VRAM) WritePRGMemory(addr types.Word, value byte) {
	if vramLoAddr <= addr && addr <= vramHiAddr {
		idx := v.getNTAddrWithMirroring(addr)
		v.data[idx] = value
	}
}

func (v *VRAM) getNTAddrWithMirroring(addr types.Word) types.Word {
	switch v.mirroring {
	case program.Vertical:
		return (addr - vramLoAddr) % verticalNametableMask
	case program.Horizontal:
		return (addr & 0x03FF) | ((addr & verticalNametableMask) >> 1)

	default:
		panic("Invalid Mirroing data")
	}
}

type PalleteRAM struct {
	data [32]byte
}

func (p *PalleteRAM) ReadPRGMemory(addr types.Word) (bool, byte) {
	if palleteRAMLoAddr <= addr && addr <= palleteRAMHiAddr {
		idx := (addr & palleteRAMIndexMask) - palleteRAMLoAddr
		return true, p.data[idx]
	}
	return false, 0
}

func (p *PalleteRAM) WritePRGMemory(addr types.Word, value byte) {
	if palleteRAMLoAddr <= addr && addr <= palleteRAMHiAddr {
		idx := (addr & palleteRAMIndexMask) - palleteRAMLoAddr
		p.data[idx] = value
	}
}
