package cpu

const (
	maxMemory = (1 << 16) - 1
)

type word = uint16

type CPU struct {
	// General purpose registers
	X byte
	Y byte
	A byte

	// Special Registers
	Flags status
	SP    byte
	PC    word

	Memory [maxMemory]byte

	// Internal to emulator

	// Only used in GetXXXAddr funtions for memory write instructions
	currentGetAddr word
}

func New() CPU {
	return CPU{
		X: 0x0,
		Y: 0x0,
		A: 0x0,
	}
}

func (c *CPU) Step() {
	opcode := c.readMemory(c.PC)
	c.PC++
	switch opcode {
	case 0x00:
		c.BRK()

	case 0x01:
		c.ORA(c.getWithXIndexIndirectAddr())
	case 0x11:
		c.ORA(c.getWithIndirectYIndexAddr())
	case 0x05:
		c.ORA(c.getWithZeroPageAddress())
	case 0x15:
		c.ORA(c.getWithZeroPageIndexedAddr(c.X))
	case 0x09:
		c.ORA(c.getWithImmediate())
	case 0x19:
		c.ORA(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x0D:
		c.ORA(c.getWithAbsoluteAddress())
	case 0x1D:
		c.ORA(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x21:
		c.AND(c.getWithXIndexIndirectAddr())
	case 0x31:
		c.AND(c.getWithIndirectYIndexAddr())
	case 0x25:
		c.AND(c.getWithZeroPageAddress())
	case 0x35:
		c.AND(c.getWithZeroPageIndexedAddr(c.X))
	case 0x29:
		c.AND(c.getWithImmediate())
	case 0x39:
		c.AND(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x2D:
		c.AND(c.getWithAbsoluteAddress())
	case 0x3D:
		c.AND(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x41:
		c.EOR(c.getWithXIndexIndirectAddr())
	case 0x51:
		c.EOR(c.getWithIndirectYIndexAddr())
	case 0x45:
		c.EOR(c.getWithZeroPageAddress())
	case 0x55:
		c.EOR(c.getWithZeroPageIndexedAddr(c.X))
	case 0x49:
		c.EOR(c.getWithImmediate())
	case 0x59:
		c.EOR(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x4D:
		c.EOR(c.getWithAbsoluteAddress())
	case 0x5D:
		c.EOR(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x24:
		c.BIT(c.getWithZeroPageAddress())
	case 0x2C:
		c.BIT(c.getWithAbsoluteAddress())

	case 0x61:
		c.ADC(c.getWithXIndexIndirectAddr())
	case 0x71:
		c.ADC(c.getWithIndirectYIndexAddr())
	case 0x65:
		c.ADC(c.getWithZeroPageAddress())
	case 0x75:
		c.ADC(c.getWithZeroPageIndexedAddr(c.X))
	case 0x69:
		c.ADC(c.getWithImmediate())
	case 0x79:
		c.ADC(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x6D:
		c.ADC(c.getWithAbsoluteAddress())
	case 0x7D:
		c.ADC(c.getWithAbsoluteIndexedAddr(c.X))

	case 0xE1:
		c.SBC(c.getWithXIndexIndirectAddr())
	case 0xF1:
		c.SBC(c.getWithIndirectYIndexAddr())
	case 0xE5:
		c.SBC(c.getWithZeroPageAddress())
	case 0xF5:
		c.SBC(c.getWithZeroPageIndexedAddr(c.X))
	case 0xE9:
		c.SBC(c.getWithImmediate())
	case 0xF9:
		c.SBC(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0xED:
		c.SBC(c.getWithAbsoluteAddress())
	case 0xFD:
		c.SBC(c.getWithAbsoluteIndexedAddr(c.X))

	case 0xE6:
		c.INC(c.getWithZeroPageAddress(), c.writeMemory)
	case 0xEE:
		c.INC(c.getWithAbsoluteAddress(), c.writeMemory)
	case 0xF6:
		c.INC(c.getWithZeroPageIndexedAddr(c.X), c.writeMemory)
	case 0xFE:
		c.INC(c.getWithAbsoluteIndexedAddr(c.X), c.writeMemory)

	case 0xC6:
		c.DEC(c.getWithZeroPageAddress(), c.writeMemory)
	case 0xCE:
		c.DEC(c.getWithAbsoluteAddress(), c.writeMemory)
	case 0xD6:
		c.DEC(c.getWithZeroPageIndexedAddr(c.X), c.writeMemory)
	case 0xDE:
		c.DEC(c.getWithAbsoluteIndexedAddr(c.X), c.writeMemory)

	case 0xE8:
		c.INX()
	case 0xC8:
		c.INY()

	case 0xCA:
		c.DEX()
	case 0x88:
		c.DEY()
	}
}

func (c *CPU) readMemory(addr word) byte {
	return c.Memory[addr]
}

func (c *CPU) writeMemory(addr word, value byte) {
	c.Memory[addr] = value
}

func (c *CPU) pushIntoStack(value byte) {
	effectiveStackAddr := 0x0100 | word(c.SP)
	c.writeMemory(effectiveStackAddr, value)
	c.SP--
}

func joinBytesToWord(lo, hi byte) word {
	const ADDR_SHIFT = 8
	return word(lo) | word(hi)<<ADDR_SHIFT
}

func splitWordToByte(w word) (byte, byte) {
	const (
		LO_ADDR_MASK = 0x00FF
		HI_ADDR_MASK = 0xFF00
		ADDR_SHIFT   = 8
	)

	return byte(w & LO_ADDR_MASK), byte((w & HI_ADDR_MASK) >> ADDR_SHIFT)
}
