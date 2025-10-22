package mappernrom

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/log"
	"github.com/dfirebird/famigom/mapper/mapperlib"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/program"
	. "github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x00

	prgRAMSize = constants.Kib4
)

var (
	WarnRequirePRGRAM = "NROM requires PRG RAM. iNES header field for PRG RAM is set 0. Assuming 4KiB of PRG RAM."
)

type MapperNROM struct {
	prgRom        []byte
	isPrgRom32kib bool

	prgRAM [prgRAMSize]byte

	chrRom [constants.Kib8]byte
}

func CreateMapperNRom(program *program.Program) *MapperNROM {
	var chr [constants.Kib8]byte

	// a prg bank has 16Kib => bank size of 2 == 32 kib
	isPrgRom32kib := program.PrgRomBankSize == 2

	if program.IsINES2 {
		// should and will panic if chr ram size is 0 and chr rom length is also 0
		chr = [constants.Kib8]byte(program.ChrRom)
		if program.ChrRAMSize > 0 {
			chr = [constants.Kib8]byte{}
		}

		if program.PrgRAMSize > 0 {
			log.Logger().Warnln(WarnRequirePRGRAM)
		}
	} else {
		if program.ChrRomBankSize == 0 {
			chr = [constants.Kib8]byte{} // CHR RAM
		} else {
			chr = [constants.Kib8]byte(program.ChrRom)
		}
	}

	mapper := MapperNROM{
		prgRom:        program.PrgRom,
		isPrgRom32kib: isPrgRom32kib,
		prgRAM:        [prgRAMSize]byte{},
		chrRom:        chr,
	}

	return &mapper
}

func (m *MapperNROM) ReadMemory(addr Word) (bool, byte) {
	if mapperlib.IsPRGRAMAddr(addr) {
		prgRAMAddr := mapperlib.CalculatePRGRAMAddr(addr)
		return true, m.prgRAM[prgRAMAddr]
	} else if mapperlib.IsPRGROMAddr(addr) {
		prgRomAddr := mapperlib.CalculatePRGROMAddr(addr)
		if !m.isPrgRom32kib {
			prgRomAddr = prgRomAddr % constants.Kib16
		}

		return true, m.prgRom[prgRomAddr]
	} else {
		return false, 0
	}
}

func (m *MapperNROM) WriteMemory(addr Word, value byte) {
	if mapperlib.IsPRGRAMAddr(addr) {
		prgRAMAddr := mapperlib.CalculatePRGRAMAddr(addr)
		m.prgRAM[prgRAMAddr] = value
	}
}

func (m *MapperNROM) ReadCHRMemory(addr Word) (bool, byte) {
	return mapperlib.GenericCHRRead(&m.chrRom, addr)
}

func (m *MapperNROM) WriteCHRMemory(addr Word, value byte) {
	mapperlib.GenericCHRWrite(&m.chrRom, addr, value)
}

func (m *MapperNROM) GetMapperNum() byte {
	return mapperNum
}

func (m *MapperNROM) SetMirroringUpdateCallback(func(nametable.NametableMirroring)) {}

func (m *MapperNROM) CPUStep() {}

func (m *MapperNROM) PPUStep() {}
