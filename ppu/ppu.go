package ppu

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/types"
)

const (
	ioLoAddr = 0x2000
	ioHiAddr = 0x3FFF

	paletteRAMLoAddr = 0x3F00
	paletteRAMHiAddr = 0x3FFF

	oamMemorySize = 256

	visibleScanlineMax = 240
	visibleDotsMax     = 256
)

type SpriteSize byte

const (
	Sprite8x8 SpriteSize = iota
	Sprite8x16
)

const (
	PPUCTRL = ioLoAddr + iota
	PPUMASK
	PPUSTATUS
	OAMADDR
	OAMDATA
	PPUSCROLL
	PPUADDR
	PPUDATA
)

type PPU struct {
	ppuCtrl byte
	ppuMask byte

	spiteOverflow bool
	spite0Hit     bool
	vblankFlag    bool

	oamAddr byte

	ppuData byte

	currentVRAMAddr types.Word
	tempVRAMAddr    types.Word
	fineX           byte
	isFirstWrite    bool

	line  uint16
	pixel uint16

	chrMemoryBus *bus.PPUBus
	oamMemory    [oamMemorySize]byte
}

func (p *PPU) readOAMMemory(addr byte) byte {
	return p.oamMemory[addr]
}

func (p *PPU) writeOAMMemory(addr, value byte) {
	p.oamMemory[addr] = value
}

func (p *PPU) readPRGMemory(addr types.Word) byte {
	return p.chrMemoryBus.ReadPRGMemory(addr)
}

func (p *PPU) writePRGMemory(addr types.Word, value byte) {
	p.chrMemoryBus.WritePRGMemory(addr, value)
}

func (p *PPU) isCurrentlyRendering() bool {
	isVisibleScanline := p.line < visibleScanlineMax
	isVisibleDot := p.pixel < visibleDotsMax

	return p.isRenderingEnabled() && isVisibleDot && isVisibleScanline
}

func (p *PPU) isRenderingEnabled() bool {
	return p.isBgRenderingEnabled() && p.isSpriteRenderingEnabled()
}
