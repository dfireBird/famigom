package program

import "bytes"

type NametableArrangement int

const (
	PRG_ROM_BANK_UNIT_SIZE = 16384
	CHR_ROM_BANK_UNIT_SIZE = 8192

	PLAYCHOICE_INST_ROM_SIZE = 8192
	PLAYCHOICE_PROM_SIZE     = 16
)

var NES_HEADER []byte = []byte{0x4E, 0x45, 0x53, 0x1A}

const (
	Vertical NametableArrangement = iota
	Horizontal
)

type Program struct {
	Mapper         byte
	PrgRomBankSize byte
	ChrRomBankSize byte

	PrgRom []byte
	ChrRom []byte

	// Not used for now but still parsed
	nametableArrangement         NametableArrangement
	isbatteryBackedRam           bool
	isTrainerPresent             bool
	isAlternativeNameTableLayout bool
	isVSUnisystem                bool

	isPlaychoice      bool
	playchoiceInstRom []byte
	playchoiceProm    []byte
}

func Parse(fileData []byte) (bool, *Program) {
	seekIdx := uint(0)

	if nesHeader := fileData[:4]; !bytes.Equal(nesHeader, NES_HEADER) {
		return false, nil
	}
	seekIdx += 4

	PrgRomBankSize := fileData[seekIdx]
	seekIdx += 1
	ChrRomBankSize := fileData[seekIdx]
	seekIdx += 1

	flags6 := fileData[seekIdx]
	seekIdx += 1
	flags7 := fileData[seekIdx]
	seekIdx += 1

	// skipping bytes from 9 - 16
	seekIdx += 8

	if nes2FormatFlag := (flags7 & 0x0C) >> 2; nes2FormatFlag == 2 {
		return false, nil
	}

	isTrainerPresent := (flags6 & 0x04) == 0x04
	if isTrainerPresent {
		// we are ignoring trainer data for now
		// only seeking the cursor
		seekIdx += 512
	}

	prgRomSize := (PRG_ROM_BANK_UNIT_SIZE * uint(PrgRomBankSize))
	PrgRom := fileData[seekIdx : seekIdx+prgRomSize]
	seekIdx += prgRomSize

	chrRomSize := (CHR_ROM_BANK_UNIT_SIZE * uint(ChrRomBankSize))
	ChrRom := fileData[seekIdx : seekIdx+chrRomSize]
	seekIdx += chrRomSize

	isPlaychoice := (flags7 & 0x02) == 0x02
	var playchoiceProm []byte
	var playchoiceInstRom []byte
	if isPlaychoice {
		playchoiceInstRom = fileData[seekIdx : seekIdx+PLAYCHOICE_INST_ROM_SIZE]
		seekIdx += PLAYCHOICE_INST_ROM_SIZE

		playchoiceProm = fileData[seekIdx : seekIdx+PLAYCHOICE_PROM_SIZE]
		seekIdx += PLAYCHOICE_PROM_SIZE
	}

	mapperLo := flags6 & 0xF0
	mapperHi := flags7 & 0xF0

	Mapper := mapperHi | (mapperLo >> 4)

	nametableArrangement := NametableArrangement(flags6 & 0x01)
	isbatteryBackedRam := (flags6 & 0x02) == 0x02
	isAlternativeNameTableLayout := (flags6 & 0x08) == 0x08

	isVSUnisystem := (flags7 & 0x01) == 0x01

	program := Program {
		Mapper:                       Mapper,
		PrgRomBankSize:               PrgRomBankSize,
		ChrRomBankSize:               ChrRomBankSize,
		PrgRom:                       PrgRom,
		ChrRom:                       ChrRom,
		nametableArrangement:         nametableArrangement,
		isbatteryBackedRam:           isbatteryBackedRam,
		isTrainerPresent:             isTrainerPresent,
		isAlternativeNameTableLayout: isAlternativeNameTableLayout,
		isVSUnisystem:                isVSUnisystem,
		isPlaychoice:                 isPlaychoice,
		playchoiceInstRom:            playchoiceInstRom,
		playchoiceProm:               playchoiceProm,
	}

	return true, &program
}
