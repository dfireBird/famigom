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

func (c *CPU) Step() {
	opcode := c.readMemory(c.PC)
	c.PC++
	switch opcode {
	case 0x00:
		c.brk()

	case 0x01:
		c.or(c.getWithXIndexIndirectAddr())
	case 0x11:
		c.or(c.getWithIndirectYIndexAddr())
	case 0x05:
		c.or(c.getWithZeroPageAddress())
	case 0x15:
		c.or(c.getWithZeroPageIndexedAddr(c.X))
	case 0x09:
		c.or(c.getWithImmediate())
	case 0x19:
		c.or(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x0D:
		c.or(c.getWithAbsoluteAddress())
	case 0x1D:
		c.or(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x21:
		c.and(c.getWithXIndexIndirectAddr())
	case 0x31:
		c.and(c.getWithIndirectYIndexAddr())
	case 0x25:
		c.and(c.getWithZeroPageAddress())
	case 0x35:
		c.and(c.getWithZeroPageIndexedAddr(c.X))
	case 0x29:
		c.and(c.getWithImmediate())
	case 0x39:
		c.and(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x2D:
		c.and(c.getWithAbsoluteAddress())
	case 0x3D:
		c.and(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x41:
		c.xor(c.getWithXIndexIndirectAddr())
	case 0x51:
		c.xor(c.getWithIndirectYIndexAddr())
	case 0x45:
		c.xor(c.getWithZeroPageAddress())
	case 0x55:
		c.xor(c.getWithZeroPageIndexedAddr(c.X))
	case 0x49:
		c.xor(c.getWithImmediate())
	case 0x59:
		c.xor(c.getWithAbsoluteIndexedAddr(c.Y))
	case 0x4D:
		c.xor(c.getWithAbsoluteAddress())
	case 0x5D:
		c.xor(c.getWithAbsoluteIndexedAddr(c.X))

	case 0x24:
		c.bit(c.getWithZeroPageAddress())
	case 0x2C:
		c.bit(c.getWithAbsoluteAddress())
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
