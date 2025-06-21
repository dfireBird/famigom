package cpu

const (
	ZERO_PAGE_MAX = 256
)

func (c *CPU) getWithAbsoluteAddress() byte {
	addr := c.getAbsoluteAddr()
	c.currentGetAddr = addr
	return c.readMemory(addr)
}

func (c *CPU) getWithZeroPageAddress() byte {
	addr := c.getZeroPageAddr()
	c.currentGetAddr = addr
	return c.readMemory(addr)
}

func (c *CPU) getWithAbsoluteIndexedAddr(idx byte) byte {
	addr := c.getAbsoluteAddr() + word(idx)
	c.currentGetAddr = addr
	return c.readMemory(addr)
}

func (c *CPU) getWithZeroPageIndexedAddr(idx byte) byte {
	addr := (c.getZeroPageAddr() + word(idx)) % ZERO_PAGE_MAX
	c.currentGetAddr = addr
	return c.readMemory(addr)
}

func (c *CPU) getWithXIndexIndirectAddr() byte {
	// The mode uses a zero page address
	addr := c.getZeroPageAddr()

	lo := c.readMemory((addr + word(c.X)) % ZERO_PAGE_MAX)
	hi := c.readMemory((addr + word(c.X) + 1) % ZERO_PAGE_MAX)

	effectiveAddr := joinBytesToWord(lo, hi)
	c.currentGetAddr = effectiveAddr
	return c.readMemory(effectiveAddr)
}

func (c *CPU) getWithIndirectYIndexAddr() byte {
	// The mode uses a zero page address
	addr := c.getZeroPageAddr()

	lo := c.readMemory(addr)
	hi := c.readMemory((addr + 1) % ZERO_PAGE_MAX)

	effectiveAddr := joinBytesToWord(lo, hi) + word(c.Y)
	c.currentGetAddr = effectiveAddr
	return c.readMemory(effectiveAddr)
}

func (c *CPU) getWithImmediate() byte {
	v := c.readMemory(c.PC)
	c.PC++
	return v
}

// Only used with branching instructions
// and only understood as relative by branching instructions
// but similar to getWithImmediate
func (c *CPU) getWithRelative() byte {
	return c.getWithImmediate()
}

// Indirect addressing mode is basically retrieving
// absolute address from the ROM and using that address
// to retrieve from memory the effective address
// Only used in JMP
func (c *CPU) getIndirectAddr() word {
	lo := c.readMemory(c.PC)
	c.PC++
	hi := c.readMemory(c.PC)
	c.PC++


	loAddr, hiAddr := joinBytesToWord(lo, hi), joinBytesToWord((lo + 1), hi)
	effectiveLo := c.readMemory(loAddr)
	effectiveHi := c.readMemory(hiAddr)
	return joinBytesToWord(effectiveLo, effectiveHi)
}

func (c *CPU) getAbsoluteAddr() word {
	lo := c.readMemory(c.PC)
	c.PC++
	hi := c.readMemory(c.PC)
	c.PC++

	return joinBytesToWord(lo, hi)
}

func (c *CPU) getZeroPageAddr() word {
	lo := c.readMemory(c.PC)
	c.PC++

	return joinBytesToWord(lo, 0x00)
}
