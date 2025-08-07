package ppu

import "github.com/dfirebird/famigom/types"

const (
	maxVRAMSize       = 2048
	maxPalleteRAMSize = 32

	vramLoAddr = 0x2000
	vramHiAddr = 0x2FFF

	palleteRamLoAddr = 0x3F00
	palleteRamHiAddr = 0x3FFF

	palletRamIndexMask = 0x3F1F
)

type VRAM struct {
	data [maxVRAMSize]byte
}

func (v *VRAM) ReadPRGMemory(addr types.Word) (bool, byte) {
	if vramLoAddr <= addr && addr <= vramHiAddr {
		idx := addr - vramLoAddr
		return true, v.data[idx]
	}

	return false, 0
}

func (v *VRAM) WritePRGMemory(addr types.Word, value byte) {
	if vramLoAddr <= addr && addr <= vramHiAddr {
		idx := addr - vramLoAddr
		v.data[idx] = value
	}
}

type PalleteRAM struct {
	data [32]byte
}

func (p *PalleteRAM) ReadPRGMemory(addr types.Word) (bool, byte) {
	if palleteRamLoAddr <= addr && addr <= palleteRamHiAddr {
		idx := (addr & palletRamIndexMask) - palleteRamLoAddr
		return true, p.data[idx]
	}
	return false, 0
}

func (p *PalleteRAM) WritePRGMemory(addr types.Word, value byte) {
	if palleteRamLoAddr <= addr && addr <= palleteRamHiAddr {
		idx := (addr & palletRamIndexMask) - palleteRamLoAddr
		p.data[idx] = value
	}
}
