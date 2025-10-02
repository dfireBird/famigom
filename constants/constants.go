package constants

const (
	LowMapperAddr  = 0x6000
	HighMapperAddr = 0xFFFF
	LowPrgRAMAddr  = LowMapperAddr
	HighPrgRAMAddr = 0x7FFF
	LowPrgROMAddr  = 0x8000
	HighPrgROMAddr = HighMapperAddr

	LowChrROMAddr  = 0x0000
	HighChrROMAddr = 0x1FFF

	Kib1  = 1024
	Kib2  = Kib1 * 2
	Kib4  = Kib2 * 2
	Kib8  = Kib4 * 2
	Kib16 = Kib8 * 2
	Kib32 = Kib16 * 2
)
