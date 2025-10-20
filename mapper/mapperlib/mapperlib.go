package mapperlib

import (
	"github.com/dfirebird/famigom/constants"
	"github.com/dfirebird/famigom/types"
)

func IsPRGRAMAddr(addr types.Word) bool {
	return constants.LowPrgRAMAddr <= addr && addr <= constants.HighPrgRAMAddr
}

func IsPRGROMAddr(addr types.Word) bool {
	return constants.LowPrgROMAddr <= addr && addr <= constants.HighPrgROMAddr
}

func CalculatePRGROMAddr(addr types.Word) types.Word {
	return addr - constants.LowPrgROMAddr
}

func CalculatePRGRAMAddr(addr types.Word) types.Word {
	return addr - constants.LowPrgRAMAddr
}

func IsCHRROMAddr(addr types.Word) bool {
	return constants.LowChrROMAddr <= addr && addr <= constants.HighChrROMAddr
}

func GenericCHRRead(chrBank *[constants.Kib8]byte, addr types.Word) (bool, byte) {
	if IsCHRROMAddr(addr) {
		return true, chrBank[addr]
	}
	return false, 0
}

func GenericCHRWrite(chrBank *[constants.Kib8]byte, addr types.Word, value byte) {
	if IsCHRROMAddr(addr) {
		chrBank[addr] = value
	}
}

func CalculateBankStartEnd(bankNum, maxBanks byte, size uint) (uint, uint) {
	bankSel := bankNum % maxBanks
	bankStartIdx := uint(bankSel) * size
	bankEndIdx := bankStartIdx + size

	return bankStartIdx, bankEndIdx
}
