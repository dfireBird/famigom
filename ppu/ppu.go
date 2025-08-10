package ppu

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/program"
	"github.com/dfirebird/famigom/types"
)

const (
	ioLoAddr = 0x2000
	ioHiAddr = 0x3FFF

	oamMemorySize = 256

	visibleScanlineMax = 240
	visibleDotsMax     = 256
	totalDots          = visibleScanlineMax * visibleDotsMax

	preRenderScanLine  = 261
	visibleScanLineLo  = 0
	visibleScanLineHi  = 239
	postRenderScanLine = 240
	vblankScanLineLo   = 241
	vblankScanLineHi   = 260

	nextScanline1stTileDot = 328
	nextScanline2ndTileDot = 336

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

	spritePaletteMSB = 1
	bgPaletteMSB     = 0
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

	nTDataLatch        byte
	aTDataLatch        byte
	patternDataLoLatch byte
	patternDataHiLatch byte

	nextTilePatternShiftRegisterLo byte
	nextTilePatternShiftRegisterHi byte
	nextTileAttributShiftRegister  byte

	currentTilePatternShiftRegisterLo byte
	currentTilePatternShiftRegisterHi byte
	currentTileAttributShiftRegister  byte

	VirtualDisplay [totalDots]byte
}

func CreatePPU(nmiCallback *func(), mirroring program.NametableArrangement) PPU {
	ppuBus := createPPUBus(mirroring)

	return PPU{
		nmiCallback:  nmiCallback,
		chrMemoryBus: &ppuBus,
	}
}

func createPPUBus(mirroring program.NametableArrangement) bus.PPUBus {
	ppuBus := bus.NewPPUBus()

	vRAM := VRAM{
		mirroring: mirroring,
		data:      [2048]byte{},
	}

	paletteRAM := PaletteRAM{
		data: [32]byte{},
	}

	ppuBus.RegisterDevice(&vRAM).RegisterDevice(&paletteRAM)

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

func (p *PPU) RegisterDevice(d bus.PPUBusDevice) {
	p.chrMemoryBus.RegisterDevice(d)
}

func (p *PPU) Step() {
	// FIXME: Do rendering of pixels
	// fmt.Printf("(%d, %d):    0x%x\n", p.line, p.dot, p.currentVRAMAddr)

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
			p.outputPixel()

			ld := p.doFetch(p.getBackgroundPatternTableAddr())
			if ld == 8 {
				p.currentTilePatternShiftRegisterLo = p.nextTilePatternShiftRegisterLo
				p.currentTilePatternShiftRegisterHi = p.nextTilePatternShiftRegisterHi
				p.currentTileAttributShiftRegister = p.nextTileAttributShiftRegister

				p.nextTilePatternShiftRegisterLo = p.patternDataLoLatch
				p.nextTilePatternShiftRegisterHi = p.patternDataHiLatch
				p.nextTileAttributShiftRegister = p.aTDataLatch
			}
		}

		if spriteDotLo <= p.dot && p.dot <= spirteDotHi {
			// FIXME: sprites
		}

		if tilesForNextScanLineLo <= p.dot && p.dot <= tilesForNextScanLineHi {
			p.doFetch(p.getBackgroundPatternTableAddr())

			if p.dot == nextScanline1stTileDot {
				p.currentTilePatternShiftRegisterLo = p.patternDataLoLatch
				p.currentTilePatternShiftRegisterHi = p.patternDataHiLatch
				p.currentTileAttributShiftRegister = p.aTDataLatch
			}

			if p.dot == nextScanline2ndTileDot {
				p.nextTilePatternShiftRegisterLo = p.patternDataLoLatch
				p.nextTilePatternShiftRegisterHi = p.patternDataHiLatch
				p.nextTileAttributShiftRegister = p.aTDataLatch
			}
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

func (p *PPU) GetBackdropColorIdx() byte {
	var colorIdx byte
	if paletteRAMLoAddr <= p.currentVRAMAddr && p.currentVRAMAddr <= paletteRAMHiAddr {
		colorIdx = p.readPRGMemory(p.currentVRAMAddr)
	} else {
		colorIdx = p.readPRGMemory(paletteRAMLoAddr)
	}

	return colorIdx
}

func (p *PPU) outputPixel() {
	paletteRAMIDx := byte(0)
	displayIdx := (p.line*visibleDotsMax + p.dot) - 1

	if p.isBgRenderingEnabled() {
		fineXBitSelect := byte(1) >> p.fineX
		tileLo := (p.currentTilePatternShiftRegisterLo & fineXBitSelect) >> p.fineX
		tileHi := (p.currentTilePatternShiftRegisterHi & fineXBitSelect) >> p.fineX

		coarseXBit1 := (p.currentVRAMAddr & coarseXScroll1Bit) >> coarseXScroll1Shift
		coarseYBit1 := (p.currentVRAMAddr & coarseYScroll1Bit) >> coarseYScroll1Shift
		atLoPos := coarseXBit1*2 + coarseYBit1*4
		atHiPos := atLoPos + 1

		attrLo := p.currentTileAttributShiftRegister & (1 << atLoPos)
		attrHi := p.currentTileAttributShiftRegister & (1 << atHiPos)

		paletteRAMIDx = bgPaletteMSB<<4 | attrHi<<3 | attrLo<<2 | tileHi<<1 | tileLo
	}

	// FIXME: Do Sprite Pixel output
	p.VirtualDisplay[displayIdx] = p.getColorIdxFromPalette(paletteRAMIDx)
}

func (p *PPU) doFetch(patternTableHalf types.Word) types.Word {
	const singleTileMaxDot = 8

	localizedDot := ((p.dot - 1) % singleTileMaxDot) + 1
	switch localizedDot {
	case 2:
		ntAddr := 0x2000 | (p.currentVRAMAddr & 0x0FFF)
		p.nTDataLatch = p.readPRGMemory(ntAddr)
	case 4:
		atAddr := 0x23C0 | (p.currentVRAMAddr & 0x0C00) | ((p.currentVRAMAddr >> 4) & 0x38) | ((p.currentVRAMAddr >> 2) & 0x07)
		p.aTDataLatch = p.readPRGMemory(atAddr)
	case 6:
		ptAddr := p.calcPatternTableAddr(0, patternTableHalf)
		p.patternDataLoLatch = p.readPRGMemory(ptAddr)
	case 8:
		ptAddr := p.calcPatternTableAddr(8, patternTableHalf)
		p.patternDataHiLatch = p.readPRGMemory(ptAddr)
		p.incrementX()
	}

	if p.dot == visibleDotsMax {
		p.incrementY()
	}

	return localizedDot
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

func (p *PPU) getColorIdxFromPalette(paletteRAMIdx byte) byte {
	return p.readPRGMemory(paletteRAMLoAddr | types.Word(paletteRAMIdx))
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
	tileNo := types.Word(p.nTDataLatch) << 4
	return patternTableHalf | tileNo | bitPlane | fineY
}
