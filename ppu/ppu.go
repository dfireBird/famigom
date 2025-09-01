package ppu

import (
	"github.com/dfirebird/famigom/bus"
	"github.com/dfirebird/famigom/log"
	"github.com/dfirebird/famigom/program"
	"github.com/dfirebird/famigom/types"
)

const (
	ioLoAddr = 0x2000
	ioHiAddr = 0x3FFF

	oamMemorySize    = 256
	secondaryOAMSize = 32

	visibleScanlineMax = 240
	visibleDotsMax     = 256
	totalDots          = visibleScanlineMax * visibleDotsMax

	dotsTillFetchesUsed = 248

	preRenderScanLine  = 261
	visibleScanLineLo  = 0
	visibleScanLineHi  = 239
	postRenderScanLine = 240
	vblankScanLineLo   = 241
	vblankScanLineHi   = 260

	nextScanline1stTileDot = 328
	nextScanline2ndTileDot = 336

	maxDotsInALine = 340

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

	dotsToDecayOpenBusReg = 3221591

	spritePritorityMask = 0x20
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

	spriteOverflow bool
	sprite0Hit     bool
	vblankFlag     bool

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

	tilePatternShiftRegisterLo   types.Word
	tilePatternShiftRegisterHi   types.Word
	tileAttributeShiftRegisterLo types.Word
	tileAttributeShiftRegisterHi types.Word

	openBusDecayRegister byte
	openBusDecayTime     uint32

	secondaryOAM      [secondaryOAMSize]byte
	secondaryOAMIdx   byte
	spriteIdx         byte
	spritePatternData [secondaryOAMSize]byte
	sprite0InNextLine bool

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
	p.sprite0Hit = false
	p.spriteOverflow = true

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
	log.GetLoggerWithSpan("ppu").Debugf("(x, y): (%03d, %03d) v: 0x%04X t: 0x%04X PPUCTRL: 0x%04X PPUMASK: 0x%04X PPUSTAT: 0x%04X ND: 0x%02X AD: 0x%02X pl: 0x%04X ph: 0x%04X al: 0x%04X ah: 0x%04X",
		p.dot, p.line, p.currentVRAMAddr, p.tempVRAMAddr, p.ppuCtrl, p.ppuMask, p.getPPUStatus(),
		p.nTDataLatch, p.aTDataLatch, p.tilePatternShiftRegisterLo, p.tilePatternShiftRegisterHi,
		p.tileAttributeShiftRegisterLo, p.tileAttributeShiftRegisterHi,
	)

	if p.openBusDecayTime == 0 {
		p.openBusDecayRegister = 0
	} else {
		p.openBusDecayTime -= 1
	}

	// FIXME: Do rendering of sprite pixels
	if (visibleScanLineLo <= p.line && p.line <= visibleScanLineHi) || p.line == preRenderScanLine {
		if p.dot == 0 {
			p.incrementDot()
			return
		}

		if p.line == preRenderScanLine && p.dot == 1 {
			p.vblankFlag = false
			p.sprite0Hit = false
			p.spriteOverflow = false
		}

		if p.isCurrentlyRendering() {
			p.outputPixel()
		}

		if p.isRenderingEnabled() && (1 <= p.dot && p.dot <= visibleDotsMax) {
			ld := p.doBackgroundFetch(p.getBackgroundPatternTableAddr())
			if ld == 8 && 1 <= p.dot && p.dot <= dotsTillFetchesUsed {
				p.tilePatternShiftRegisterLo &= 0xFF00
				p.tilePatternShiftRegisterHi &= 0xFF00
				p.tileAttributeShiftRegisterLo &= 0xFF00
				p.tileAttributeShiftRegisterHi &= 0xFF00

				p.tilePatternShiftRegisterLo |= types.Word(p.patternDataLoLatch)
				p.tilePatternShiftRegisterHi |= types.Word(p.patternDataHiLatch)
				p.tileAttributeShiftRegisterLo |= types.Word(extractNthBitAndRepeat(0, p.aTDataLatch))
				p.tileAttributeShiftRegisterHi |= types.Word(extractNthBitAndRepeat(1, p.aTDataLatch))
			}
			p.spriteEvaulvation()
		}

		if p.isRenderingEnabled() && (spriteDotLo <= p.dot && p.dot <= spirteDotHi) {
			p.doSpriteFetch()
		}

		if p.isRenderingEnabled() && (tilesForNextScanLineLo <= p.dot && p.dot <= tilesForNextScanLineHi) {
			p.doBackgroundFetch(p.getBackgroundPatternTableAddr())

			if p.dot == nextScanline1stTileDot {
				p.tilePatternShiftRegisterLo = types.Word(p.patternDataLoLatch) << 8
				p.tilePatternShiftRegisterHi = types.Word(p.patternDataHiLatch) << 8
				p.tileAttributeShiftRegisterLo = types.Word(extractNthBitAndRepeat(0, p.aTDataLatch)) << 8
				p.tileAttributeShiftRegisterHi = types.Word(extractNthBitAndRepeat(1, p.aTDataLatch)) << 8
			}

			if p.dot == nextScanline2ndTileDot {
				p.tilePatternShiftRegisterLo |= types.Word(p.patternDataLoLatch)
				p.tilePatternShiftRegisterHi |= types.Word(p.patternDataHiLatch)
				p.tileAttributeShiftRegisterLo |= types.Word(extractNthBitAndRepeat(0, p.aTDataLatch))
				p.tileAttributeShiftRegisterHi |= types.Word(extractNthBitAndRepeat(1, p.aTDataLatch))
			}
		}

		if p.isRenderingEnabled() && p.dot == horiPosCopyDot {
			p.currentVRAMAddr &^= coarseXScrollMask                   // coarse X = 0
			p.currentVRAMAddr |= (p.tempVRAMAddr & coarseXScrollMask) // copy coarse X

			p.currentVRAMAddr &^= horizontalNametableMask                   // clear horizontal Nametable bit
			p.currentVRAMAddr |= (p.tempVRAMAddr & horizontalNametableMask) // copy	nt bit
		}

		if p.isRenderingEnabled() && p.line == preRenderScanLine && (veriPosCopyDotLo <= p.dot && p.dot <= veriPosCopyDotHi) {
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
			if p.isNMIEnabled() {
				(*p.nmiCallback)()
			}
		}
	}

	p.incrementDot()
}

func (p *PPU) outputPixel() {
	paletteRAMIDx := byte(0)
	displayIdx := (p.line*visibleDotsMax + p.dot) - 1
	isBgOpaque := false

	if p.isBgRenderingEnabled() {
		shiftOfFineX := types.Word(15 - (p.fineX % 8))
		fineXBitSelect := types.Word(1) << shiftOfFineX // to make fineX value 0 to pull out MSB
		tileLo := (p.tilePatternShiftRegisterLo & fineXBitSelect) >> shiftOfFineX
		tileHi := (p.tilePatternShiftRegisterHi & fineXBitSelect) >> shiftOfFineX

		attrLo := (p.tileAttributeShiftRegisterLo & fineXBitSelect) >> shiftOfFineX
		attrHi := (p.tileAttributeShiftRegisterHi & fineXBitSelect) >> shiftOfFineX
		log.GetLoggerWithSpan("ppu").Debugf("Attribute: %d %d", attrHi, attrLo)

		isBgOpaque = (tileHi<<1 | tileLo) != 0

		paletteRAMIDx = byte(bgPaletteMSB<<4 | attrHi<<3 | attrLo<<2 | tileHi<<1 | tileLo)

		p.tilePatternShiftRegisterLo = (p.tilePatternShiftRegisterLo << 1) | 1
		p.tilePatternShiftRegisterHi = (p.tilePatternShiftRegisterHi << 1) | 1
		p.tileAttributeShiftRegisterLo = (p.tileAttributeShiftRegisterLo << 1) | 1
		p.tileAttributeShiftRegisterHi = (p.tileAttributeShiftRegisterHi << 1) | 1
	}

	if p.isSpriteRenderingEnabled() {
		spritePaletteRAMIdx := byte(0)
		for i := range byte(8) {
			dotVal := p.dot - 1
			xCoord := uint16(p.spritePatternData[i*4+0])
			priority := p.spritePatternData[i*4+1] & spritePritorityMask

			isInXCoordRange := xCoord <= dotVal && dotVal <= xCoord+7
			if xCoord != 0xFF && isInXCoordRange {
				attribute := p.spritePatternData[i*4+1] & 0x3
				spritePtDataLo := &p.spritePatternData[i*4+2]
				spritePtDataHi := &p.spritePatternData[i*4+3]

				shift := 7
				bitSelect := byte(1) << shift
				tileLo := (*spritePtDataLo & bitSelect) >> shift
				tileHi := (*spritePtDataHi & bitSelect) >> shift

				log.GetLoggerWithSpan("ppu").Debugf("x: %d, tl: 0b%08b th: 0b%08b at: 0b%02b", xCoord, *spritePtDataLo, *spritePtDataHi, attribute)

				*spritePtDataLo = (*spritePtDataLo << 1) | 1
				*spritePtDataHi = (*spritePtDataHi << 1) | 1

				isSpriteOpaque := (tileHi<<1 | tileLo) != 0

				isLeftClip := (p.isLeftClipBg() || p.isLeftClipSprite()) && (dotVal <= 7)
				isRightEdge := p.dot >= 256
				if !p.sprite0Hit && i == 0 && p.sprite0InNextLine && !isLeftClip && !isRightEdge && isBgOpaque && isSpriteOpaque {
					p.sprite0Hit = true
				}

				isInFg := priority == 0
				isInBg := priority == spritePritorityMask
				if spritePaletteRAMIdx == 0 && isSpriteOpaque && (isInFg || (!isBgOpaque && isInBg)){
					spritePaletteRAMIdx = spritePaletteMSB<<4 | attribute<<2 | tileHi<<1 | tileLo
					paletteRAMIDx = spritePaletteMSB<<4 | attribute<<2 | tileHi<<1 | tileLo
				}
			}
		}
	}

	p.VirtualDisplay[displayIdx] = p.getColorIdxFromPalette(paletteRAMIDx)
}

func (p *PPU) incrementDot() {
	newDot := (p.dot + 1) % (maxDotsInALine + 1)
	newLineVal := p.line

	if p.oddFrame && newDot == (maxDotsInALine-1) && newLineVal == preRenderScanLine {
		newDot = 0
		newLineVal = 0
	} else if newDot == 0 {
		newLineVal = (newLineVal + 1) % (preRenderScanLine + 1)
	}

	p.dot = newDot
	p.line = newLineVal
}

func (p *PPU) getColorIdxFromPalette(paletteRAMIdx byte) byte {
	return p.readCHRMemory(paletteRAMLoAddr | types.Word(paletteRAMIdx))
}

func (p *PPU) readOAMMemory(addr byte) byte {
	value := p.oamMemory[addr]
	if addr%4 == 2 {
		value &= 0b1110_0011
	}
	return value
}

func (p *PPU) writeOAMMemory(addr, value byte) {
	p.oamMemory[addr] = value
}

func (p *PPU) readCHRMemory(addr types.Word) byte {
	return p.chrMemoryBus.ReadCHRMemory(addr)
}

func (p *PPU) writeCHRMemory(addr types.Word, value byte) {
	p.chrMemoryBus.WriteCHRMemory(addr, value)
}

func (p *PPU) isCurrentlyRendering() bool {
	isVisibleScanline := p.line < visibleScanlineMax
	isVisibleDot := p.dot <= visibleDotsMax

	return p.isRenderingEnabled() && (isVisibleDot && isVisibleScanline)
}

func (p *PPU) isRenderingEnabled() bool {
	return p.isBgRenderingEnabled() || p.isSpriteRenderingEnabled()
}

func extractNthBitAndRepeat(n, val byte) byte {
	mask := byte(1) << n
	nthBit := (val & mask) >> n

	if nthBit == 0 {
		return 0
	} else {
		return 0xFF
	}
}
