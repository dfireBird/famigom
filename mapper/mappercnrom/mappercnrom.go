package mappercnrom

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x03
)

type MapperCNROM struct {
	prgROM []byte
	chrROM []byte

	chrBank []byte

	// Used for only one game:  Hayauchi Super Igo
	isPrgRAM bool
	prgRAM   []byte
}

func CreateMapperCNROM(prgROM, chrROM []byte) *MapperCNROM {
	chrBank := chrROM[:constants.Kib8]

	return &MapperCNROM{
		prgROM:  prgROM,
		chrROM:  chrROM,
		chrBank: chrBank,

		// FIXME: implement NES 2.0 header and then implement it properly here
		isPrgRAM: false,
		prgRAM:   []byte{},
	}
}

func (m *MapperCNROM) ReadMemory(addr types.Word) (bool, byte) {
	if m.isPrgRAM && (constants.LowPrgRAMAddr <= addr && addr <= constants.HighPrgRAMAddr) {
		idx := addr - constants.LowPrgRAMAddr
		return true, m.prgRAM[idx]
	}

	if constants.LowPrgROMAddr <= addr && addr <= constants.HighPrgROMAddr {
		idx := addr - constants.LowPrgROMAddr
		return true, m.prgROM[idx]
	}

	return false, 0
}

func (m *MapperCNROM) WriteMemory(addr types.Word, value byte) {
	if constants.LowPrgROMAddr <= addr && addr <= constants.HighPrgROMAddr {
		chrBankSel := value & 0x03
		chrBankStartIdx := uint(chrBankSel) * constants.Kib8
		chrBankEndIdx := chrBankStartIdx + constants.Kib8
		m.chrBank = m.chrROM[chrBankStartIdx:chrBankEndIdx]
	}
}

func (m *MapperCNROM) ReadCHRMemory(addr types.Word) (bool, byte) {
	if constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr {
		return true, m.chrBank[addr]
	}
	return false, 0
}

func (m *MapperCNROM) WriteCHRMemory(addr types.Word, value byte) {
	if constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr {
		m.chrBank[addr] = value
	}
}

func (m *MapperCNROM) GetMapperNum() byte {
	return mapperNum
}

func (m *MapperCNROM) SetMirroringUpdateCallback(func(nametable.NametableMirroring)) {}
