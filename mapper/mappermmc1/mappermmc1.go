package mappermmc1

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/mapper/mapperlib"
	"github.com/dfirebird/famigom/ppu/nametable"
	"github.com/dfirebird/famigom/types"
)

const (
	mapperNum = 1

	clearBitMask    = 0x80
	registerSelMask = 0x6000
)

type mapperShiftRegister struct {
	noOfWrites byte
	value      byte
}

type prgROMMode uint
type chrROMMode uint

const (
	switchAsWhole prgROMMode = iota
	switchAsWhole_
	fixFirstBank
	fixLastBank
)

const (
	kib8Mode chrROMMode = iota
	kib4Mode
)

const (
	controlRegister = iota
	chrBank0Register
	chrBank1Register
	prgBankRegister
)

type MapperMMC1 struct {
	prgROM []byte
	chrROM []byte
	prgRAM []byte

	prgBanks byte
	chrBanks byte

	writeCounter      uint
	interfaceRegister mapperShiftRegister

	prgMode prgROMMode
	chrMode chrROMMode

	chrBank0Num byte
	chrBank1Num byte

	prgRAMBank []byte

	prgROMBank0 []byte
	prgROMBank1 []byte

	chrROMBank0 []byte
	chrROMBank1 []byte

	updateMirroringCallback *func(nametable.NametableMirroring)
}

func CreateMapperMMC1(prgROM, chrROM []byte, prgBanks, chrBanks byte) *MapperMMC1 {
	var chr []byte
	if len(chrROM) == 0 {
		chrRAM := [constants.Kib8]byte{}
		chr = chrRAM[:]
		chrBanks = 2
	} else {
		chr = chrROM
		chrBanks = chrBanks * 2
	}

	prgRAM := [constants.Kib8]byte{}
	interfaceRegister := mapperShiftRegister{}
	mapper := MapperMMC1{
		prgROM:                  prgROM,
		chrROM:                  chr,
		prgRAM:                  prgRAM[:],
		prgBanks:                prgBanks,
		chrBanks:                chrBanks,
		chrBank0Num:             0,
		chrBank1Num:             0,
		interfaceRegister:       interfaceRegister,
		prgRAMBank:              prgRAM[:],
		updateMirroringCallback: nil,

		prgMode: fixLastBank,
	}

	mapper.FixFirstBank()
	mapper.FixLastBank()
	mapper.BankCHRKib4(0, &mapper.chrROMBank0)
	mapper.BankCHRKib4(0, &mapper.chrROMBank1)

	return &mapper
}

func (m *MapperMMC1) ReadMemory(addr types.Word) (bool, byte) {
	if mapperlib.IsPRGRAMAddr(addr) {
		idx := mapperlib.CalculatePRGRAMAddr(addr)
		return true, m.prgRAMBank[idx]
	}
	if mapperlib.IsPRGROMAddr(addr) {
		idx := mapperlib.CalculatePRGROMAddr(addr)
		// idx => 0x0000..=0x7FFF
		if idx < constants.Kib16 {
			return true, m.prgROMBank0[idx]
		} else {
			return true, m.prgROMBank1[idx-constants.Kib16]
		}
	}
	return false, 0
}

func (m *MapperMMC1) WriteMemory(addr types.Word, value byte) {
	if mapperlib.IsPRGRAMAddr(addr) {
		idx := mapperlib.CalculatePRGRAMAddr(addr)
		m.prgRAMBank[idx] = value
	}
	if mapperlib.IsPRGROMAddr(addr) {
		if value&clearBitMask == 0 {
			if m.writeCounter == 0 {
				is5thWrite, value := m.interfaceRegister.Input(value & 0x01)
				if is5thWrite {
					registerSel := (addr & registerSelMask) >> 13
					switch registerSel {
					case controlRegister:
						m.HandleControlRegister(value)
					case chrBank0Register:
						m.HandleCHRBank0Register(value)
					case chrBank1Register:
						m.HandleCHRBank1Register(value)
					case prgBankRegister:
						m.HandlePRGBankRegister(value)
					}
				}
			}
		} else {
			m.prgMode = fixLastBank
			m.FixLastBank()
			m.interfaceRegister.Clear()
		}
		m.writeCounter += 1
	}
}

func (m *MapperMMC1) ReadCHRMemory(addr types.Word) (bool, byte) {
	if mapperlib.IsCHRROMAddr(addr) {
		// addr => 0x0000..=0x1FFF
		if addr < constants.Kib4 {
			return true, m.chrROMBank0[addr]
		} else {
			return true, m.chrROMBank1[addr-constants.Kib4]
		}
	}
	return false, 0
}

func (m *MapperMMC1) WriteCHRMemory(addr types.Word, value byte) {
	if mapperlib.IsCHRROMAddr(addr) {
		// addr => 0x0000..=0x1FFF
		if addr < constants.Kib4 {
			m.chrROMBank0[addr] = value
		} else {
			m.chrROMBank1[addr-constants.Kib4] = value
		}
	}
}

func (m *MapperMMC1) GetMapperNum() byte {
	return mapperNum
}

func (m *MapperMMC1) SetMirroringUpdateCallback(callback func(nametable.NametableMirroring)) {
	m.updateMirroringCallback = &callback
}

func (m *MapperMMC1) CPUStep() {
	m.writeCounter = 0
}

func (m *MapperMMC1) PPUStep() {}

func (m *MapperMMC1) HandleControlRegister(value byte) {
	nametableSelect := value & 0x3
	switch callback := *(m.updateMirroringCallback); nametableSelect {
	case 0:
		callback(nametable.SingleScreenLo)
	case 1:
		callback(nametable.SingleScreenHi)
	case 2:
		callback(nametable.Vertical)
	case 3:
		callback(nametable.Horizontal)
	}

	prgMode := (value >> 2) & 0x3
	chrMode := (value >> 4) & 0x1
	isCHRModeChanged := m.chrMode != chrROMMode(chrMode)
	m.prgMode = prgROMMode(prgMode)
	m.chrMode = chrROMMode(chrMode)

	if isCHRModeChanged {
		m.HandleCHRBank0Register(m.chrBank0Num)
		m.HandleCHRBank1Register(m.chrBank1Num)
	}
}

func (m *MapperMMC1) HandleCHRBank0Register(value byte) {
	m.chrBank0Num = value
	if m.chrMode == kib8Mode {
		bankSel0 := (value &^ 0x1)
		bankSel1 := (value &^ 0x1) | 0x1
		m.BankCHRKib4(byte(bankSel0), &m.chrROMBank0)
		m.BankCHRKib4(byte(bankSel1), &m.chrROMBank1)
	} else {
		m.BankCHRKib4(value, &m.chrROMBank0)
	}
}

func (m *MapperMMC1) HandleCHRBank1Register(value byte) {
	m.chrBank1Num = value
	if m.chrMode == kib4Mode {
		m.BankCHRKib4(value, &m.chrROMBank1)
	}
}

func (m *MapperMMC1) HandlePRGBankRegister(value byte) {
	value = value &^ 0x10
	switch m.prgMode {
	case switchAsWhole, switchAsWhole_:
		bankSel0 := (value &^ 0x1)
		bankSel1 := (value &^ 0x1) | 0x1
		m.BankPRGKib16(bankSel0, &m.prgROMBank0)
		m.BankPRGKib16(bankSel1, &m.prgROMBank1)
	case fixFirstBank:
		m.FixFirstBank()
		m.BankPRGKib16(value, &m.prgROMBank1)
	case fixLastBank:
		m.FixLastBank()
		m.BankPRGKib16(value, &m.prgROMBank0)
	}
}

func (m *MapperMMC1) BankPRGKib16(value byte, bank *[]byte) {
	bankStartIdx, bankEndIdx := mapperlib.CalculateBankStartEnd(value, m.prgBanks, constants.Kib16)
	*bank = m.prgROM[bankStartIdx:bankEndIdx]
}

func (m *MapperMMC1) BankCHRKib4(value byte, bank *[]byte) {
	bankStartIdx, bankEndIdx := mapperlib.CalculateBankStartEnd(value, m.chrBanks, constants.Kib4)
	*bank = m.chrROM[bankStartIdx:bankEndIdx]
}

func (m *MapperMMC1) FixFirstBank() {
	m.prgROMBank0 = m.prgROM[0:constants.Kib16]
}

func (m *MapperMMC1) FixLastBank() {
	startIdx := len(m.prgROM) - constants.Kib16
	m.prgROMBank1 = m.prgROM[startIdx:]
}

func (sr *mapperShiftRegister) Input(bit byte) (bool, byte) {
	sr.noOfWrites += 1
	sr.value = (sr.value >> 1) | (bit << 4)

	if sr.noOfWrites == 5 {
		defer sr.Clear()
	}

	return sr.noOfWrites == 5, sr.value
}

func (sr *mapperShiftRegister) Clear() {
	sr.noOfWrites = 0
	sr.value = 0
}
