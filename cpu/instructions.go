package cpu

const (
	BIT_7_NEG_SET = 0b1000_0000
	BIT_8_OVL_SET = 0x100

	IRQ_BRK_HANDLER_ADDR = 0xFFFE
)

func (c *CPU) BRK() {
	c.pushAddrIntoStack(c.PC + 1)
	c.pushIntoStack(c.SP | BREAK_BIT_MASK)
	c.Flags.SetInterruptDisable(true)

	c.PC = IRQ_BRK_HANDLER_ADDR
}

func (c *CPU) ORA(op byte) {
	c.A |= op

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) AND(op byte) {
	c.A &= op

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) EOR(op byte) {
	c.A ^= op

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) BIT(op byte) {
	res := c.A & op
	bit6 := (op & (1 << 6)) >> 6
	bit7 := (op & (1 << 7)) >> 7

	c.Flags.SetZero(res == 0)
	c.Flags.SetOverflow(bit6 == 1)
	c.Flags.SetNegative(bit7 == 1)
}

func (c *CPU) ADC(op byte) {
	c.add(op)
}

func (c *CPU) SBC(op byte) {
	c.add(^op)
}

func (c *CPU) INC(op byte, memSet func(word, byte)) {
	res := c.increment(op)
	memSet(c.currentGetAddr, res)
}

func (c *CPU) INX() {
	c.X = c.increment(c.X)
}

func (c *CPU) INY() {
	c.Y = c.increment(c.Y)
}

func (c *CPU) DEC(op byte, memSet func(word, byte)) {
    res := c.decrement(op)
	memSet(c.currentGetAddr, res)
}

func (c *CPU) DEX() {
	c.X = c.decrement(c.X)
}

func (c *CPU) DEY() {
	c.Y = c.decrement(c.Y)
}

func (c *CPU) LDA(value byte) {
	c.A = c.load(value)
}

func (c *CPU) LDX(value byte) {
	c.X = c.load(value)
}

func (c *CPU) LDY(value byte) {
	c.Y = c.load(value)
}

// Stores use getXXXAddr functions result as parameter but doesn't use it
// But uses the currentGetAddr as address value to write to
func (c *CPU) STA(_ byte) {
	c.writeMemory(c.currentGetAddr, c.A)
}

func (c *CPU) STX(_ byte) {
	c.writeMemory(c.currentGetAddr, c.X)
}

func (c *CPU) STY(_ byte) {
	c.writeMemory(c.currentGetAddr, c.Y)
}

func (c *CPU) TAX() {
	c.X = c.A

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetZero(isNegative(c.A))
}

func (c *CPU) TAY() {
	c.Y = c.A

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetZero(isNegative(c.A))
}

func (c *CPU) TXA() {
	c.A = c.X

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetZero(isNegative(c.A))
}

func (c *CPU) TYA() {
	c.A = c.Y

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetZero(isNegative(c.A))
}

func (c *CPU) load(value byte) byte {
	c.Flags.SetZero(value == 0)
	c.Flags.SetNegative(isNegative(value))

	return value
}

func (c *CPU) increment(op byte) byte {
	res := op + 1

	c.Flags.SetZero(res == 0)
	c.Flags.SetNegative(isNegative(res))

	return res
}

func (c *CPU) decrement(op byte) byte {
	res := op - 1

	c.Flags.SetZero(res == 0)
	c.Flags.SetNegative(isNegative(res))

	return res
}

func (c *CPU) add(op byte) {
	carry := c.Flags.GetCarryNum()
	result := c.A + op + carry

	// Masking the 8th bit in 16bit value for carry
	resultAs16Bit := word(c.A) + word(op) + word(carry)
	c.Flags.SetCarry((resultAs16Bit & BIT_8_OVL_SET) == BIT_8_OVL_SET)

	c.Flags.SetZero(result == 0)
	c.Flags.SetNegative(isNegative(result))
	c.Flags.SetOverflow(((result ^ c.A) & (result ^ op) & BIT_7_NEG_SET) == BIT_7_NEG_SET)

	c.A = result
}

func (c *CPU) pushAddrIntoStack(value word) {
	lo, hi := splitWordToByte(value)
	c.pushIntoStack(lo)
	c.pushIntoStack(hi)
}

func isNegative(val byte) bool {
	return val&BIT_7_NEG_SET == BIT_7_NEG_SET
}
