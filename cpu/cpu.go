package cpu

const maxMemory = 2048

type word = uint16

type CPU struct {
	// General purpose registers
	X byte
	Y byte
	A byte

	// Special Registers
	SR status
	SP byte
	PC word

	// TODO: should have references to PPU and APU as well here

	Memory [maxMemory]byte
}

func New() CPU {
	return CPU{
		X: 0x0,
		Y: 0x0,
		A: 0x0,
	}
}
