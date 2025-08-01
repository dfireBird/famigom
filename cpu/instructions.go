package cpu

import (
	. "github.com/dfirebird/famigom/types"
)

const (
	BIT_7_NEG_SET = 0b1000_0000
	BIT_8_OVL_SET = 0x100

	IRQ_BRK_HANDLER_ADDR = 0xFFFE
)

func (c *CPU) i_BRK() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.PC++
	c.pushAddrIntoStack(c.PC)
	c.pushIntoStack(byte(c.Flags) | BREAK_BIT_MASK)
	c.Flags.SetInterruptDisable(true)

	handler_routine_addr := joinBytesToWord(c.ReadMemory(IRQ_BRK_HANDLER_ADDR), c.ReadMemory(IRQ_BRK_HANDLER_ADDR+1))
	c.PC = handler_routine_addr
}

func (c *CPU) i_ORA(op byte) {
	c.A |= op

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) i_AND(op byte) {
	c.A &= op

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) i_EOR(op byte) {
	c.A ^= op

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) i_BIT(op byte) {
	res := c.A & op
	bit6 := (op & (1 << 6)) >> 6
	bit7 := (op & (1 << 7)) >> 7

	c.Flags.SetZero(res == 0)
	c.Flags.SetOverflow(bit6 == 1)
	c.Flags.SetNegative(bit7 == 1)
}

func (c *CPU) i_ADC(op byte) {
	c.add(op)
}

func (c *CPU) i_SBC(op byte) {
	c.add(^op)
}

func (c *CPU) i_INC(op byte, memSet func(Word, byte)) {
	res := c.increment(op)
	memSet(c.currentGetAddr, op)
	memSet(c.currentGetAddr, res)
}

func (c *CPU) i_INX() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.X = c.increment(c.X)
}

func (c *CPU) i_INY() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Y = c.increment(c.Y)
}

func (c *CPU) i_DEC(op byte, memSet func(Word, byte)) {
	res := c.decrement(op)
	memSet(c.currentGetAddr, op)
	memSet(c.currentGetAddr, res)
}

func (c *CPU) i_DEX() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.X = c.decrement(c.X)
}

func (c *CPU) i_DEY() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Y = c.decrement(c.Y)
}

func (c *CPU) i_LDA(value byte) {
	c.A = c.load(value)
}

func (c *CPU) i_LDX(value byte) {
	c.X = c.load(value)
}

func (c *CPU) i_LDY(value byte) {
	c.Y = c.load(value)
}

func (c *CPU) i_STA(addr Word) {
	c.WriteMemory(addr, c.A)
}

func (c *CPU) i_STX(addr Word) {
	c.WriteMemory(addr, c.X)
}

func (c *CPU) i_STY(addr Word) {
	c.WriteMemory(addr, c.Y)
}

func (c *CPU) i_TAX() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.X = c.A

	c.Flags.SetZero(c.X == 0)
	c.Flags.SetNegative(isNegative(c.X))
}

func (c *CPU) i_TAY() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Y = c.A

	c.Flags.SetZero(c.Y == 0)
	c.Flags.SetNegative(isNegative(c.Y))
}

func (c *CPU) i_TXA() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.A = c.X

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) i_TYA() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.A = c.Y

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) i_ASL(value byte, resSet func(byte)) {
	c.shift(value, resSet, true)
}

func (c *CPU) i_LSR(value byte, resSet func(byte)) {
	c.shift(value, resSet, false)
}

func (c *CPU) i_ROL(value byte, resSet func(byte)) {
	c.rotate(value, resSet, true)
}

func (c *CPU) i_ROR(value byte, resSet func(byte)) {
	c.rotate(value, resSet, false)
}

func (c *CPU) i_CMP(value byte) {
	c.compare(c.A, value)
}

func (c *CPU) i_CPX(value byte) {
	c.compare(c.X, value)
}

func (c *CPU) i_CPY(value byte) {
	c.compare(c.Y, value)
}

func (c *CPU) i_JMP(addr Word) {
	c.PC = addr
}

func (c *CPU) i_JSR(addr Word) {
	c.dummyOperationStack()
	c.pushAddrIntoStack(c.PC + 1)
	c.PC = addr
}

func (c *CPU) i_RTS() {
	c.ReadMemory(c.PC)      // dummy read for implied addressing
	c.dummyOperationStack() // dummy pull to increment SP
	pc := c.pullAddrFromStack()

	c.ReadMemory(pc) // dummy read for incrementing PC
	c.PC = pc + 1
}

func (c *CPU) i_RTI() {
	c.ReadMemory(c.PC)      // dummy read for implied addressing
	c.dummyOperationStack() // dummy pull to increment SP
	flags := c.pullFromStack()
	pc := c.pullAddrFromStack()

	c.Flags.SetValueFromStack(flags)
	c.PC = pc
}

func (c *CPU) i_BCS(offset byte) {
	if c.Flags.GetCarry() {
		c.branch(offset)
	}
}

func (c *CPU) i_BCC(offset byte) {
	if !c.Flags.GetCarry() {
		c.branch(offset)
	}
}

func (c *CPU) i_BEQ(offset byte) {
	if c.Flags.GetZero() {
		c.branch(offset)
	}
}

func (c *CPU) i_BNE(offset byte) {
	if !c.Flags.GetZero() {
		c.branch(offset)
	}
}

func (c *CPU) i_BMI(offset byte) {
	if c.Flags.GetNegative() {
		c.branch(offset)
	}
}

func (c *CPU) i_BPL(offset byte) {
	if !c.Flags.GetNegative() {
		c.branch(offset)
	}
}

func (c *CPU) i_BVS(offset byte) {
	if c.Flags.GetOverflow() {
		c.branch(offset)
	}
}

func (c *CPU) i_BVC(offset byte) {
	if !c.Flags.GetOverflow() {
		c.branch(offset)
	}
}

func (c *CPU) i_PHA() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.pushIntoStack(c.A)
}

func (c *CPU) i_PLA() {
	c.ReadMemory(c.PC)      // dummy read for implied addressing
	c.dummyOperationStack() // dummy pull to increment SP
	c.A = c.pullFromStack()

	c.Flags.SetZero(c.A == 0)
	c.Flags.SetNegative(isNegative(c.A))
}

func (c *CPU) i_PHP() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.pushIntoStack(byte(c.Flags) | BREAK_BIT_MASK)
}

func (c *CPU) i_PLP() {
	c.ReadMemory(c.PC)      // dummy read for implied addressing
	c.dummyOperationStack() // dummy pull to increment SP
	c.Flags.SetValueFromStack(c.pullFromStack())
}

func (c *CPU) i_TXS() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.SP = c.X
}

func (c *CPU) i_TSX() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.X = c.SP

	c.Flags.SetZero(c.X == 0)
	c.Flags.SetNegative(isNegative(c.X))
}

func (c *CPU) i_CLC() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Flags.SetCarry(false)
}

func (c *CPU) i_CLD() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Flags.SetDecimal(false)
}

func (c *CPU) i_CLI() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Flags.SetInterruptDisable(false)
}

func (c *CPU) i_CLV() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Flags.SetOverflow(false)
}

func (c *CPU) i_SEC() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Flags.SetCarry(true)
}

func (c *CPU) i_SED() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Flags.SetDecimal(true)
}

func (c *CPU) i_SEI() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
	c.Flags.SetInterruptDisable(true)
}

func (c *CPU) i_NOP() {
	c.ReadMemory(c.PC) // dummy read for implied addressing
}

func (c *CPU) i_JAM() {
	// FIXME: It should set a variable in CPU so that CPU runs but doesn't
	// execute instructions only poll for hardware interupts
	c.isJammed = true
	panic("JAM instruction")
}

func (c *CPU) branch(offset byte) {
	c.ReadMemory(c.PC)
	signedOffset := int32(int8(offset))

	oldPC := c.PC
	c.PC = uint16(int32(c.PC) + signedOffset)

	if c.PC&HI_BYTE_MASK != oldPC&HI_BYTE_MASK {
		c.ReadMemory(c.PC)
	}
}

func (c *CPU) compare(register, value byte) {
	c.Flags.SetCarry(register >= value)
	c.Flags.SetZero(register == value)
	c.Flags.SetNegative(isNegative(register - value))
}

func (c *CPU) shift(value byte, resSet func(byte), isLeft bool) {
	var result byte
	var carry bool
	if isLeft {
		result = value << 1
		carry = isNegative(value) // using isNegative since they both use bit 7
	} else {
		result = value >> 1
		carry = (value & 1) == 1
	}

	c.Flags.SetCarry(carry)
	c.Flags.SetNegative(isNegative(result))
	c.Flags.SetZero(result == 0)

	resSet(value) // dummy write only emulate cycle, also writes to accumulator twice not sure of implications
	resSet(result)
}

func (c *CPU) rotate(value byte, resSet func(byte), isLeft bool) {
	var result byte
	var carry bool
	if carryValue := c.Flags.GetCarryNum(); isLeft {
		result = (value << 1) | carryValue
		carry = isNegative(value) // using isNegative since they both use bit 7
	} else {
		result = (value >> 1) | (carryValue << 7)
		carry = (value & 1) == 1
	}

	c.Flags.SetCarry(carry)
	c.Flags.SetNegative(isNegative(result))
	c.Flags.SetZero(result == 0)

	resSet(value) // dummy write only emulate cycle, also writes to accumulator twice not sure of implications
	resSet(result)
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
	resultAs16Bit := Word(c.A) + Word(op) + Word(carry)
	c.Flags.SetCarry((resultAs16Bit & BIT_8_OVL_SET) == BIT_8_OVL_SET)

	c.Flags.SetZero(result == 0)
	c.Flags.SetNegative(isNegative(result))
	c.Flags.SetOverflow(((result ^ c.A) & (result ^ op) & BIT_7_NEG_SET) == BIT_7_NEG_SET)

	c.A = result
}

func (c *CPU) pushAddrIntoStack(value Word) {
	lo, hi := splitWordToByte(value)
	c.pushIntoStack(hi)
	c.pushIntoStack(lo)
}

func (c *CPU) pullAddrFromStack() Word {
	lo := c.pullFromStack()
	hi := c.pullFromStack()

	return joinBytesToWord(lo, hi)
}

func isNegative(val byte) bool {
	return val&BIT_7_NEG_SET == BIT_7_NEG_SET
}
