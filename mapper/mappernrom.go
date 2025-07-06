package mapper

import (
	"github.com/dfirebird/famigom/program"
	. "github.com/dfirebird/famigom/types"
)

const (
	prgRAMSize = 4096

	lowAddr  = 0x6000
	highAddr = 0xFFFF

	lowPrgRamAddr  = lowAddr
	highPrgRamAddr = 0x7FFF

	lowPrgRomAddr  = 0x8000
	highPrgRomAddr = highAddr

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

func CreateMapperNRom(prgRom []byte, chrRom [kib8]byte, nametableMirroring program.NametableArrangement) (*MapperNROM, AddrRange) {
	isPrgRom32kib := len(prgRom) == kib32
	mapper := MapperNROM{
		prgRom:             prgRom,
		isPrgRom32kib:      isPrgRom32kib,
		prgRam:             [4096]byte{},
		chrRom:             chrRom,
		nametableMirroring: nametableMirroring,
	}

	addrRange := AddrRange{
		LowAddr:  lowAddr,
		HighAddr: highAddr,
	}

	return &mapper, addrRange
}

func (m *MapperNROM) ReadMemory(addr Word) byte {
	if lowPrgRamAddr <= addr && addr <= highPrgRamAddr {
		prgRamAddr := addr - lowPrgRamAddr
		return m.prgRam[prgRamAddr]
	} else {
		prgRomAddr := addr - lowPrgRomAddr
		if !m.isPrgRom32kib {
			prgRomAddr = prgRomAddr % kib16
		}

		return m.prgRom[prgRomAddr]
	}
}

func (m *MapperNROM) WriteMemory(addr Word, value byte) {
	if lowPrgRamAddr <= addr && addr <= highPrgRamAddr {
		prgRamAddr := addr - lowPrgRamAddr
		m.prgRam[prgRamAddr] = value
	}
}
