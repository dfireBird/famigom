package bus

import (
	. "github.com/dfirebird/famigom/types"
)

type MainBusDevice interface {
	ReadMemory(addr Word) (bool, byte)
	WriteMemory(addr Word, value byte)
}

type MainBus struct {
	devicesMap []*MainBusDevice
}

func (b *MainBus) ReadMemory(addr Word) byte {
	for _, device := range b.devicesMap {
		if isRead, value := (*device).ReadMemory(addr); isRead {
			return value
		}
	}
	return 0x00
}

func (b *MainBus) WriteMemory(addr Word, value byte) {
	for _, device := range b.devicesMap {
		(*device).WriteMemory(addr, value)
	}
}

func (b *MainBus) RegisterDevice(deviceStruct MainBusDevice) *MainBus {
	b.devicesMap = append(b.devicesMap, &deviceStruct)
	return b
}

func CreateMainBus() MainBus {
	return MainBus{
		devicesMap: []*MainBusDevice{},
	}
}
