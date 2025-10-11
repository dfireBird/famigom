package mapperaxrom

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 0x07
)

type MapperAxROM struct {
	prgROM []byte
	chrROM [constants.Kib8]byte

	maxBanks                byte
	currentBank             []byte
	updateMirroringCallback *func(nametable.NametableMirroring)
}

func CreateMapperAxROM(prgROM, chrROM []byte) *MapperAxROM {
	var chr [constants.Kib8]byte
	if len(chrROM) == 0 {
		chr = [constants.Kib8]byte{} // CHR RAM
	} else {
		chr = [constants.Kib8]byte(chrROM)
	}

	return &MapperAxROM{
		prgROM:                  prgROM,
		chrROM:                  chr,
		maxBanks:                byte(len(prgROM) / constants.Kib32),
		currentBank:             prgROM[:constants.Kib32],
		updateMirroringCallback: nil,
	}
}

func (m *MapperAxROM) ReadMemory(addr types.Word) (bool, byte) {
	if constants.LowPrgROMAddr <= addr && addr <= constants.HighPrgROMAddr {
		idx := addr - constants.LowPrgROMAddr
		return true, m.currentBank[idx]
	}
	return false, 0
}

func (m *MapperAxROM) WriteMemory(addr types.Word, value byte) {
	if constants.LowPrgROMAddr <= addr && addr <= constants.HighPrgROMAddr {
		bankSel := (value & 0x07) % m.maxBanks
		bankStartIdx := uint(bankSel) * constants.Kib32
		bankEndIdx := bankStartIdx + constants.Kib32
		m.currentBank = m.prgROM[bankStartIdx:bankEndIdx]

		if m.updateMirroringCallback == nil {
			// Should be an unreachable condition
			panic("Mirroring callback is not set before called")
		}

		mirroringCallback := *m.updateMirroringCallback
		if nametableSel := (value & 0x10); nametableSel == 0 {
			mirroringCallback(nametable.SingleScreenLo)
		} else {
			mirroringCallback(nametable.SingleScreenHi)
		}
	}
}

func (m *MapperAxROM) ReadCHRMemory(addr types.Word) (bool, byte) {
	if constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr {
		return true, m.chrROM[addr]
	}
	return false, 0
}

func (m *MapperAxROM) WriteCHRMemory(addr types.Word, value byte) {
	if constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr {
		m.chrROM[addr] = value
	}
}

func (m *MapperAxROM) GetMapperNum() byte {
	return mapperNum
}

func (m *MapperAxROM) SetMirroringUpdateCallback(callback func(nametable.NametableMirroring)) {
	m.updateMirroringCallback = &callback
}
