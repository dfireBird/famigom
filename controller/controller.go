package controller

import (
	"github.com/dfirebird/famigom/types"
)

type Controllers struct {
	strobe bool

	port1 byte
	port2 byte
}

func CreateControllers() *Controllers {
	return &Controllers{
		strobe: false,
	}
}

func (c *Controllers) LoadButtonData(port1, port2 byte) {
	if c.strobe {
		// log.GetLoggerWithSpan("controller").Debugf("Load Controller Data p1: 0b%08b p2: 0b%08b", port1, port2)
		c.port1 = port1
		c.port2 = port2
	}
}

func (c *Controllers) ReadMemory(addr types.Word) (bool, byte) {
	if addr == 0x4016 {
		// log.GetLoggerWithSpan("controller").Debugf("Read Controller port: %d strobe: %t value: 0b%08b", 1, c.strobe, c.port1)
		return true, handleControllerRead(&c.port1, c.strobe)
	}
	if addr == 0x4017 {
		return true, handleControllerRead(&c.port2, c.strobe)
	}
	return false, 0
}

func (c *Controllers) WriteMemory(addr types.Word, value byte) {
	if addr == 0x4016 {
		c.strobe = (value & 1) == 1
	}
}

func handleControllerRead(port *byte, strobe bool) byte {
	value := readController(*port)

	if !strobe {
		*port = shiftAndSet(*port)
	}
	return value
}

func readController(portValue byte) byte {
	return portValue & 1
}

func shiftAndSet(portValue byte) byte {
	return (portValue >> 1) | 0b1000_0000
}
