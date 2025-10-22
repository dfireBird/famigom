package mapperuxrom

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/mapper/mapperlib"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/program"
	"github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x02

	fixedBankStart = 0xC000
)

type MapperUxROM struct {
	maxBanks types.Word

	prgROM []byte
	chrROM [constants.Kib8]byte

	isPRGRAM   bool
	prgRAMSize types.Word
	prgRAM     []byte

	switchableBank []byte
	fixedBank      []byte
}

func CreateMapperUxROM(program *program.Program) *MapperUxROM {
	var chr [constants.Kib8]byte

	var prgRAM []byte
	var prgRAMSize types.Word
	isPRGRAM := false

	if program.IsINES2 {
		chr = [constants.Kib8]byte(program.ChrRom)
		if program.ChrRAMSize > 0 {
			chr = [constants.Kib8]byte{}
		}

		if program.PrgRAMSize > 0 {
			isPRGRAM = true
			prgRAMSize = program.PrgRAMSize
			prgRAM = make([]byte, program.PrgRAMSize)
		}
	} else {
		if program.ChrRomBankSize == 0 {
			chr = [constants.Kib8]byte{} // CHR RAM
		} else {
			chr = [constants.Kib8]byte(program.ChrRom)
		}
	}

	fixedBankStartIdx := (program.PrgRomBankSize - 1) * constants.Kib16
	switchableBank := program.PrgRom[:constants.Kib16]
	fixedBank := program.PrgRom[fixedBankStartIdx:]
	maxBanks := program.PrgRomBankSize

	mapper := MapperUxROM{
		maxBanks:       maxBanks,
		prgROM:         program.PrgRom,
		chrROM:         chr,
		switchableBank: switchableBank,
		fixedBank:      fixedBank,

		isPRGRAM:   isPRGRAM,
		prgRAMSize: prgRAMSize,
		prgRAM:     prgRAM,
	}

	return &mapper
}

func (m *MapperUxROM) ReadMemory(addr types.Word) (bool, byte) {
	if m.isPRGRAM && mapperlib.IsPRGRAMAddr(addr) {
		prgRAMAddr := mapperlib.CalculatePRGRAMAddr(addr)
		return true, m.prgRAM[prgRAMAddr%m.prgRAMSize]
	} else if mapperlib.IsPRGROMAddr(addr) {
		if addr >= fixedBankStart {
			idx := addr - fixedBankStart
			return true, m.fixedBank[idx]
		} else {
			idx := addr - constants.LowPrgROMAddr
			return true, m.switchableBank[idx]
		}
	}
	return false, 0
}

func (m *MapperUxROM) WriteMemory(addr types.Word, value byte) {
	if m.isPRGRAM && mapperlib.IsPRGRAMAddr(addr) {
		prgRAMAddr := mapperlib.CalculatePRGRAMAddr(addr)
		m.prgRAM[prgRAMAddr%m.prgRAMSize] = value
	} else if mapperlib.IsPRGROMAddr(addr) {
		bankSel := types.Word(value&0x0F) % m.maxBanks
		bankStartIdx := uint(bankSel) * constants.Kib16
		bankEndIdx := bankStartIdx + constants.Kib16
		m.switchableBank = m.prgROM[bankStartIdx:bankEndIdx]
	}
}

func (m *MapperUxROM) ReadCHRMemory(addr types.Word) (bool, byte) {
	return mapperlib.GenericCHRRead(&m.chrROM, addr)
}

func (m *MapperUxROM) WriteCHRMemory(addr types.Word, value byte) {
	mapperlib.GenericCHRWrite(&m.chrROM, addr, value)
}

func (m *MapperUxROM) GetMapperNum() byte {
	return mapperNum
}

func (m *MapperUxROM) SetMirroringUpdateCallback(func(nametable.NametableMirroring)) {}

func (m *MapperUxROM) CPUStep() {}

func (m *MapperUxROM) PPUStep() {}
