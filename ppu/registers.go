package ppu

import (
	"github.com/dfirebird/famigom/types"
)

const (
	maxIoRegisterAddr = 0x2007

	nametableSelectMask = 0x03

	fineScrollMask   byte = 0x07
	coarseScrollMask      = ^fineScrollMask

	coarseXScrollMask  = 0x001F
	coarseYScrollMask  = 0x03E0
	fineYScrollMask    = 0x7000
	coarseYScrollShift = 5
	fineYScrollShift   = 12

	firstPPUAddrWriteMask = 0x1F
	hiByteVRAMMask        = 0xFF00
	loByteVRAMMask        = 0x00FF
	hiByteVRAMShift       = 8
)

func (p *PPU) ReadMemory(addr types.Word) (bool, byte) {
	if ioLoAddr <= addr && addr <= ioHiAddr {
		effectiveAddr := addr & maxIoRegisterAddr
		switch effectiveAddr {
		case PPUSTATUS:
			status := p.getPPUStatus()

			p.vblankFlag = false
			p.isFirstWrite = true

			return true, status
		case PPUDATA:
			returnData := p.ppuData

			p.ppuData = p.readPRGMemory(p.currentVRAMAddr)
			p.currentVRAMAddr += p.getVRAMAddrIncr()
			if paletteRAMLoAddr <= addr && addr <= paletteRAMHiAddr {
				returnData = p.ppuData
			}

			return true, returnData
		default:
			return true, 0xFF
		}
	}
	return false, 0
}

func (p *PPU) WriteMemory(addr types.Word, value byte) {
	if ioLoAddr <= addr && addr <= ioHiAddr {
		effectiveAddr := addr & maxIoRegisterAddr
		switch effectiveAddr {
		case PPUCTRL:
			nametableBits := value & nametableSelectMask

			p.tempVRAMAddr &^= nametableSelectMask << 10
			p.tempVRAMAddr |= types.Word(nametableBits) << 10

			p.ppuCtrl = value

		case PPUMASK:
			p.ppuMask = value

		case OAMADDR:
			p.oamAddr = value

		case OAMDATA:
			p.oamMemory[p.oamAddr] = value
			p.oamAddr += 1

		case PPUSCROLL:
			fineScroll := value & fineScrollMask
			coarseScroll := (value & coarseScrollMask) >> 3
			if p.isFirstWrite {
				p.fineX = fineScroll

				p.tempVRAMAddr &^= coarseXScrollMask
				p.tempVRAMAddr |= types.Word(coarseScroll)
			} else {
				p.tempVRAMAddr &^= fineYScrollMask
				p.tempVRAMAddr |= types.Word(fineScroll) << fineYScrollShift

				p.tempVRAMAddr &^= coarseYScrollMask
				p.tempVRAMAddr |= types.Word(coarseScroll) << coarseYScrollShift
			}
			p.isFirstWrite = !p.isFirstWrite

		case PPUADDR:
			if p.isFirstWrite {
				maskedVal := value & firstPPUAddrWriteMask

				p.tempVRAMAddr &^= hiByteVRAMMask
				p.tempVRAMAddr |= types.Word(maskedVal) << hiByteVRAMShift
			} else {
				p.tempVRAMAddr &^= loByteVRAMMask
				p.tempVRAMAddr |= types.Word(value)

				p.currentVRAMAddr = p.tempVRAMAddr
			}
			p.isFirstWrite = !p.isFirstWrite

		case PPUDATA:
			p.writePRGMemory(p.currentVRAMAddr, value)
			p.currentVRAMAddr += p.getVRAMAddrIncr()
		}

	}
}

func (p *PPU) getBaseNametableAddr() types.Word {
	nametablebits := p.ppuCtrl & nametableSelectMask

	var addr types.Word
	switch nametablebits {
	case 0:
		addr = 0x2000
	case 1:
		addr = 0x2400
	case 2:
		addr = 0x2800
	case 3:
		addr = 0x2C00
	}

	return addr
}

func (p *PPU) getVRAMAddrIncr() types.Word {
	const (
		mask   byte       = 0x04
		valIf0 types.Word = 1
		valIf1 types.Word = 32
	)

	return maskAndGetValueFromPpuCtrl(p.ppuCtrl, mask, valIf1, valIf0)
}

func (p *PPU) getSpritePatternTableAddr() types.Word {
	const (
		mask   byte       = 0x08
		valIf0 types.Word = 0x0000
		valIf1 types.Word = 0x1000
	)

	return maskAndGetValueFromPpuCtrl(p.ppuCtrl, mask, valIf1, valIf0)
}

func (p *PPU) getBackgroundPatternTableAddr() types.Word {
	const (
		mask   byte       = 0x10
		valIf0 types.Word = 0x0000
		valIf1 types.Word = 0x1000
	)

	return maskAndGetValueFromPpuCtrl(p.ppuCtrl, mask, valIf1, valIf0)
}

func (p *PPU) getSpriteSize() SpriteSize {
	const mask byte = 0x20
	return maskAndGetValueFromPpuCtrl(p.ppuCtrl, mask, Sprite8x16, Sprite8x8)
}

func (p *PPU) getPPUStatus() (ppuStatus byte) {
	ppuStatus |= boolToByte(p.vblankFlag) << 7
	ppuStatus |= boolToByte(p.spite0Hit) << 6
	ppuStatus |= boolToByte(p.spiteOverflow) << 5
	return
}

func (p *PPU) isNMIEnabled() bool {
	const mask byte = 0x80
	return p.ppuCtrl&mask == mask
}

func (p *PPU) isGreyscale() bool {
	const mask byte = 0x01
	return p.ppuMask&mask == mask
}

func (p *PPU) isBgRenderingEnabled() bool {
	const mask byte = 0x08
	return p.ppuMask&mask == mask
}

func (p *PPU) isSpriteRenderingEnabled() bool {
	const mask byte = 0x10
	return p.ppuMask&mask == mask
}

func maskAndGetValueFromPpuCtrl[T any](ppuCtrl, mask byte, valIf1, valIf0 T) T {
	if ppuCtrl&mask == mask {
		return valIf1
	} else {
		return valIf0
	}
}

func boolToByte(val bool) byte {
	if val {
		return 1
	} else {
		return 0
	}
}
