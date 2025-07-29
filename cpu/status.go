package cpu

type Status byte

const (
	NEGATIVE_SHIFT    = 7
	NEGATIVE_BIT_MASK = 0b1000_0000

	OVERFLOW_SHIFT    = 6
	OVERFLOW_BIT_MASK = 0b0100_0000

	BREAK_BIT_MASK = 0b0001_0000

	DECIMAL_SHIFT    = 3
	DECIMAL_BIT_MASK = 0b0000_1000

	INTERRUPT_SHIFT    = 2
	INTERRUPT_BIT_MASK = 0b0000_0100

	ZERO_SHIFT    = 1
	ZERO_BIT_MASK = 0b0000_0010

	CARRY_SHIFT    = 0
	CARRY_BIT_MASK = 0b0000_0001

	INITIAL_STATUS = 0b0010_0000
)

func (s *Status) SetValueFromStack(value byte) {
	newStatus := Status(value)
	s.SetCarry(newStatus.GetCarry())
	s.SetZero(newStatus.GetZero())
	s.SetInterruptDisable(newStatus.GetInterruptDisable())
	s.SetDecimal(newStatus.GetDecimal())
	s.SetOverflow(newStatus.GetOverflow())
	s.SetNegative(newStatus.GetNegative())
}

func (s *Status) SetNegative(value bool) {
	s.setXBit(value, NEGATIVE_BIT_MASK)
}

func (s *Status) GetNegative() bool {
	return s.getXBit(NEGATIVE_BIT_MASK, NEGATIVE_SHIFT)
}

func (s *Status) SetOverflow(value bool) {
	s.setXBit(value, OVERFLOW_BIT_MASK)
}

func (s *Status) GetOverflow() bool {
	return s.getXBit(OVERFLOW_BIT_MASK, OVERFLOW_SHIFT)
}

func (s *Status) SetDecimal(value bool) {
	s.setXBit(value, DECIMAL_BIT_MASK)
}

func (s *Status) GetDecimal() bool {
	return s.getXBit(DECIMAL_BIT_MASK, DECIMAL_SHIFT)
}

func (s *Status) SetInterruptDisable(value bool) {
	s.setXBit(value, INTERRUPT_BIT_MASK)
}

func (s *Status) GetInterruptDisable() bool {
	return s.getXBit(INTERRUPT_BIT_MASK, INTERRUPT_SHIFT)
}

func (s *Status) SetZero(value bool) {
	s.setXBit(value, ZERO_BIT_MASK)
}

func (s *Status) GetZero() bool {
	return s.getXBit(ZERO_BIT_MASK, ZERO_SHIFT)
}

func (s *Status) SetCarry(value bool) {
	s.setXBit(value, CARRY_BIT_MASK)
}

func (s *Status) GetCarry() bool {
	return s.getXBit(CARRY_BIT_MASK, CARRY_SHIFT)
}

func (s *Status) GetCarryNum() byte {
	if s.GetCarry() {
		return 1
	} else {
		return 0
	}
}

func (s *Status) setXBit(value bool, mask byte) {
	if value {
		*s |= Status(mask)
	} else {
		*s &^= Status(mask)
	}
}

func (s *Status) getXBit(mask, shift byte) bool {
	return (*s&Status(mask))>>shift == 1
}
