package cpu

import (
	. "github.com/dfirebird/famigom/types"
)

const (
	maxMemory = (1 << 16)
)

type CPU struct {
	// General purpose registers
	X byte
	Y byte
	A byte

	// Special Registers
	Flags Status
	SP    byte
	PC    Word

	Memory [maxMemory]byte

	// Internal to emulator

	// Only used in GetXXXAddr funtions for memory write instructions
	currentGetAddr Word
	isJammed bool
}

func New() CPU {
	return CPU{
		X: 0x0,
		Y: 0x0,
		A: 0x0,

		Flags: Status(INITIAL_STATUS),
		SP: 0xFF,
		PC: 0x00,

		// other fields default value as per go spec is fine
	}
}

func (c *CPU) Step() {
	setResultFactory := func(c *CPU) func(byte) {
		return func(v byte) { c.WriteMemory(c.currentGetAddr, v) }
	}

	opcode := c.ReadMemory(c.PC)
	c.PC++
	switch opcode {
	case 0x00:
		c.i_BRK()
	case 0x40:
		c.i_RTI()

	case 0x01:
		c.i_ORA(c.getWithXIndexIndirectAddr())
	case 0x11:
		c.i_ORA(c.getWithIndirectYIndexAddr())
	case 0x05:
		c.i_ORA(c.getWithZeroPageAddress())
	case 0x15:
		c.i_ORA(c.getWithZeroPageIndexedAddr(c.X))
	case 0x09:
		c.i_ORA(c.getWithImmediate())
	case 0x19:
		c.i_ORA(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x0D:
		c.i_ORA(c.getWithAbsoluteAddress())
	case 0x1D:
		c.i_ORA(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x21:
		c.i_AND(c.getWithXIndexIndirectAddr())
	case 0x31:
		c.i_AND(c.getWithIndirectYIndexAddr())
	case 0x25:
		c.i_AND(c.getWithZeroPageAddress())
	case 0x35:
		c.i_AND(c.getWithZeroPageIndexedAddr(c.X))
	case 0x29:
		c.i_AND(c.getWithImmediate())
	case 0x39:
		c.i_AND(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x2D:
		c.i_AND(c.getWithAbsoluteAddress())
	case 0x3D:
		c.i_AND(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x41:
		c.i_EOR(c.getWithXIndexIndirectAddr())
	case 0x51:
		c.i_EOR(c.getWithIndirectYIndexAddr())
	case 0x45:
		c.i_EOR(c.getWithZeroPageAddress())
	case 0x55:
		c.i_EOR(c.getWithZeroPageIndexedAddr(c.X))
	case 0x49:
		c.i_EOR(c.getWithImmediate())
	case 0x59:
		c.i_EOR(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x4D:
		c.i_EOR(c.getWithAbsoluteAddress())
	case 0x5D:
		c.i_EOR(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x24:
		c.i_BIT(c.getWithZeroPageAddress())
	case 0x2C:
		c.i_BIT(c.getWithAbsoluteAddress())

	case 0x61:
		c.i_ADC(c.getWithXIndexIndirectAddr())
	case 0x71:
		c.i_ADC(c.getWithIndirectYIndexAddr())
	case 0x65:
		c.i_ADC(c.getWithZeroPageAddress())
	case 0x75:
		c.i_ADC(c.getWithZeroPageIndexedAddr(c.X))
	case 0x69:
		c.i_ADC(c.getWithImmediate())
	case 0x79:
		c.i_ADC(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x6D:
		c.i_ADC(c.getWithAbsoluteAddress())
	case 0x7D:
		c.i_ADC(c.getWithAbsoluteIndexedAddr(c.X))

	case 0xE1:
		c.i_SBC(c.getWithXIndexIndirectAddr())
	case 0xF1:
		c.i_SBC(c.getWithIndirectYIndexAddr())
	case 0xE5:
		c.i_SBC(c.getWithZeroPageAddress())
	case 0xF5:
		c.i_SBC(c.getWithZeroPageIndexedAddr(c.X))
	case 0xE9:
		c.i_SBC(c.getWithImmediate())
	case 0xF9:
		c.i_SBC(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0xED:
		c.i_SBC(c.getWithAbsoluteAddress())
	case 0xFD:
		c.i_SBC(c.getWithAbsoluteIndexedAddr(c.X))

	case 0xE6:
		c.i_INC(c.getWithZeroPageAddress(), c.WriteMemory)
	case 0xEE:
		c.i_INC(c.getWithAbsoluteAddress(), c.WriteMemory)
	case 0xF6:
		c.i_INC(c.getWithZeroPageIndexedAddr(c.X), c.WriteMemory)
	case 0xFE:
		c.i_INC(c.getWithAbsoluteIndexedAddr(c.X), c.WriteMemory)

	case 0xC6:
		c.i_DEC(c.getWithZeroPageAddress(), c.WriteMemory)
	case 0xCE:
		c.i_DEC(c.getWithAbsoluteAddress(), c.WriteMemory)
	case 0xD6:
		c.i_DEC(c.getWithZeroPageIndexedAddr(c.X), c.WriteMemory)
	case 0xDE:
		c.i_DEC(c.getWithAbsoluteIndexedAddr(c.X), c.WriteMemory)

	case 0xE8:
		c.i_INX()
	case 0xC8:
		c.i_INY()
	case 0xCA:
		c.i_DEX()
	case 0x88:
		c.i_DEY()

	case 0x81:
		c.i_STA(c.getWithXIndexIndirectAddr())
	case 0x91:
		c.i_STA(c.getWithIndirectYIndexAddr())
	case 0x85:
		c.i_STA(c.getWithZeroPageAddress())
	case 0x95:
		c.i_STA(c.getWithZeroPageIndexedAddr(c.X))
	case 0x99:
		c.i_STA(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x8D:
		c.i_STA(c.getWithAbsoluteAddress())
	case 0x9D:
		c.i_STA(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x86:
		c.i_STX(c.getWithZeroPageAddress())
	case 0x96:
		c.i_STX(c.getWithZeroPageIndexedAddr(c.Y))
	case 0x8E:
		c.i_STX(c.getWithAbsoluteAddress())

	case 0x84:
		c.i_STY(c.getWithZeroPageAddress())
	case 0x94:
		c.i_STY(c.getWithZeroPageIndexedAddr(c.X))
	case 0x8C:
		c.i_STY(c.getWithAbsoluteAddress())

	case 0xA1:
		c.i_LDA(c.getWithXIndexIndirectAddr())
	case 0xB1:
		c.i_LDA(c.getWithIndirectYIndexAddr())
	case 0xA5:
		c.i_LDA(c.getWithZeroPageAddress())
	case 0xB5:
		c.i_LDA(c.getWithZeroPageIndexedAddr(c.X))
	case 0xA9:
		c.i_LDA(c.getWithImmediate())
	case 0xB9:
		c.i_LDA(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0xAD:
		c.i_LDA(c.getWithAbsoluteAddress())
	case 0xBD:
		c.i_LDA(c.getWithAbsoluteIndexedAddr(c.X))

	case 0xA2:
		c.i_LDX(c.getWithImmediate())
	case 0xA6:
		c.i_LDX(c.getWithZeroPageAddress())
	case 0xB6:
		c.i_LDX(c.getWithZeroPageIndexedAddr(c.Y))
	case 0xAE:
		c.i_LDX(c.getWithAbsoluteAddress())
	case 0xBE:
		c.i_LDX(c.getWithAbsoluteIndexedAddr(c.Y))

	case 0xA0:
		c.i_LDY(c.getWithImmediate())
	case 0xA4:
		c.i_LDX(c.getWithZeroPageAddress())
	case 0xB4:
		c.i_LDY(c.getWithZeroPageIndexedAddr(c.Y))
	case 0xAC:
		c.i_LDY(c.getWithAbsoluteAddress())
	case 0xBC:
		c.i_LDY(c.getWithAbsoluteIndexedAddr(c.Y))

	case 0xAA:
		c.i_TAX()
	case 0xA8:
		c.i_TAY()
	case 0x8A:
		c.i_TXA()
	case 0x98:
		c.i_TYA()
	case 0x9A:
		c.i_TXS()
	case 0xBA:
		c.i_TSX()

	case 0x06:
		c.i_ASL(c.getWithZeroPageAddress(), setResultFactory(c))
	case 0x16:
		c.i_ASL(c.getWithZeroPageIndexedAddr(c.X), setResultFactory(c))
	case 0x0A:
		c.i_ASL(c.A, func(v byte) { c.A = v })
	case 0x0E:
		c.i_ASL(c.getWithAbsoluteAddress(), setResultFactory(c))
	case 0x1E:
		c.i_ASL(c.getWithAbsoluteIndexedAddr(c.X), setResultFactory(c))

	case 0x46:
		c.i_LSR(c.getWithZeroPageAddress(), setResultFactory(c))
	case 0x56:
		c.i_LSR(c.getWithZeroPageIndexedAddr(c.X), setResultFactory(c))
	case 0x4A:
		c.i_LSR(c.A, func(v byte) { c.A = v })
	case 0x4E:
		c.i_LSR(c.getWithAbsoluteAddress(), setResultFactory(c))
	case 0x5E:
		c.i_LSR(c.getWithAbsoluteIndexedAddr(c.X), setResultFactory(c))

	case 0x26:
		c.i_ROL(c.getWithZeroPageAddress(), setResultFactory(c))
	case 0x36:
		c.i_ROL(c.getWithZeroPageIndexedAddr(c.X), setResultFactory(c))
	case 0x2A:
		c.i_ROL(c.A, func(v byte) { c.A = v })
	case 0x2E:
		c.i_ROL(c.getWithAbsoluteAddress(), setResultFactory(c))
	case 0x3E:
		c.i_ROL(c.getWithAbsoluteIndexedAddr(c.X), setResultFactory(c))

	case 0x66:
		c.i_ROR(c.getWithZeroPageAddress(), setResultFactory(c))
	case 0x76:
		c.i_ROR(c.getWithZeroPageIndexedAddr(c.X), setResultFactory(c))
	case 0x6A:
		c.i_ROR(c.A, func(v byte) { c.A = v })
	case 0x6E:
		c.i_ROR(c.getWithAbsoluteAddress(), setResultFactory(c))
	case 0x7E:
		c.i_ROR(c.getWithAbsoluteIndexedAddr(c.X), setResultFactory(c))

	case 0xC1:
		c.i_CMP(c.getWithXIndexIndirectAddr())
	case 0xD1:
		c.i_CMP(c.getWithIndirectYIndexAddr())
	case 0xC5:
		c.i_CMP(c.getWithZeroPageAddress())
	case 0xD5:
		c.i_CMP(c.getWithZeroPageIndexedAddr(c.X))
	case 0xC9:
		c.i_CMP(c.getWithImmediate())
	case 0xD9:
		c.i_CMP(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0xCD:
		c.i_CMP(c.getWithAbsoluteAddress())
	case 0xDD:
		c.i_CMP(c.getWithAbsoluteIndexedAddr(c.X))

	case 0xC0:
		c.i_CPY(c.getWithImmediate())
	case 0xC4:
		c.i_CPY(c.getWithZeroPageAddress())
	case 0xCC:
		c.i_CPY(c.getWithAbsoluteAddress())

	case 0xE0:
		c.i_CPX(c.getWithImmediate())
	case 0xE4:
		c.i_CPX(c.getWithZeroPageAddress())
	case 0xEC:
		c.i_CPX(c.getWithAbsoluteAddress())

	case 0x10:
		c.i_BPL(c.getWithRelative())
	case 0x30:
		c.i_BMI(c.getWithRelative())
	case 0x50:
		c.i_BVC(c.getWithRelative())
	case 0x70:
		c.i_BVS(c.getWithRelative())
	case 0x90:
		c.i_BCC(c.getWithRelative())
	case 0xB0:
		c.i_BCS(c.getWithRelative())
	case 0xD0:
		c.i_BNE(c.getWithRelative())
	case 0xF0:
		c.i_BEQ(c.getWithRelative())

	case 0x4C:
		c.i_JMP(c.getAbsoluteAddr())
	case 0x6C:
		c.i_JMP(c.getIndirectAddr())

	case 0x20:
		c.i_JSR(c.getAbsoluteAddrNoPC())
	case 0x60:
		c.i_RTS()

	case 0x48:
		c.i_PHA()
	case 0x08:
		c.i_PHP()
	case 0x68:
		c.i_PLA()
	case 0x28:
		c.i_PLP()

	case 0x18:
		c.i_CLC()
	case 0x58:
		c.i_CLI()
	case 0xB8:
		c.i_CLV()
	case 0xD8:
		c.i_CLD()
	case 0x38:
		c.i_SEC()
	case 0x78:
		c.i_SEI()
	case 0xF8:
		c.i_SED()

	case 0xEA:
		c.i_NOP()

	default:
		c.i_JAM()
	}
}

func (c *CPU) ReadMemory(addr Word) byte {
	return c.Memory[addr]
}

func (c *CPU) WriteMemory(addr Word, value byte) {
	c.Memory[addr] = value
}

func (c *CPU) pushIntoStack(value byte) {
	effectiveStackAddr := 0x0100 | Word(c.SP)
	c.WriteMemory(effectiveStackAddr, value)
	c.SP--
}

func (c *CPU) pullFromStack() byte {
	c.SP++
	effectiveStackAddr := 0x0100 | Word(c.SP)
	value := c.ReadMemory(effectiveStackAddr)
	return value
}

func joinBytesToWord(lo, hi byte) Word {
	const ADDR_SHIFT = 8
	return Word(lo) | Word(hi)<<ADDR_SHIFT
}

func splitWordToByte(w Word) (byte, byte) {
	const (
		LO_ADDR_MASK = 0x00FF
		HI_ADDR_MASK = 0xFF00
		ADDR_SHIFT   = 8
	)

	return byte(w & LO_ADDR_MASK), byte((w & HI_ADDR_MASK) >> ADDR_SHIFT)
}
