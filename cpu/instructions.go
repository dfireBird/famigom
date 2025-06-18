package cpu

const (
	HI_BIT_SET  = 0b1000_0000
	HI_BIT_MASK = 1 << 7

	IRQ_BRK_HANDLER_ADDR = 0xFFFE
)

func (c *CPU) brk() {
	c.pushAddrIntoStack(c.PC + 1)
	c.pushIntoStack(c.SP | BREAK_BIT_MASK)
	c.SR.SetInterruptDisable(true)

	c.PC = IRQ_BRK_HANDLER_ADDR
}

func (c *CPU) or(op byte) {
	c.A |= op

	c.SR.SetZero(c.A == 0)
	c.SR.SetNegative(c.A&HI_BIT_MASK == HI_BIT_SET)
}

func (c *CPU) and(op byte) {
	c.A &= op

	c.SR.SetZero(c.A == 0)
	c.SR.SetNegative(c.A&HI_BIT_MASK == HI_BIT_SET)
}

func (c *CPU) xor(op byte) {
	c.A ^= op

	c.SR.SetZero(c.A == 0)
	c.SR.SetNegative(c.A&HI_BIT_MASK == HI_BIT_SET)
}

func (c *CPU) bit(op byte) {
    res := c.A & op
	bit6 := (op & (1<<6)) >> 6
	bit7 := (op & (1<<7)) >> 7

	c.SR.SetZero(res == 0)
	c.SR.SetOverflow(bit6 == 1)
	c.SR.SetNegative(bit7 == 1)
}

func (c *CPU) pushAddrIntoStack(value word) {
	lo, hi := splitWordToByte(value)
	c.pushIntoStack(lo)
	c.pushIntoStack(hi)
}
