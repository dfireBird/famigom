package program

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dfirebird/famigom/log"
	"github.com/klauspost/compress/zip"
)

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
	ErrInvalidNesRom      = fmt.Errorf("ROM file is invalid/corrupted")
	ErrInsupportedVersion = fmt.Errorf("NES 2.0 ROM is not supported yet")
	ErrNotSingleFileZip   = fmt.Errorf("input zip archive has more than one files")
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
	logger := log.Logger()
	seekIdx := uint(0)

	if unzippedData, err := unzipIfPossible(romData); err != nil {
		logger.Infoln("Unzipping file met with an error:", err.Error())
		logger.Infoln("Considering the input file as RAW NES file.")
	} else {
		romData = unzippedData
	}

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

	/* Temoprarily commenting it out, since without this error some ines 2.0 header ROMs do work correctly */
	// if nes2FormatFlag := (flags7 & 0x0C) >> 2; nes2FormatFlag == 2 {
	// 	return nil, ErrInsupportedVersion
	// }

	mapperLo := flags6 & 0xF0
	mapperHi := flags7 & 0xF0

	Mapper := mapperHi | (mapperLo >> 4)

	nametableArrangement := NametableArrangement(flags6 & 0x01)

	logger.Infof("Reading Header")
	logger.Infof("PRG ROM Bank Size %d", PrgRomBankSize)
	logger.Infof("CHR ROM Bank Size %d", ChrRomBankSize)
	logger.Infof("Mapper #%d", Mapper)
	logger.Infof("Name Table Mirroring %v", nametableArrangement.getMirroring())

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

	logger.Infof("Read %d bytes of data", seekIdx)
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

func unzipIfPossible(romData []byte) ([]byte, error) {
	reader := bytes.NewReader(romData)
	zipReader, err := zip.NewReader(reader, int64(len(romData)))
	if err != nil {
		return nil, err
	}

	if len(zipReader.File) != 1 {
		return nil, ErrNotSingleFileZip
	}

	name := zipReader.File[0].FileInfo().Name()
	f, err := zipReader.Open(name)
	if err != nil {
		return nil, err
	}

	unzipped, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return unzipped, nil
}

func (n NametableArrangement) getMirroring() NametableArrangement {
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
