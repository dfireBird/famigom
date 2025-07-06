package bus

import (
	. "github.com/dfirebird/famigom/types"
)

type MemoryBus interface {
	ReadMemory(addr Word) byte
	WriteMemory(addr Word, value byte)
}

type MainBus struct {
	devicesMap map[AddrRange](*MemoryBus)
}

func (b *MainBus) ReadMemory(addr Word) byte {
	for addrRange, device := range b.devicesMap {
		if addrRange.LowAddr <= addr && addr <= addrRange.HighAddr {
			return (*device).ReadMemory(addr)
		}
	}
	return 0x00
}

func (b *MainBus) WriteMemory(addr Word, value byte) {
	for addrRange, device := range b.devicesMap {
		if addrRange.LowAddr <= addr && addr <= addrRange.HighAddr {
			(*device).WriteMemory(addr, value)
		}
	}
}

func (b *MainBus) RegisterDevice(addrRange AddrRange, deviceStruct MemoryBus) *MainBus {
	b.devicesMap[addrRange] = &deviceStruct
	return b
}

func CreateMainBus() MainBus {
	return MainBus{
		devicesMap: map[AddrRange]*MemoryBus{},
	}
}
