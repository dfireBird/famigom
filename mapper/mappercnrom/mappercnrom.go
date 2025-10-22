package mappercnrom

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/mapper/mapperlib"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/program"
	"github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x03
)

type MapperCNROM struct {
	prgROM []byte
	chrROM []byte

	chrBank []byte

	// Used for only one game: Hayauchi Super Igo
	isPrgRAM   bool
	prgRAM     []byte
	prgRAMSize types.Word
}

func CreateMapperCNROM(program *program.Program) *MapperCNROM {
	isPRGRAM := false
	prgRAM := []byte{}
	prgRAMSize := types.Word(0)

	if program.IsINES2 && program.PrgRAMSize > 0 {
		isPRGRAM = true
		prgRAMSize = program.PrgRAMSize
		prgRAM = make([]byte, prgRAMSize)
	}
	chrBank := program.ChrRom[:constants.Kib8]

	return &MapperCNROM{
		prgROM:  program.PrgRom,
		chrROM:  program.ChrRom,
		chrBank: chrBank,

		isPrgRAM:   isPRGRAM,
		prgRAM:     prgRAM,
		prgRAMSize: prgRAMSize,
	}
}

func (m *MapperCNROM) ReadMemory(addr types.Word) (bool, byte) {
	if m.isPrgRAM && mapperlib.IsPRGRAMAddr(addr) {
		idx := addr - constants.LowPrgRAMAddr
		return true, m.prgRAM[idx%m.prgRAMSize]
	} else if mapperlib.IsPRGROMAddr(addr) {
		idx := addr - constants.LowPrgROMAddr
		return true, m.prgROM[idx]
	}

	return false, 0
}

func (m *MapperCNROM) WriteMemory(addr types.Word, value byte) {
	if m.isPrgRAM && mapperlib.IsPRGRAMAddr(addr) {
		idx := addr - constants.LowPrgRAMAddr
		m.prgRAM[idx%m.prgRAMSize] = value
	} else if mapperlib.IsPRGROMAddr(addr) {
		chrBankSel := value & 0x03
		chrBankStartIdx := uint(chrBankSel) * constants.Kib8
		chrBankEndIdx := chrBankStartIdx + constants.Kib8
		m.chrBank = m.chrROM[chrBankStartIdx:chrBankEndIdx]
	}
}

func (m *MapperCNROM) ReadCHRMemory(addr types.Word) (bool, byte) {
	chrBank := [constants.Kib8]byte(m.chrBank)
	return mapperlib.GenericCHRRead(&chrBank, addr)
}

func (m *MapperCNROM) WriteCHRMemory(addr types.Word, value byte) {
	chrBank := [constants.Kib8]byte(m.chrBank)
	mapperlib.GenericCHRWrite(&chrBank, addr, value)
}

func (m *MapperCNROM) GetMapperNum() byte {
	return mapperNum
}

func (m *MapperCNROM) SetMirroringUpdateCallback(func(nametable.NametableMirroring)) {}

func (m *MapperCNROM) CPUStep() {}

func (m *MapperCNROM) PPUStep() {}
