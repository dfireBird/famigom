package cpu

const (
	ZERO_PAGE_MAX = 256
)

func (c *CPU) getWithAbsoluteAddress() byte {
	addr := c.getAbsoluteAddr()
	return c.readMemory(addr)
}

func (c *CPU) getWithZeroPageAddress() byte {
	addr := c.getZeroPageAddr()
	return c.readMemory(addr)
}

func (c *CPU) getWithAbsoluteIndexedAddr(idx byte) byte {
	addr := c.getAbsoluteAddr() + word(idx)
	return c.readMemory(addr)
}

func (c *CPU) getWithZeroPageIndexedAddr(idx byte) byte {
	addr := (c.getZeroPageAddr() + word(idx)) % ZERO_PAGE_MAX
	return c.readMemory(addr)
}

func (c *CPU) getWithXIndexIndirectAddr() byte {
	// The mode uses a zero page address
	addr := c.getZeroPageAddr()

	lo := c.readMemory((addr + word(c.X)) % ZERO_PAGE_MAX)
	hi := c.readMemory((addr + word(c.X) + 1) % ZERO_PAGE_MAX)

	effectiveAddr := joinBytesToWord(lo, hi)
	return c.readMemory(effectiveAddr)
}

func (c *CPU) getWithIndirectYIndexAddr() byte {
	// The mode uses a zero page address
	addr := c.getZeroPageAddr()

	lo := c.readMemory(addr)
	hi := c.readMemory((addr + 1) % ZERO_PAGE_MAX)

	effectiveAddr := joinBytesToWord(lo, hi) + word(c.Y)
	return c.readMemory(effectiveAddr)
}

func (c *CPU) getWithImmediate() byte {
	v := c.readMemory(c.PC)
	c.PC++
	return v
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
