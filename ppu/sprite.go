package ppu

import (
	"math/bits"

	"github.com/dfirebird/famigom/log"
	"github.com/dfirebird/famigom/types"
)

const (
	horizontalFlipOAMMask = 0x40
	verticalFlipOAMMask   = 0x80
)

func (p *PPU) spriteEvaluation() {
	if 1 <= p.line && p.line <= visibleDotsMax {
		if 1 <= p.dot && p.dot <= 64 {
			if p.dot%2 == 0 {
				if idx := (p.dot / 2) - 1; idx%4 == 0 {
					p.secondaryOAM[idx] = 0xFF
				} else {
					p.secondaryOAM[idx] = 0xFF
				}
			}

			if p.dot == 64 {
				p.secondaryOAMIdx = 0
				p.spriteIdx = 0
			}
		} else if 65 <= p.dot && p.dot <= 256 {
			yCoord := uint16(p.readOAMMemory(p.spriteIdx * 4))
			spriteHeight := getSpriteHeight(p.getSpriteSize())
			if p.spriteIdx < 64 && yCoord <= p.line && p.line <= yCoord+spriteHeight { // we can check with current line num, since Y-coord is subtracted by 1
				if p.secondaryOAMIdx <= 7 {
					p.sprite0InNextLine = p.spriteIdx == 0

					for i := range byte(4) {
						p.secondaryOAM[p.secondaryOAMIdx*4+i] = p.readOAMMemory(p.spriteIdx*4 + i)
					}
					log.GetLoggerWithSpan("ppu").Debugf("Y: %d, tid: 0x%02X, attr: 0b%08b x: %d", p.secondaryOAM[p.secondaryOAMIdx*4+0], p.secondaryOAM[p.secondaryOAMIdx*4+1], p.secondaryOAM[p.secondaryOAMIdx*4+2], p.secondaryOAM[p.secondaryOAMIdx*4+3])
				} else {
					p.spriteOverflow = true
				}
				p.secondaryOAMIdx += 1
			}
			p.spriteIdx = incrementSpriteIdx(p.spriteIdx)

			if p.dot == 256 {
				p.secondaryOAMIdx = 0
			}
		}
	}
}

func (p *PPU) doSpriteFetch() {
	const singleTileMaxDot = 8

	spriteTileNo := p.secondaryOAM[p.secondaryOAMIdx*4+1]
	yCoord := p.secondaryOAM[p.secondaryOAMIdx*4+0]

	attributes := p.secondaryOAM[p.secondaryOAMIdx*4+2]
	vertFlip := (attributes & verticalFlipOAMMask) == verticalFlipOAMMask
	horiFlip := (attributes & horizontalFlipOAMMask) == horizontalFlipOAMMask

	localizedDot := ((p.dot - 1) % singleTileMaxDot) + 1
	switch localizedDot {
	case 2:
		p.spritePatternData[p.secondaryOAMIdx*4+0] = p.secondaryOAM[p.secondaryOAMIdx*4+3] // x-coord
		p.readCHRMemory(0x2000 | (p.currentVRAMAddr & 0x0FFF))                             // dummy reads
	case 4:
		p.spritePatternData[p.secondaryOAMIdx*4+1] = attributes
		p.readCHRMemory(0x2000 | (p.currentVRAMAddr & 0x0FFF)) // dummy reads
	case 6:
		ptAddr := p.calcSpritePatternAddr(0, types.Word(spriteTileNo), yCoord, vertFlip)
		ptData := flipHorizontallyIfTrue(horiFlip, p.readCHRMemory(ptAddr))

		p.spritePatternData[p.secondaryOAMIdx*4+2] = ptData
	case 8:
		ptAddr := p.calcSpritePatternAddr(8, types.Word(spriteTileNo), yCoord, vertFlip)
		ptData := flipHorizontallyIfTrue(horiFlip, p.readCHRMemory(ptAddr))

		p.spritePatternData[p.secondaryOAMIdx*4+3] = ptData
		p.secondaryOAMIdx = incrementSecondaryOAMIdx(p.secondaryOAMIdx)
	}
}

func (p *PPU) calcSpritePatternAddr(bitPlane, tileNo types.Word, yCoord byte, vertFlip bool) types.Word {
	fineY := p.line - uint16(yCoord)
	spriteHeight := getSpriteHeight(p.getSpriteSize())
	if vertFlip {
		fineY = spriteHeight - fineY
	}

	isNextTileIn8x16 := false
	if fineY > 7 {
		fineY %= 8
		isNextTileIn8x16 = true
	}

	var patternTableHalf types.Word
	if p.getSpriteSize() == Sprite8x8 {
		patternTableHalf = p.getSpritePatternTableAddr()
	} else {
		bit0 := tileNo & 1
		if bit0 == 0 {
			patternTableHalf = 0x0000
		} else {
			patternTableHalf = 0x1000
			tileNo &^= 1
		}
	}

	if isNextTileIn8x16 {
		tileNo |= 1
	}
	return patternTableHalf | (tileNo << 4) | bitPlane | fineY
}

func flipHorizontallyIfTrue(isFlip bool, ptData byte) byte {
	if isFlip {
		return bits.Reverse8(ptData)
	}
	return ptData
}

func incrementSpriteIdx(spriteIdx byte) byte {
	if spriteIdx == 64 {
		return spriteIdx
	}
	return spriteIdx + 1
}

func incrementSecondaryOAMIdx(idx byte) byte {
	if idx == 7 {
		return idx
	}
	return idx + 1
}

func getSpriteHeight(size SpriteSize) types.Word {
	if size == Sprite8x8 {
		return 7
	} else {
		return 15
	}
}
