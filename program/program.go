package program

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

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

var (
	ErrInvalidNesRom = fmt.Errorf("ROM file is invalid/corrupted")
	ErrInsupportedVersion = fmt.Errorf("emulator does not support NES 2.0 ROM yet")
)

type Program struct {
	Mapper         byte
	PrgRomBankSize byte
	ChrRomBankSize byte

	NametableArrangement NametableArrangement

	PrgRom []byte
	ChrRom []byte
}

func Parse(romData []byte) (*Program, error) {
	seekIdx := uint(0)

	if nesHeader := romData[:4]; !bytes.Equal(nesHeader, NES_HEADER) {
		return nil, ErrInvalidNesRom
	}
	seekIdx += 4

	PrgRomBankSize := romData[seekIdx]
	seekIdx += 1
	ChrRomBankSize := romData[seekIdx]
	seekIdx += 1

	flags6 := romData[seekIdx]
	seekIdx += 1
	flags7 := romData[seekIdx]
	seekIdx += 1

	// skipping bytes from 9 - 16
	seekIdx += 8

	if nes2FormatFlag := (flags7 & 0x0C) >> 2; nes2FormatFlag == 2 {
		return nil, ErrInsupportedVersion
	}

	mapperLo := flags6 & 0xF0
	mapperHi := flags7 & 0xF0

	Mapper := mapperHi | (mapperLo >> 4)

	nametableArrangement := NametableArrangement(flags6 & 0x01)

	logger.Printf("Reading Header")
	logger.Printf("PRG ROM Bank Size %d", PrgRomBankSize)
	logger.Printf("CHR ROM Bank Size %d", ChrRomBankSize)
	logger.Printf("Mapper #%d", Mapper)
	logger.Printf("Name Table Mirroring %v", nametableArrangement.GetMirroring())

	isTrainerPresent := (flags6 & 0x04) == 0x04
	if isTrainerPresent {
		// we are ignoring trainer data for now
		// only seeking the cursor
		seekIdx += 512
	}

	prgRomSize := (PRG_ROM_BANK_UNIT_SIZE * uint(PrgRomBankSize))
	PrgRom := romData[seekIdx : seekIdx+prgRomSize]
	seekIdx += prgRomSize

	chrRomSize := (CHR_ROM_BANK_UNIT_SIZE * uint(ChrRomBankSize))
	ChrRom := romData[seekIdx : seekIdx+chrRomSize]
	seekIdx += chrRomSize

	isPlaychoice := (flags7 & 0x02) == 0x02
	if isPlaychoice {
		seekIdx += PLAYCHOICE_INST_ROM_SIZE
		seekIdx += PLAYCHOICE_PROM_SIZE
	}

	logger.Printf("Read %d bytes of data", seekIdx)
	program := Program{
		Mapper:               Mapper,
		PrgRomBankSize:       PrgRomBankSize,
		ChrRomBankSize:       ChrRomBankSize,
		PrgRom:               PrgRom,
		ChrRom:               ChrRom,
		NametableArrangement: nametableArrangement,
	}

	return &program, nil
}

func (n NametableArrangement) GetMirroring() NametableArrangement {
    if n == Horizontal {
		return Vertical
	} else {
		return Horizontal
	}
}

func (n NametableArrangement) String() string {
	if n == Horizontal {
		return "Horizontal"
	} else {
		return "Vertical"
	}
}
