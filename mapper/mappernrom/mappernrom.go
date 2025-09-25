package mappernrom

import (
	"github.com/dfirebird/famigom/constants"
	. "github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x00

	prgRAMSize = 4096

	kib8  = 8192
	kib16 = kib8 * 2
	kib32 = kib16 * 2
)

type MapperNROM struct {
	prgRom        []byte
	isPrgRom32kib bool

	prgRAM [prgRAMSize]byte

	chrRom [kib8]byte
}

func CreateMapperNRom(prgRom []byte, chrRom []byte) *MapperNROM {
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
		prgRAM:             [4096]byte{},
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
			prgRomAddr = prgRomAddr % kib16
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
