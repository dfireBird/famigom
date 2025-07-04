package cpu

import (
	. "github.com/dfirebird/famigom/types"
)

const (
	ZERO_PAGE_MAX = 256

	HI_BYTE_MASK = 0xFF00
)

func (c *CPU) getWithAbsoluteAddress() byte {
	addr := c.getAbsoluteAddr()
	c.currentGetAddr = addr
	return c.ReadMemory(addr)
}

func (c *CPU) getWithZeroPageAddress() byte {
	addr := c.getZeroPageAddr()
	c.currentGetAddr = addr
	return c.ReadMemory(addr)
}

func (c *CPU) getWithAbsoluteIndexedAddr(idx byte, isAlwaysFixHi bool) byte {
	addr := c.getAbsoluteIndexedAddr(idx, isAlwaysFixHi)
	c.currentGetAddr = addr
	return c.ReadMemory(addr)
}

func (c *CPU) getWithZeroPageIndexedAddr(idx byte) byte {
	addr := c.getZeroPageIndexedAddr(idx)
	c.currentGetAddr = addr
	return c.ReadMemory(addr)
}

func (c *CPU) getWithXIndexIndirectAddr() byte {
	effectiveAddr := c.getXIndexIndirectAddr()
	c.currentGetAddr = effectiveAddr
	return c.ReadMemory(effectiveAddr)
}

func (c *CPU) getWithIndirectYIndexAddr(isAlwaysFixHi bool) byte {
	effectiveAddr := c.getIndirectYIndexAddr(isAlwaysFixHi)
	c.currentGetAddr = effectiveAddr
	return c.ReadMemory(effectiveAddr)
}

func (c *CPU) getWithImmediate() byte {
	v := c.ReadMemory(c.PC)
	c.PC++
	return v
}

func (c *CPU) getWithImplied(val byte) byte {
	c.ReadMemory(c.PC) // dummy read
	return val
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
func (c *CPU) getIndirectAddr() Word {
	lo := c.ReadMemory(c.PC)
	c.PC++
	hi := c.ReadMemory(c.PC)
	c.PC++

	loAddr, hiAddr := joinBytesToWord(lo, hi), joinBytesToWord((lo+1), hi)
	effectiveLo := c.ReadMemory(loAddr)
	effectiveHi := c.ReadMemory(hiAddr)
	return joinBytesToWord(effectiveLo, effectiveHi)
}

func (c *CPU) getAbsoluteAddr() Word {
	lo := c.ReadMemory(c.PC)
	c.PC++
	hi := c.ReadMemory(c.PC)
	c.PC++

	return joinBytesToWord(lo, hi)
}

func (c *CPU) getAbsoluteAddrNoPC() Word {
	lo := c.ReadMemory(c.PC)
	hi := c.ReadMemory(c.PC + 1)
	return joinBytesToWord(lo, hi)
}

func (c *CPU) getZeroPageAddr() Word {
	lo := c.ReadMemory(c.PC)
	c.PC++

	return joinBytesToWord(lo, 0x00)
}

func (c *CPU) getAbsoluteIndexedAddr(idx byte, isAlwaysFixHi bool) Word {
	baseAddr := c.getAbsoluteAddr()
	effectiveAddr := baseAddr + Word(idx)

	if isAlwaysFixHi || baseAddr&HI_BYTE_MASK != effectiveAddr&HI_BYTE_MASK {
		c.ReadMemory(effectiveAddr)
	}
	return effectiveAddr
}

func (c *CPU) getZeroPageIndexedAddr(idx byte) Word {
	baseAddr := c.getZeroPageAddr()
	addr := (baseAddr + Word(idx)) % ZERO_PAGE_MAX

	c.ReadMemory(baseAddr)

	return addr
}

func (c *CPU) getXIndexIndirectAddr() Word {
	// The mode uses a zero page address
	addr := c.getZeroPageAddr()

	c.ReadMemory(addr) // dummy read
	lo := c.ReadMemory((addr + Word(c.X)) % ZERO_PAGE_MAX)
	hi := c.ReadMemory((addr + Word(c.X) + 1) % ZERO_PAGE_MAX)

	return joinBytesToWord(lo, hi)
}

func (c *CPU) getIndirectYIndexAddr(isAlwaysFixHi bool) Word {
	// The mode uses a zero page address
	addr := c.getZeroPageAddr()

	lo := c.ReadMemory(addr)
	hi := c.ReadMemory((addr + 1) % ZERO_PAGE_MAX)

	baseAddr := joinBytesToWord(lo, hi)
	effectiveAddr := baseAddr + Word(c.Y)

	if isAlwaysFixHi || baseAddr&HI_BYTE_MASK != effectiveAddr&HI_BYTE_MASK {
		c.ReadMemory(effectiveAddr)
	}

	return effectiveAddr
}
