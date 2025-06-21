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
	setResultFactory := func(c *CPU) func(byte) {
		return func(v byte) { c.writeMemory(c.currentGetAddr, v) }
	}

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

	case 0x81:
		c.STA(c.getWithXIndexIndirectAddr())
	case 0x91:
		c.STA(c.getWithIndirectYIndexAddr())
	case 0x85:
		c.STA(c.getWithZeroPageAddress())
	case 0x95:
		c.STA(c.getWithZeroPageIndexedAddr(c.X))
	case 0x99:
		c.STA(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x8D:
		c.STA(c.getWithAbsoluteAddress())
	case 0x9D:
		c.STA(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x86:
		c.STX(c.getWithZeroPageAddress())
	case 0x96:
		c.STX(c.getWithZeroPageIndexedAddr(c.Y))
	case 0x8E:
		c.STX(c.getWithAbsoluteAddress())

	case 0x84:
		c.STY(c.getWithZeroPageAddress())
	case 0x94:
		c.STY(c.getWithZeroPageIndexedAddr(c.X))
	case 0x8C:
		c.STY(c.getWithAbsoluteAddress())

	case 0xA1:
		c.LDA(c.getWithXIndexIndirectAddr())
	case 0xB1:
		c.LDA(c.getWithIndirectYIndexAddr())
	case 0xA5:
		c.LDA(c.getWithZeroPageAddress())
	case 0xB5:
		c.LDA(c.getWithZeroPageIndexedAddr(c.X))
	case 0xA9:
		c.LDA(c.getWithImmediate())
	case 0xB9:
		c.LDA(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0xAD:
		c.LDA(c.getWithAbsoluteAddress())
	case 0xBD:
		c.LDA(c.getWithAbsoluteIndexedAddr(c.X))

	case 0xA2:
		c.LDX(c.getWithImmediate())
	case 0xA6:
		c.LDX(c.getWithZeroPageAddress())
	case 0xB6:
		c.LDX(c.getWithZeroPageIndexedAddr(c.Y))
	case 0xAE:
		c.LDX(c.getWithAbsoluteAddress())
	case 0xBE:
		c.LDX(c.getWithAbsoluteIndexedAddr(c.Y))

	case 0xA0:
		c.LDY(c.getWithImmediate())
	case 0xA4:
		c.LDX(c.getWithZeroPageAddress())
	case 0xB4:
		c.LDY(c.getWithZeroPageIndexedAddr(c.Y))
	case 0xAC:
		c.LDY(c.getWithAbsoluteAddress())
	case 0xBC:
		c.LDY(c.getWithAbsoluteIndexedAddr(c.Y))

	case 0xAA:
		c.TAX()
	case 0xA8:
		c.TAY()
	case 0x8A:
		c.TXA()
	case 0x98:
		c.TYA()

	case 0x06:
		c.ASL(c.getWithZeroPageAddress(), setResultFactory(c))
	case 0x16:
		c.ASL(c.getWithZeroPageIndexedAddr(c.X), setResultFactory(c))
	case 0x0A:
		c.ASL(c.A, func(v byte) { c.A = v })
	case 0x0E:
		c.ASL(c.getWithAbsoluteAddress(), setResultFactory(c))
	case 0x1E:
		c.ASL(c.getWithAbsoluteIndexedAddr(c.X), setResultFactory(c))

	case 0x46:
		c.LSR(c.getWithZeroPageAddress(), setResultFactory(c))
	case 0x56:
		c.LSR(c.getWithZeroPageIndexedAddr(c.X), setResultFactory(c))
	case 0x4A:
		c.LSR(c.A, func(v byte) { c.A = v })
	case 0x4E:
		c.LSR(c.getWithAbsoluteAddress(), setResultFactory(c))
	case 0x5E:
		c.LSR(c.getWithAbsoluteIndexedAddr(c.X), setResultFactory(c))

	case 0x26:
		c.ROL(c.getWithZeroPageAddress(), setResultFactory(c))
	case 0x36:
		c.ROL(c.getWithZeroPageIndexedAddr(c.X), setResultFactory(c))
	case 0x2A:
		c.ROL(c.A, func(v byte) { c.A = v })
	case 0x2E:
		c.ROL(c.getWithAbsoluteAddress(), setResultFactory(c))
	case 0x3E:
		c.ROL(c.getWithAbsoluteIndexedAddr(c.X), setResultFactory(c))

	case 0x66:
		c.ROR(c.getWithZeroPageAddress(), setResultFactory(c))
	case 0x76:
		c.ROR(c.getWithZeroPageIndexedAddr(c.X), setResultFactory(c))
	case 0x6A:
		c.ROR(c.A, func(v byte) { c.A = v })
	case 0x6E:
		c.ROR(c.getWithAbsoluteAddress(), setResultFactory(c))
	case 0x7E:
		c.ROR(c.getWithAbsoluteIndexedAddr(c.X), setResultFactory(c))

	case 0xC1:
		c.CMP(c.getWithXIndexIndirectAddr())
	case 0xD1:
		c.CMP(c.getWithIndirectYIndexAddr())
	case 0xC5:
		c.CMP(c.getWithZeroPageAddress())
	case 0xD5:
		c.CMP(c.getWithZeroPageIndexedAddr(c.X))
	case 0xC9:
		c.CMP(c.getWithImmediate())
	case 0xD9:
		c.CMP(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0xCD:
		c.CMP(c.getWithAbsoluteAddress())
	case 0xDD:
		c.CMP(c.getWithAbsoluteIndexedAddr(c.X))

	case 0xC0:
		c.CPY(c.getWithImmediate())
	case 0xC4:
		c.CPY(c.getWithZeroPageAddress())
	case 0xCC:
		c.CPY(c.getWithAbsoluteAddress())

	case 0xE0:
		c.CPX(c.getWithImmediate())
	case 0xE4:
		c.CPX(c.getWithZeroPageAddress())
	case 0xEC:
		c.CPX(c.getWithAbsoluteAddress())
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
