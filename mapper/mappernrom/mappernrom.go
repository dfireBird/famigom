package mappernrom

import (
	"github.com/dfirebird/famigom/constants"
	. "github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x00

	prgRAMSize = constants.Kib4
)

type MapperNROM struct {
	prgRom        []byte
	isPrgRom32kib bool

	prgRAM [prgRAMSize]byte

	chrRom [constants.Kib8]byte
}

func CreateMapperNRom(prgRom []byte, chrRom []byte) *MapperNROM {
	isPrgRom32kib := len(prgRom) == constants.Kib32

	var chr [constants.Kib8]byte
	if (len(chrRom) == 0) {
		chr = [constants.Kib8]byte{} // CHR RAM
	} else {
		chr = [constants.Kib8]byte(chrRom)
	}
	mapper := MapperNROM{
		prgRom:             prgRom,
		isPrgRom32kib:      isPrgRom32kib,
		prgRAM:             [prgRAMSize]byte{},
		chrRom:             chr,
	}

	return &mapper
}

func (m *MapperNROM) ReadMemory(addr Word) (bool, byte) {
	if constants.LowPrgRAMAddr <= addr && addr <= constants.HighPrgRAMAddr {
		prgRAMAddr := addr - constants.LowPrgRAMAddr
		return true, m.prgRAM[prgRAMAddr]
	} else if constants.LowPrgROMAddr <= addr && addr <= constants.HighPrgROMAddr {
		prgRomAddr := addr - constants.LowPrgROMAddr
		if !m.isPrgRom32kib {
			prgRomAddr = prgRomAddr % constants.Kib16
		}

		return true, m.prgRom[prgRomAddr]
	} else {
		return false, 0
	}
}

func (m *MapperNROM) WriteMemory(addr Word, value byte) {
	if constants.LowPrgRAMAddr <= addr && addr <= constants.HighPrgRAMAddr {
		prgRAMAddr := addr - constants.LowPrgRAMAddr
		m.prgRAM[prgRAMAddr] = value
	}
}

func (m *MapperNROM) ReadCHRMemory(addr Word) (bool, byte) {
	if constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr {
		return true, m.chrRom[addr]
	}
	return false, 0
}

func (m *MapperNROM) WriteCHRMemory(addr Word, value byte) {
	if constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr {
		m.chrRom[addr] = value
	}
}

func (m *MapperNROM) GetMapperNum() byte {
	return mapperNum
}
