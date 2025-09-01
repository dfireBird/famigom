package ppu

import "github.com/dfirebird/famigom/types"

func (p *PPU) doBackgroundFetch(patternTableHalf types.Word) types.Word {
	const singleTileMaxDot = 8

	localizedDot := ((p.dot - 1) % singleTileMaxDot) + 1
	switch localizedDot {
	case 2:
		ntAddr := 0x2000 | (p.currentVRAMAddr & 0x0FFF)
		p.nTDataLatch = p.readCHRMemory(ntAddr)
	case 4:
		atAddr := 0x23C0 | (p.currentVRAMAddr & 0x0C00) | ((p.currentVRAMAddr >> 4) & 0x38) | ((p.currentVRAMAddr >> 2) & 0x07)
		p.aTDataLatch = p.extractBgAttr(p.readCHRMemory(atAddr))
	case 6:
		ptAddr := p.calcPatternTableAddr(0, patternTableHalf)
		p.patternDataLoLatch = p.readCHRMemory(ptAddr)
	case 8:
		ptAddr := p.calcPatternTableAddr(8, patternTableHalf)
		p.patternDataHiLatch = p.readCHRMemory(ptAddr)
		p.incrementX()
	}

	if p.dot == visibleDotsMax {
		p.incrementY()
	}

	return localizedDot
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

		p.currentVRAMAddr = (p.currentVRAMAddr &^ coarseYScrollMask) | (coarseY << coarseYScrollShift)
	}
}

func (p *PPU) extractBgAttr(attribute byte) byte {
	coarseX := p.currentVRAMAddr & coarseXScrollMask
	coarseY := (p.currentVRAMAddr & coarseYScrollMask) >> coarseYScrollShift
	atLoPos := (coarseX & 0x02) + ((coarseY & 0x02) << 1)
	atHiPos := atLoPos + 1
	attrLo := (attribute & (1 << atLoPos)) >> byte(atLoPos)
	attrHi := (attribute & (1 << atHiPos)) >> byte(atHiPos)
	return (attrHi<<1 | attrLo)
}

func (p *PPU) calcPatternTableAddr(bitPlane, patternTableHalf types.Word) types.Word {
	fineY := (p.currentVRAMAddr & fineYScrollMask) >> fineYScrollShift
	tileNo := types.Word(p.nTDataLatch) << 4
	return patternTableHalf | tileNo | bitPlane | fineY
}

func (p *PPU) DrawNametable() {
	oldVramAddr := p.currentVRAMAddr
	oldNtDataLatch := p.nTDataLatch
	p.currentVRAMAddr = 0x00

	for i := range totalDots {
		ntAddr := 0x2000 | (p.currentVRAMAddr & 0x0FFF)
		ntData := p.readCHRMemory(ntAddr)

		atAddr := 0x23C0 | (p.currentVRAMAddr & 0x0C00) | ((p.currentVRAMAddr >> 4) & 0x38) | ((p.currentVRAMAddr >> 2) & 0x07)
		atData := p.readCHRMemory(atAddr)

		p.nTDataLatch = ntData
		ptAddrLo := p.calcPatternTableAddr(0, p.getBackgroundPatternTableAddr())
		pTDataLo := p.readCHRMemory(ptAddrLo)

		ptAddrHi := p.calcPatternTableAddr(8, p.getBackgroundPatternTableAddr())
		pTDataHi := p.readCHRMemory(ptAddrHi)

		fineX := 7 - (i % 8) // to emulate shifting not actually scrolling
		fineXBitSelect := byte(1) << fineX
		tileLo := (pTDataLo & fineXBitSelect) >> fineX
		tileHi := (pTDataHi & fineXBitSelect) >> fineX

		coarseXBit1 := (p.currentVRAMAddr & coarseXScroll1Bit) >> coarseXScroll1Shift
		coarseYBit1 := (p.currentVRAMAddr & coarseYScroll1Bit) >> coarseYScroll1Shift
		atLoPos := coarseXBit1*2 + coarseYBit1*4
		atHiPos := atLoPos + 1

		attrLo := (atData & (1 << atLoPos)) >> byte(atLoPos)
		attrHi := (atData & (1 << atHiPos)) >> byte(atHiPos)

		paletteRAMIDx := bgPaletteMSB<<4 | attrHi<<3 | attrLo<<2 | tileHi<<1 | tileLo
		p.VirtualDisplay[i] = p.getColorIdxFromPalette(paletteRAMIDx)

		if i%8 == 7 {
			if p.currentVRAMAddr&coarseXScrollMask == coarseXScrollMask {
				p.currentVRAMAddr &^= coarseXScrollMask
			} else {
				p.currentVRAMAddr += 1
			}
		}

		if i%visibleDotsMax == 255 {
			p.incrementY()
		}
	}

	p.nTDataLatch = oldNtDataLatch
	p.currentVRAMAddr = oldVramAddr
}
