package ppu

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/types"
)

const (
	maxVRAMSize       = 2048
	maxPaletteRAMSize = 32

	vramLoAddr = 0x2000
	vramHiAddr = 0x2FFF

	paletteRAMLoAddr = 0x3F00
	paletteRAMHiAddr = 0x3FFF
)

type VRAM struct {
	mirroring nametable.NametableMirroring
	data      [maxVRAMSize]byte
}

func (v *VRAM) ReadCHRMemory(addr types.Word) (bool, byte) {
	if vramLoAddr <= addr && addr <= vramHiAddr {
		idx := v.getNTAddrWithMirroring(addr)
		return true, v.data[idx]
	}

	return false, 0
}

func (v *VRAM) WriteCHRMemory(addr types.Word, value byte) {
	if vramLoAddr <= addr && addr <= vramHiAddr {
		idx := v.getNTAddrWithMirroring(addr)
		v.data[idx] = value
	}
}

func (v *VRAM) getNTAddrWithMirroring(addr types.Word) types.Word {
	switch v.mirroring {
	case nametable.Vertical:
		return (addr - vramLoAddr) % verticalNametableMask
	case nametable.Horizontal:
		return (addr & 0x03FF) | ((addr & verticalNametableMask) >> 1)
	case nametable.SingleScreenLo:
		return (addr & 0x03FF)
	case nametable.SingleScreenHi:
		return (addr & 0x03FF) + constants.Kib1

	default:
		panic("Invalid Mirroing data")
	}
}

func (v *VRAM) UpdateMirroringCallback() func(nametable.NametableMirroring) {
	return func(m nametable.NametableMirroring) {
		v.mirroring = m
	}
}

type PaletteRAM struct {
	data [32]byte
}

func (p *PaletteRAM) ReadCHRMemory(addr types.Word) (bool, byte) {
	if paletteRAMLoAddr <= addr && addr <= paletteRAMHiAddr {
		idx := (addr - paletteRAMLoAddr) & 31
		if idx%4 == 0 {
			return true, p.data[0]
		}
		return true, p.data[idx]
	}
	return false, 0
}

func (p *PaletteRAM) WriteCHRMemory(addr types.Word, value byte) {
	if paletteRAMLoAddr <= addr && addr <= paletteRAMHiAddr {
		idx := (addr - paletteRAMLoAddr) & 31
		if idx%4 == 0 {
			p.data[idx&0x0F] = value
		} else {
			p.data[idx] = value
		}
	}
}
