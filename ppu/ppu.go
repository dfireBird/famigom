package ppu

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/program"
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

	preRenderScanLine  = 261
	visibleScanLineLo  = 0
	visibleScanLineHi  = 239
	postRenderScanLine = 240
	vblankScanLineLo   = 241
	vblankScanLineHi   = 260

	maxDots = 340

	horiPosCopyDot   = 257
	veriPosCopyDotLo = 280
	veriPosCopyDotHi = 304

	spriteDotLo = 257
	spirteDotHi = 320

	tilesForNextScanLineLo = 321
	tilesForNextScanLineHi = 336

	horizontalNametableMask = 0x0400
	verticalNametableMask   = 0x0800
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
	dot      uint16
	oddFrame bool

	chrMemoryBus *bus.PPUBus
	oamMemory    [oamMemorySize]byte

	nmiCallback *func()

	currentNTData byte
	currentATData byte
	currentTileLo byte
	currentTileHi byte
}

func CreatePPU(nmiCallback *func(), mirroring program.NametableArrangement) PPU {
	ppuBus := createPPUBus(mirroring)

	return PPU{
		nmiCallback: nmiCallback,
		chrMemoryBus: &ppuBus,
	}
}

func createPPUBus(mirroring program.NametableArrangement) bus.PPUBus {
	ppuBus := bus.NewPPUBus()

	vRAM := VRAM{
		mirroring: mirroring,
		data: [2048]byte{},
	}

	palleteRAM := PalleteRAM{
		data: [32]byte{},
	}

	ppuBus.RegisterDevice(&vRAM).RegisterDevice(&palleteRAM)

	return ppuBus
}

func (p *PPU) PowerUp() {
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

func (p *PPU) Step() {
	// FIXME: Do rendering of pixels

	if (visibleScanLineLo <= p.line && p.line <= visibleScanLineHi) || p.line == preRenderScanLine {
		if p.dot == 0 {
			p.incrementDot()
			return
		}

		if p.line == preRenderScanLine && p.dot == 1 {
			p.vblankFlag = false
			p.spite0Hit = false
			p.spiteOverflow = false
		}

		if 0 < p.dot && p.dot <= visibleDotsMax {
			p.doFetch(p.getBackgroundPatternTableAddr())
		}

		if spriteDotLo <= p.dot && p.dot <= spirteDotHi {
			// FIXME: sprites
		}

		if tilesForNextScanLineLo <= p.dot && p.dot <= tilesForNextScanLineHi {
			p.doFetch(p.getBackgroundPatternTableAddr())
		}

		if p.isRenderingEnabled() && p.dot == horiPosCopyDot {
			p.currentVRAMAddr &^= coarseXScrollMask                   // coarse X = 0
			p.currentVRAMAddr |= (p.tempVRAMAddr & coarseXScrollMask) // copy coarse X

			p.currentVRAMAddr &^= horizontalNametableMask                   // clear horizontal Nametable bit
			p.currentVRAMAddr |= (p.tempVRAMAddr & horizontalNametableMask) // copy	nt bit
		}

		if p.isRenderingEnabled() && (veriPosCopyDotLo <= p.dot && p.dot <= veriPosCopyDotHi) {
			p.currentVRAMAddr &^= coarseYScrollMask                   // coarse Y = 0
			p.currentVRAMAddr |= (p.tempVRAMAddr & coarseYScrollMask) // copy coarse Y

			p.currentVRAMAddr &^= verticalNametableMask                   // clear vertical nt bit
			p.currentVRAMAddr |= (p.tempVRAMAddr & verticalNametableMask) // copy nt bit

			p.currentVRAMAddr &^= fineYScrollMask                   // fine Y = 0
			p.currentVRAMAddr |= (p.tempVRAMAddr & fineYScrollMask) // copy fine Y
		}
	} else if p.line == vblankScanLineLo {
		if p.dot == 1 {
			p.vblankFlag = true
			(*p.nmiCallback)()
		}
	}

	p.incrementDot()
}

func (p *PPU) doFetch(patternTableHalf types.Word) {
	const singleTileMaxDot = 8

	switch localizedDot := ((p.dot - 1) % singleTileMaxDot) + 1; localizedDot {
	case 2:
		ntAddr := 0x2000 | (p.currentVRAMAddr & 0x0FFF)
		p.currentNTData = p.readPRGMemory(ntAddr)
	case 4:
		atAddr := 0x23C0 | (p.currentVRAMAddr & 0x0C00) | ((p.currentVRAMAddr >> 4) & 0x38) | ((p.currentVRAMAddr >> 2) & 0x07)
		p.currentATData = p.readPRGMemory(atAddr)
	case 6:
		ptAddr := p.calcPatternTableAddr(0, patternTableHalf)
		p.currentTileLo = p.readPRGMemory(ptAddr)
	case 8:
		ptAddr := p.calcPatternTableAddr(8, patternTableHalf)
		p.currentTileHi = p.readPRGMemory(ptAddr)
		p.incrementX()
	}

	if p.dot == visibleDotsMax {
		p.incrementY()
	}
}

func (p *PPU) incrementDot() {
	newDot := (p.dot + 1) % (maxDots + 1)
	newLineVal := p.line

	if p.oddFrame && newDot == (maxDots-1) && newLineVal == preRenderScanLine {
		newDot = 0
		newLineVal = 0
	} else if newDot == 0 {
		newLineVal = (newLineVal + 1) % (preRenderScanLine + 1)
	}

	p.dot = newDot
	p.line = newLineVal
}

func (p *PPU) incrementX() {
	if p.currentVRAMAddr&coarseXScrollMask == coarseXScrollMask {
		p.currentVRAMAddr &^= coarseXScrollMask
		p.currentVRAMAddr ^= horizontalNametableMask
	} else {
		p.currentVRAMAddr += 1
	}
}

func (p *PPU) incrementY() {
	if p.currentVRAMAddr&fineYScrollMask != fineYScrollMask {
		p.currentVRAMAddr += 0x1000
	} else {
		p.currentVRAMAddr &^= fineYScrollMask // fine Y = 0
		coarseY := (p.currentVRAMAddr & coarseYScrollMask) >> coarseYScrollShift
		switch coarseY {
		case 29:
			coarseY = 0
			p.currentVRAMAddr ^= verticalNametableMask
		case 31:
			coarseY = 0 // nametable swtich is not allowed
		default:
			coarseY += 1
		}

		p.currentVRAMAddr = (p.currentVRAMAddr &^ coarseXScrollMask) | (coarseY << coarseYScrollShift)
	}
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
	isVisibleDot := p.dot < visibleDotsMax

	return p.isRenderingEnabled() && isVisibleDot && isVisibleScanline
}

func (p *PPU) isRenderingEnabled() bool {
	return p.isBgRenderingEnabled() && p.isSpriteRenderingEnabled()
}

func (p *PPU) calcPatternTableAddr(bitPlane, patternTableHalf types.Word) types.Word {
	fineY := (p.currentVRAMAddr & fineYScrollMask) >> fineYScrollShift
	tileNo := types.Word(p.currentNTData) << 4
	return patternTableHalf | tileNo | bitPlane | fineY
}
