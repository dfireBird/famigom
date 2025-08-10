package ppu

import (
	"github.com/dfirebird/famigom/program"
	"github.com/dfirebird/famigom/types"
)

const (
	maxVRAMSize       = 2048
	maxPaletteRAMSize = 32

	vramLoAddr = 0x2000
	vramHiAddr = 0x2FFF

	paletteRAMLoAddr = 0x3F00
	paletteRAMHiAddr = 0x3FFF

	paletteRAMIndexMask = 0x3F1F
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

type PaletteRAM struct {
	data [32]byte
}

func (p *PaletteRAM) ReadPRGMemory(addr types.Word) (bool, byte) {
	if paletteRAMLoAddr <= addr && addr <= paletteRAMHiAddr {
		idx := (addr & paletteRAMIndexMask) - paletteRAMLoAddr
		return true, p.data[idx]
	}
	return false, 0
}

func (p *PaletteRAM) WritePRGMemory(addr types.Word, value byte) {
	if paletteRAMLoAddr <= addr && addr <= paletteRAMHiAddr {
		idx := (addr & paletteRAMIndexMask) - paletteRAMLoAddr
		if idx == 0x00 {
			p.data[0x10] = value
		}
		if idx == 0x10 {
			p.data[0x00] = value
		}
		p.data[idx] = value
	}
}
