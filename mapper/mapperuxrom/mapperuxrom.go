package mapperuxrom

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x02

	fixedBankStart = 0xC000
)

type MapperUxROM struct {
	maxBanks byte

	prgROM []byte
	chrROM [constants.Kib8]byte

	bank      []byte
	fixedBank []byte
}

func CreateMapperUxROM(prgROM, chrROM []byte) *MapperUxROM {
	var chr [constants.Kib8]byte
	if len(chrROM) == 0 {
		chr = [constants.Kib8]byte{} // CHR RAM
	} else {
		chr = [constants.Kib8]byte(chrROM)
	}

	size := len(prgROM)
	fixedBank := prgROM[size-constants.Kib16:]
	maxBanks := size / constants.Kib16
	mapper := MapperUxROM{
		maxBanks:  byte(maxBanks),
		prgROM:    prgROM,
		chrROM:    chr,
		bank:      prgROM[:constants.Kib16],
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
		bankSel := (value & 0x0F) % m.maxBanks
		bankStartIdx := uint(bankSel) * constants.Kib16
		bankEndIdx := bankStartIdx + constants.Kib16
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

func (m *MapperUxROM) SetMirroringUpdateCallback(func(nametable.NametableMirroring)) {}
