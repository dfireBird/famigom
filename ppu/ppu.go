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

	line     uint16
	pixel    uint16
	oddFrame bool

	chrMemoryBus *bus.PPUBus
	oamMemory    [oamMemorySize]byte
}

func CreatePPU() PPU {
	ppuBus := bus.CreatePPUBus()

	return PPU{
		chrMemoryBus: &ppuBus,
	}
}

func createPPUBus() bus.PPUBus {
	ppuBus := bus.PPUBus{
		devicesMap: []*bus.PPUBusDevice{},
	}

	vRAM := VRAM{
		data: [2048]byte{},
	}

	palleteRAM := PalleteRAM{
		data: [32]byte{},
	}

	ppuBus.RegisterDevice(&vRAM).RegisterDevice(&palleteRAM)

	return ppuBus
}

func (p *PPU) PowerUP() {
	p.ppuCtrl = 0
	p.ppuMask = 0
	p.vblankFlag = true
	p.spite0Hit = false
	p.spiteOverflow = true

	p.oamAddr = 0

	p.tempVRAMAddr = 0
	p.currentVRAMAddr = 0
	p.ppuData = 0

	p.oddFrame = false
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
