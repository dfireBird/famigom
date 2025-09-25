package mapperuxrom

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x02

	kib   = 1024
	kib8  = 8 * kib
	kib16 = 16 * kib

	fixedBankStart = 0xC000
)

type MapperUxROM struct {
	prgROM []byte
	chrROM [kib8]byte

	bank      []byte
	fixedBank []byte
}

func CreateMapperUxROM(prgROM, chrROM []byte) *MapperUxROM {
	var chr [kib8]byte
	if len(chrROM) == 0 {
		chr = [kib8]byte{} // CHR RAM
	} else {
		chr = [kib8]byte(chrROM)
	}

	size := len(prgROM)
	fixedBank := prgROM[size-kib16:]
	mapper := MapperUxROM{
		prgROM:    prgROM,
		chrROM:    chr,
		bank:      prgROM[:kib16],
		fixedBank: fixedBank,
	}

	return &mapper
}

func (m *MapperUxROM) ReadMemory(addr types.Word) (bool, byte) {
	if constants.LowPrgROMAddr <= addr && addr <= constants.HighPrgROMAddr {
		if addr >= fixedBankStart {
			idx := addr - fixedBankStart
			return true, m.fixedBank[idx]
		} else {
			idx := addr - constants.LowPrgROMAddr
			return true, m.bank[idx]
		}
	}
	return false, 0
}

func (m *MapperUxROM) WriteMemory(addr types.Word, value byte) {
	if constants.LowPrgROMAddr <= addr && addr <= constants.HighPrgROMAddr {
		bankSel := value & 0x0F
		bankStartIdx := uint(bankSel) * kib16
		bankEndIdx := bankStartIdx + kib16
		m.bank = m.prgROM[bankStartIdx:bankEndIdx]
	}
}

func (m *MapperUxROM) ReadCHRMemory(addr types.Word) (bool, byte) {
	if constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr {
		return true, m.chrROM[addr]
	}
	return false, 0
}

func (m *MapperUxROM) WriteCHRMemory(addr types.Word, value byte) {
	if constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr {
		m.chrROM[addr] = value
	}
}

func (m *MapperUxROM) GetMapperNum() byte {
	return mapperNum
}
