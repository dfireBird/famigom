package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/dfirebird/famigom/program"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func main() {
	verbose := flag.Bool("v", false, "Enables verbose logging")
	flag.Parse()

	if len(flag.Args()) != 1 {
		logger.Printf("ERROR: NES ROM file path is not passed as an argument")
		os.Exit(1)
	}

	romPath, _ := filepath.Abs(flag.Arg(0))
	logger.Printf("Reading ROM/Program file from path: %s", romPath)
	romData, err := os.ReadFile(romPath)
	if err != nil {
		logger.Printf("ERROR: Reading the ROM file %v", err)
		os.Exit(1)
	}

	if *verbose {
		logger.Printf("Parsing ROM/Program file of size %d", len(romData))
	}
	program.Parse(romData)
}
