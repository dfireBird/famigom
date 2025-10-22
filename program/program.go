package program

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dfirebird/famigom/log"
	"github.com/dfirebird/famigom/types"
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

var (
	warnHasTrainer           = "ROM file has trainer data. Trainer support is not implemented yet. Running the ROM may have issues."
	warnHasNonVolatileMemory = "ROM file has non-volatile memory flag set. Non-volatile memory support is not implemented yet. Running the ROM will not have any issues but the RAM contents will not be saved."

	warnNonNESFamicom = "ROM file indicates the console type is not NES/Famicom. Support for other types of consoles is not implemented. Running the ROM may cause panics or have issues."

	warnTimingType = "ROM file indicates timing system of %s. Support has not been not implemented yet. Running the ROM may have issues."
)

type Program struct {
	Mapper         types.Word
	PrgRomBankSize types.Word
	ChrRomBankSize types.Word

	NametableArrangement NametableArrangement

	PrgRom []byte
	ChrRom []byte

	IsINES2 bool

	// common unused fields
	hasNonVolatileMemory   bool
	hasTrainer             bool
	isAlternativeNametable bool

	// iNES 2.0 specific fields
	SubMapperNumber byte

	PrgRAMSize   types.Word
	PrgNVRAMSize types.Word

	ChrRAMSize   types.Word
	ChrNVRAMSize types.Word
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

	headers := [12]byte{}
	for seekIdx < 16 {
		headers[seekIdx-4] = romData[seekIdx]
		seekIdx++
	}

	if isINES2 := headers[3]&0x0C == 0x08; isINES2 {
		header6 := headers[2]
		NametableArrangement := NametableArrangement(header6 & 0x01)
		hasNonVolatileMemory := header6&0x02 == 0x02
		hasTrainer := header6&0x04 == 0x04
		isAlternativeNametable := header6&0x08 == 0x08

		if hasTrainer {
			logger.Warnln(warnHasTrainer)
			seekIdx += 512
		}
		if hasNonVolatileMemory {
			logger.Warn(warnHasNonVolatileMemory)
		}

		header7 := headers[3]
		consoleType := header7 & 0x03
		if consoleType != 0 {
			logger.Warnln(warnNonNESFamicom)
		}

		header8 := headers[4]
		SubMapperNumber := header8 & 0xF0 >> 4

		mapperLSBLo := types.Word(header6 & 0xF0)
		mapperLSBHi := types.Word(header7 & 0xF0)
		mapperMSB := types.Word(header8 & 0x0F)

		Mapper := (mapperMSB << 8) | mapperLSBHi | (mapperLSBLo >> 4)

		prgROMBankSizeLSB := types.Word(headers[0])
		chrROMBankSizeLSB := types.Word(headers[1])

		prgROMBankSizeMSB := types.Word(headers[5]&0x0F) << 8
		chrROMBankSizeMSB := types.Word(headers[5]&0xF0) << 4

		PrgRomBankSize := prgROMBankSizeMSB | prgROMBankSizeLSB
		ChrRomBankSize := chrROMBankSizeMSB | chrROMBankSizeLSB

		var PrgRAMSize, PrgNVRAMSize, ChrRAMSize, ChrNVRAMSize types.Word

		if shiftCount := headers[6] & 0x0F; shiftCount != 0 {
			PrgRAMSize = 64 << shiftCount
		}
		if shiftCount := (headers[6] & 0xF0) >> 4; shiftCount != 0 {
			PrgNVRAMSize = 64 << shiftCount
		}

		if shiftCount := headers[7] & 0x0F; shiftCount != 0 {
			ChrRAMSize = 64 << shiftCount
		}
		if shiftCount := (headers[7] & 0xF0) >> 4; shiftCount != 0 {
			ChrNVRAMSize = 64 << shiftCount
		}

		systemTimingType := headers[8] & 0x03
		switch systemTimingType {
		case 0x01:
			logger.Warnf(warnTimingType, "PAL")
		case 0x03:
			logger.Warnf(warnTimingType, "Dendy")
		}

		prgRomSize := (PRG_ROM_BANK_UNIT_SIZE * uint(PrgRomBankSize))
		PrgRom := romData[seekIdx : seekIdx+prgRomSize]
		seekIdx += prgRomSize

		chrRomSize := (CHR_ROM_BANK_UNIT_SIZE * uint(ChrRomBankSize))
		ChrRom := romData[seekIdx : seekIdx+chrRomSize]
		seekIdx += chrRomSize

		logger.Infof("Reading Header")
		logger.Infof("PRG ROM Bank Size %d", PrgRomBankSize)
		logger.Infof("CHR ROM Bank Size %d", ChrRomBankSize)
		logger.Infof("Mapper #%d", Mapper)
		logger.Infof("Name Table Mirroring %v", NametableArrangement.getMirroring())

		program := Program{
			IsINES2:         true,
			Mapper:          Mapper,
			SubMapperNumber: SubMapperNumber,

			NametableArrangement: NametableArrangement,

			PrgRomBankSize: PrgRomBankSize,
			ChrRomBankSize: ChrRomBankSize,
			PrgRom:         PrgRom,
			ChrRom:         ChrRom,

			PrgRAMSize:   PrgRAMSize,
			PrgNVRAMSize: PrgNVRAMSize,
			ChrRAMSize:   ChrRAMSize,
			ChrNVRAMSize: ChrNVRAMSize,

			hasNonVolatileMemory:   hasNonVolatileMemory,
			hasTrainer:             hasTrainer,
			isAlternativeNametable: isAlternativeNametable,
		}

		return &program, nil
	} else {
		logger.Warnln("ROM file has iNES 1.0 format. Some ROMs or features might have issues while running.")
		PrgRomBankSize := headers[0]
		ChrRomBankSize := headers[1]

		flags6 := headers[2]
		flags7 := headers[3]

		hasNonVolatileMemory := flags6&0x02 == 0x02
		hasTrainer := flags6&0x04 == 0x04
		isAlternativeNametable := flags6&0x08 == 0x08

		mapperLo := flags6 & 0xF0
		mapperHi := flags7 & 0xF0

		Mapper := mapperHi | (mapperLo >> 4)

		nametableArrangement := NametableArrangement(flags6 & 0x01)

		if hasTrainer {
			logger.Warnln(warnHasTrainer)
			// we are ignoring trainer data for now
			// only seeking the cursor
			seekIdx += 512
		}
		if hasNonVolatileMemory {
			logger.Warnln(warnHasNonVolatileMemory)
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

		logger.Infof("Reading Header")
		logger.Infof("PRG ROM Bank Size %d", PrgRomBankSize)
		logger.Infof("CHR ROM Bank Size %d", ChrRomBankSize)
		logger.Infof("Mapper #%d", Mapper)
		logger.Infof("Name Table Mirroring %v", nametableArrangement.getMirroring())
		program := Program{
			IsINES2:              false,
			Mapper:               types.Word(Mapper),
			PrgRomBankSize:       types.Word(PrgRomBankSize),
			ChrRomBankSize:       types.Word(ChrRomBankSize),
			PrgRom:               PrgRom,
			ChrRom:               ChrRom,
			NametableArrangement: nametableArrangement,

			hasNonVolatileMemory:   hasNonVolatileMemory,
			hasTrainer:             hasTrainer,
			isAlternativeNametable: isAlternativeNametable,
		}

		return &program, nil
	}
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
