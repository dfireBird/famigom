package mappernrom

import (
	"github.com/dfirebird/famigom/program"
	. "github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x00

	prgRAMSize = 4096

	lowAddr  = 0x6000
	highAddr = 0xFFFF

	lowPrgRamAddr  = lowAddr
	highPrgRamAddr = 0x7FFF

	lowPrgRomAddr  = 0x8000
	highPrgRomAddr = highAddr

	lowChrRomAddr = 0x0000
	hiChrRomAddr  = 0x1FFF

	kib8  = 8192
	kib16 = kib8 * 2
	kib32 = kib16 * 2
)

type MapperNROM struct {
	prgRom        []byte
	isPrgRom32kib bool

	prgRam [prgRAMSize]byte

	chrRom [kib8]byte
	// can use NametableArrangement but should be inversed i.e if V -> H or if H -> V
	nametableMirroring program.NametableArrangement
}

func CreateMapperNRom(prgRom []byte, chrRom []byte, nametableMirroring program.NametableArrangement) *MapperNROM {
	isPrgRom32kib := len(prgRom) == kib32

	var chr [kib8]byte
	if (len(chrRom) == 0) {
		chr = [kib8]byte{} // CHR RAM
	} else {
		chr = [8192]byte(chrRom)
	}
	mapper := MapperNROM{
		prgRom:             prgRom,
		isPrgRom32kib:      isPrgRom32kib,
		prgRam:             [4096]byte{},
		chrRom:             chr,
		nametableMirroring: nametableMirroring,
	}

	return &mapper
}

func (m *MapperNROM) ReadMemory(addr Word) (bool, byte) {
	if lowPrgRamAddr <= addr && addr <= highPrgRamAddr {
		prgRamAddr := addr - lowPrgRamAddr
		return true, m.prgRam[prgRamAddr]
	} else if lowPrgRomAddr <= addr && addr <= highPrgRomAddr {
		prgRomAddr := addr - lowPrgRomAddr
		if !m.isPrgRom32kib {
			prgRomAddr = prgRomAddr % kib16
		}

		return true, m.prgRom[prgRomAddr]
	} else {
		return false, 0
	}
}

func (m *MapperNROM) WriteMemory(addr Word, value byte) {
	if lowPrgRamAddr <= addr && addr <= highPrgRamAddr {
		prgRamAddr := addr - lowPrgRamAddr
		m.prgRam[prgRamAddr] = value
	}
}

func (m *MapperNROM) ReadCHRMemory(addr Word) (bool, byte) {
	if lowChrRomAddr <= addr && addr <= hiChrRomAddr {
		return true, m.chrRom[addr]
	}
	return false, 0
}

func (m *MapperNROM) WriteCHRMemory(addr Word, value byte) {
	if lowChrRomAddr <= addr && addr <= hiChrRomAddr {
		m.chrRom[addr] = value
	}
}

func (m *MapperNROM) GetMapperNum() byte {
	return mapperNum
}
