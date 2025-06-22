package cpu

type status byte

const (
	NEGATIVE_SHIFT     = 7
	NEGATIVE_BIT_MASK  = 0b1000_0000

	OVERFLOW_SHIFT     = 6
	OVERFLOW_BIT_MASK  = 0b0100_0000

	BREAK_BIT_MASK     = 0b0001_0000

	DECIMAL_SHIFT      = 3
	DECIMAL_BIT_MASK   = 0b0000_1000

	INTERRUPT_SHIFT    = 2
	INTERRUPT_BIT_MASK = 0b0000_0100

	ZERO_SHIFT         = 1
	ZERO_BIT_MASK      = 0b0000_0010

	CARRY_SHIFT        = 0
	CARRY_BIT_MASK     = 0b0000_0001

	INITIAL_STATUS = 0b0010_0000
)

func not(value byte) byte {
	return 0b1111_1111 ^ value
}

func (s *status) SetValueFromStack(value byte) {
	newStatus := status(value)
	s.SetCarry(newStatus.GetCarry())
	s.SetZero(newStatus.GetZero())
	s.SetInterruptDisable(newStatus.GetInterruptDisable())
	s.SetDecimal(newStatus.GetDecimal())
	s.SetOverflow(newStatus.GetOverflow())
	s.SetNegative(newStatus.GetNegative())
}

func (s *status) SetNegative(value bool) {
	s.setXBit(value, NEGATIVE_BIT_MASK)
}

func (s *status) GetNegative() bool {
	return s.getXBit(NEGATIVE_BIT_MASK, NEGATIVE_SHIFT)
}

func (s *status) SetOverflow(value bool) {
	s.setXBit(value, OVERFLOW_BIT_MASK)
}

func (s *status) GetOverflow() bool {
	return s.getXBit(OVERFLOW_BIT_MASK, OVERFLOW_SHIFT)
}

func (s *status) SetDecimal(value bool) {
	s.setXBit(value, DECIMAL_BIT_MASK)
}

func (s *status) GetDecimal() bool {
	return s.getXBit(DECIMAL_BIT_MASK, DECIMAL_SHIFT)
}

func (s *status) SetInterruptDisable(value bool) {
	s.setXBit(value, INTERRUPT_BIT_MASK)
}

func (s *status) GetInterruptDisable() bool {
	return s.getXBit(INTERRUPT_BIT_MASK, INTERRUPT_SHIFT)
}

func (s *status) SetZero(value bool) {
	s.setXBit(value, ZERO_BIT_MASK)
}

func (s *status) GetZero() bool {
	return s.getXBit(ZERO_BIT_MASK, ZERO_SHIFT)
}

func (s *status) SetCarry(value bool) {
	s.setXBit(value, CARRY_BIT_MASK)
}

func (s *status) GetCarry() bool {
	return s.getXBit(CARRY_BIT_MASK, CARRY_SHIFT)
}

func (s *status) GetCarryNum() byte {
	if s.GetCarry() {
		return 1
	} else {
		return 0
	}
}

func (s *status) setXBit(value bool, mask byte) {
	if value {
		*s |= status (mask)
	} else {
		*s &= status (not(mask))
	}
}

func (s *status) getXBit(mask, shift byte) bool {
	return (*s & status(mask)) >> shift == 1
}
